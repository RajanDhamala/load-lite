package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// Config holds all CLI configuration
type Config struct {
	URL      string
	Method   string
	RPS      int
	Duration int
	Body     string
	Headers  string
}

// Stats tracks request metrics
type Stats struct {
	mu        sync.Mutex
	latencies []time.Duration
	errors    int
	requests  int
}

func (s *Stats) recordSuccess(latency time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.latencies = append(s.latencies, latency)
	s.requests++
}

func (s *Stats) recordError() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.errors++
	s.requests++
}

func (s *Stats) calculate() Summary {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.latencies) == 0 {
		return Summary{Requests: s.requests, Errors: s.errors}
	}

	// Sort for percentile calculation
	sorted := make([]time.Duration, len(s.latencies))
	copy(sorted, s.latencies)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	var total time.Duration
	min := sorted[0]
	max := sorted[len(sorted)-1]

	for _, d := range sorted {
		total += d
	}

	avg := total / time.Duration(len(sorted))
	p95 := percentile(sorted, 0.95)
	p99 := percentile(sorted, 0.99)

	return Summary{
		Requests: s.requests,
		Errors:   s.errors,
		Avg:      avg,
		Min:      min,
		Max:      max,
		P95:      p95,
		P99:      p99,
	}
}

// Summary holds final calculated statistics
type Summary struct {
	Requests int
	Errors   int
	Avg      time.Duration
	Min      time.Duration
	Max      time.Duration
	P95      time.Duration
	P99      time.Duration
}

func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}

	idx := int(float64(len(sorted)) * p)
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

type TrafficGenerator struct {
	config *Config
	client *http.Client
	stats  *Stats
}

func NewTrafficGenerator(cfg *Config) *TrafficGenerator {
	return &TrafficGenerator{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		stats: &Stats{
			latencies: make([]time.Duration, 0, cfg.RPS*cfg.Duration),
		},
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
	var wg sync.WaitGroup

	fmt.Printf("Starting traffic generation: %d rps for %d seconds to %s\n",
		tg.config.RPS, tg.config.Duration, tg.config.URL)

	for {
		select {
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				tg.makeRequest()
			}()
		case <-timeout:
			wg.Wait()
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
	fmt.Printf("  Errors         : %d\n", summary.Errors)

	if summary.Requests > summary.Errors {
		fmt.Println()
		fmt.Println("  Latency")
		fmt.Println("  ------------------------------")
		fmt.Printf("  Avg            : %v\n", summary.Avg.Round(time.Millisecond))
		fmt.Printf("  Min            : %v\n", summary.Min.Round(time.Millisecond))
		fmt.Printf("  Max            : %v\n", summary.Max.Round(time.Millisecond))
		fmt.Printf("  P95            : %v\n", summary.P95.Round(time.Millisecond))
		fmt.Printf("  P99            : %v\n", summary.P99.Round(time.Millisecond))
	}

	fmt.Println("========================================")
}
