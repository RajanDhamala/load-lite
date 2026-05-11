package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Config struct {
	URL         string
	Method      string
	RPS         int
	Duration    int
	Timeout     int
	Concurrency int
	Body        string
	Headers     string
}

type Stats struct {
	mu           sync.Mutex
	latencyHist  []int
	requests     int
	successes    int
	errors       int
	dropped      int
	totalLatency time.Duration
	minLatency   time.Duration
	maxLatency   time.Duration
}

func (s *Stats) recordSuccess(latency time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.requests++
	s.successes++
	s.totalLatency += latency

	if s.successes == 1 || latency < s.minLatency {
		s.minLatency = latency
	}
	if latency > s.maxLatency {
		s.maxLatency = latency
	}

	bucket := int(latency / time.Millisecond)
	if bucket >= len(s.latencyHist) {
		bucket = len(s.latencyHist) - 1
	}
	s.latencyHist[bucket]++
}

func (s *Stats) recordError() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.errors++
	s.requests++
}

func (s *Stats) recordDropped() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dropped++
}

func (s *Stats) calculate() Summary {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.successes == 0 {
		return Summary{Requests: s.requests, Errors: s.errors, Dropped: s.dropped}
	}

	return Summary{
		Requests:  s.requests,
		Successes: s.successes,
		Errors:    s.errors,
		Dropped:   s.dropped,
		Avg:       s.totalLatency / time.Duration(s.successes),
		Min:       s.minLatency,
		Max:       s.maxLatency,
		P95:       percentile(s.latencyHist, s.successes, 95),
		P99:       percentile(s.latencyHist, s.successes, 99),
	}
}

type Summary struct {
	Requests  int
	Successes int
	Errors    int
	Dropped   int
	Avg       time.Duration
	Min       time.Duration
	Max       time.Duration
	P95       time.Duration
	P99       time.Duration
}

func NewStats(timeout time.Duration) *Stats {
	buckets := int(timeout/time.Millisecond) + 2
	if buckets < 2 {
		buckets = 2
	}

	return &Stats{
		latencyHist: make([]int, buckets),
	}
}

func percentile(hist []int, samples int, percentile int) time.Duration {
	if samples == 0 || len(hist) == 0 {
		return 0
	}

	rank := (samples*percentile + 99) / 100
	seen := 0
	for bucket, count := range hist {
		seen += count
		if seen >= rank {
			return time.Duration(bucket) * time.Millisecond
		}
	}

	return time.Duration(len(hist)-1) * time.Millisecond
}

type TrafficGenerator struct {
	config *Config
	client *http.Client
	stats  *Stats
}

func NewTrafficGenerator(cfg *Config) *TrafficGenerator {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5
	}
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = cfg.RPS
	}

	timeout := time.Duration(cfg.Timeout) * time.Second
	idleConns := cfg.Concurrency
	if idleConns < 1 {
		idleConns = 1
	}

	return &TrafficGenerator{
		config: cfg,
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        idleConns,
				MaxIdleConnsPerHost: idleConns,
				MaxConnsPerHost:     cfg.Concurrency,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		stats: NewStats(timeout),
	}
}

func (tg *TrafficGenerator) makeRequest() {
	start := time.Now()

	req, err := tg.buildRequest()
	if err != nil {
		tg.stats.recordError()
		return
	}

	resp, err := tg.client.Do(req)
	if err != nil {
		tg.stats.recordError()
		return
	}
	defer resp.Body.Close()

	_, _ = io.Copy(io.Discard, resp.Body)

	latency := time.Since(start)

	if resp.StatusCode >= 400 {
		tg.stats.recordError()
	} else {
		tg.stats.recordSuccess(latency)
	}
}

func (tg *TrafficGenerator) buildRequest() (*http.Request, error) {
	var body io.Reader
	if tg.config.Method == "POST" && tg.config.Body != "" {
		body = bytes.NewBufferString(tg.config.Body)
	}

	req, err := http.NewRequest(tg.config.Method, tg.config.URL, body)
	if err != nil {
		return nil, err
	}

	if tg.config.Headers != "" {
		headers := strings.Split(tg.config.Headers, ",")
		for _, h := range headers {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) == 2 {
				req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}
	}

	return req, nil
}

func (tg *TrafficGenerator) Run() Summary {
	duration := time.Duration(tg.config.Duration) * time.Second
	interval := time.Second / time.Duration(tg.config.RPS)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	timeout := time.After(duration)
	jobs := make(chan struct{})
	var workers sync.WaitGroup
	for i := 0; i < tg.config.Concurrency; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for range jobs {
				tg.makeRequest()
			}
		}()
	}

	fmt.Printf("Starting traffic generation: %d rps for %d seconds to %s (timeout=%ds, concurrency=%d)\n",
		tg.config.RPS, tg.config.Duration, tg.config.URL, tg.config.Timeout, tg.config.Concurrency)

	for {
		select {
		case <-ticker.C:
			select {
			case jobs <- struct{}{}:
			default:
				tg.stats.recordDropped()
			}
		case <-timeout:
			close(jobs)
			workers.Wait()
			return tg.stats.calculate()
		}
	}
}

func parseConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.URL, "url", "http://localhost:8001/devices", "Target URL")
	flag.StringVar(&cfg.Method, "method", "GET", "HTTP method (GET or POST)")
	flag.IntVar(&cfg.RPS, "rps", 400, "Requests per second")
	flag.IntVar(&cfg.Duration, "duration", 10, "Duration in seconds")
	flag.IntVar(&cfg.Timeout, "timeout", 5, "Request timeout in seconds")
	flag.IntVar(&cfg.Concurrency, "concurrency", 0, "Max concurrent in-flight requests (default: rps)")
	flag.StringVar(&cfg.Body, "body", "", "Request body for POST")
	flag.StringVar(&cfg.Headers, "headers", "", "Headers (comma-separated, e.g. 'Content-Type:application/json,X-Custom:value')")

	flag.Parse()

	// Validate
	cfg.Method = strings.ToUpper(cfg.Method)
	if cfg.Method != "GET" && cfg.Method != "POST" {
		fmt.Println("Error: method must be GET or POST")
		flag.Usage()
		return nil
	}

	if cfg.RPS <= 0 || cfg.Duration <= 0 {
		fmt.Println("Error: rps and duration must be positive")
		flag.Usage()
		return nil
	}
	if time.Second/time.Duration(cfg.RPS) <= 0 {
		fmt.Println("Error: rps is too high")
		flag.Usage()
		return nil
	}

	if cfg.Timeout <= 0 {
		fmt.Println("Error: timeout must be positive")
		flag.Usage()
		return nil
	}

	if cfg.Concurrency < 0 {
		fmt.Println("Error: concurrency cannot be negative")
		flag.Usage()
		return nil
	}
	if cfg.Concurrency == 0 {
		cfg.Concurrency = cfg.RPS
	}

	return cfg
}

func main() {
	cfg := parseConfig()
	if cfg == nil {
		return
	}

	generator := NewTrafficGenerator(cfg)
	summary := generator.Run()
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("         Traffic Test Summary")
	fmt.Println("========================================")

	fmt.Printf("  Requests       : %d\n", summary.Requests)
	fmt.Printf("  Successes      : %d\n", summary.Successes)
	fmt.Printf("  Errors         : %d\n", summary.Errors)
	fmt.Printf("  Dropped        : %d\n", summary.Dropped)

	if summary.Successes > 0 {
		fmt.Println()
		fmt.Println("  Successful Response Latency")
		fmt.Println("  ------------------------------")
		fmt.Printf("  Avg            : %v\n", summary.Avg.Round(time.Millisecond))
		fmt.Printf("  Min            : %v\n", summary.Min.Round(time.Millisecond))
		fmt.Printf("  Max            : %v\n", summary.Max.Round(time.Millisecond))
		fmt.Printf("  P95            : %v\n", summary.P95.Round(time.Millisecond))
		fmt.Printf("  P99            : %v\n", summary.P99.Round(time.Millisecond))
	}

	fmt.Println("========================================")
}
