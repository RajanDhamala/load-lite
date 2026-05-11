<p align="center">
  <h1 align="center"> Load Lite</h1>
  <p align="center">
    A lightweight, production-grade CLI tool for generating controlled HTTP traffic
    <br />
    <a href="#quick-start">Quick Start</a>
    ·
    <a href="#cli-reference">CLI Reference</a>
    ·
    <a href="./examples">Examples</a>
  </p>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License">
  <img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg" alt="PRs Welcome">
</p>

---

## Overview

Traffic Simulator is a **traffic generation tool** (not a load tester) designed to:

-  Generate controlled HTTP traffic at specified rates
-  Compute accurate latency statistics (Avg, Min, Max, P95, P99)
-  Support observability demos and experiments
-  Create visible spikes in monitoring systems

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/RajanDhamala/load-lite.git
cd load-lite 

# Build the binary
go build -o load-lite ./main.go


# Or run directly
go run main.go --url http://localhost:8080 --rps 10 --duration 10
```

### Basic Usage

```bash
# Simple GET request at 10 requests/second for 30 seconds
./load-lite --url http://localhost:8080/api --rps 10 --duration 30
```

## CLI Reference

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--url` | string | `http://localhost:8001/devices` | Target URL to send requests |
| `--method` | string | `GET` | HTTP method (`GET` or `POST`) |
| `--rps` | int | `400` | Requests per second |
| `--duration` | int | `10` | Test duration in seconds |
| `--timeout` | int | `5` | Per-request timeout in seconds |
| `--concurrency` | int | `rps` | Maximum concurrent in-flight requests |
| `--body` | string | `""` | Request body (for POST requests) |
| `--headers` | string | `""` | Custom headers (comma-separated) |

### Flag Details

#### `--url`
The target endpoint to send traffic to.

```bash
./load-lite --url http://api.example.com/v1/users
./load-lite --url http://localhost:3000/health
./load-lite --url https://myservice.internal:8443/endpoint
```

#### `--method`
HTTP method for requests. Supports `GET` and `POST`.

```bash
# GET request (default)
./load-lite --url http://localhost:8080/api --method GET

# POST request
./load-lite --url http://localhost:8080/api --method POST --body '{"key":"value"}'
```

#### `--rps`
Requests Per Second - controls the traffic rate.

```bash
# Low traffic (10 req/s)
./load-lite --url http://localhost:8080 --rps 10

# Medium traffic (100 req/s)
./load-lite --url http://localhost:8080 --rps 100

# High traffic (500 req/s)
./load-lite --url http://localhost:8080 --rps 500
```

#### `--duration`
How long to run the test in seconds.

```bash
# Quick test (10 seconds)
./load-lite --url http://localhost:8080 --duration 10

# Extended test (5 minutes)
./load-lite --url http://localhost:8080 --duration 300
```

#### `--timeout`
Per-request timeout in seconds. The default is intentionally short so slow or broken targets do not keep the CLI alive for a long time after the traffic window ends.

```bash
# Fail slow requests after 2 seconds
./load-lite --url http://localhost:8080 --timeout 2
```

#### `--concurrency`
Maximum number of in-flight HTTP requests. Defaults to `--rps`.

If every worker is busy, the scheduled tick is counted as `Dropped` instead of starting unbounded goroutines. Increase this when your target has higher normal latency and you still want to sustain the configured RPS.

```bash
# Allow up to 500 active requests while generating 200 RPS
./load-lite --url http://localhost:8080 --rps 200 --concurrency 500
```

#### `--body`
Request body for POST requests. Use with `--method POST`.

```bash
# JSON body
./load-lite --url http://localhost:8080/api \
  --method POST \
  --body '{"username":"john","email":"john@example.com"}'

# Form data style
./load-lite --url http://localhost:8080/api \
  --method POST \
  --body 'name=John&age=30'
```

#### `--headers`
Custom HTTP headers. Multiple headers are comma-separated, key:value format.

```bash
# Single header
./load-lite --url http://localhost:8080 \
  --headers 'Authorization:Bearer mytoken123'

# Multiple headers
./load-lite --url http://localhost:8080 \
  --headers 'Content-Type:application/json,Authorization:Bearer token,X-Request-ID:abc123'

# API key authentication
./load-lite --url http://api.example.com \
  --headers 'X-API-Key:your-api-key-here'
```

## Output

```
Starting traffic generation: 100 rps for 30 seconds to http://localhost:8080/api (timeout=5s, concurrency=100)

========================================
         Traffic Test Summary
========================================
  Requests       : 3000
  Successes      : 2995
  Errors         : 5
  Dropped        : 0

  Successful Response Latency
  ------------------------------
  Avg            : 42ms
  Min            : 8ms
  Max            : 285ms
  P95            : 125ms
  P99            : 210ms
========================================
```

## Examples

See the [`examples/`](./examples) directory for complete usage examples:

- [Basic GET Request](./examples/01-basic-get.sh)
- [POST with JSON Body](./examples/02-post-json.sh)
- [Custom Headers & Auth](./examples/03-custom-headers.sh)
- [High Traffic Simulation](./examples/04-high-traffic.sh)
- [API Testing Scenarios](./examples/05-api-testing.sh)

## Use Cases

### 1. Validate Monitoring Systems
```bash
# Generate traffic spike visible in Prometheus/Grafana
./load-lite --url http://myapp:8080/api --rps 50 --duration 60
```

### 2. Test Latency SLOs
```bash
# Verify P95 latency meets requirements
./load-lite --url http://myapp:8080/health --rps 100 --duration 30
```

### 3. Demo Auto-scaling
```bash
# Create sustained load to trigger HPA/auto-scaling
./load-lite --url http://myapp/api --rps 200 --duration 300
```

### 4. Pre-deployment Smoke Test
```bash
# Quick validation of new deployment
./load-lite --url http://staging.myapp.com/health --rps 10 --duration 10
```

## Performance

| Metric | Value | Notes |
|--------|-------|-------|
| CPU | <5% | At 100 RPS on modern hardware |
| Memory | Bounded | Latency statistics use a fixed millisecond histogram based on `--timeout`, not one entry per request |
| Connections | Bounded, pooled & reused | `--concurrency` caps active requests and connection usage |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
