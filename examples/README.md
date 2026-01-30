# Traffic Simulator Examples

This directory contains example scripts demonstrating various usage patterns for the traffic simulator.

## Examples Overview

| File | Description |
|------|-------------|
| [01-basic-get.sh](./01-basic-get.sh) | Simple GET requests - the basics |
| [02-post-json.sh](./02-post-json.sh) | POST requests with JSON body |
| [03-custom-headers.sh](./03-custom-headers.sh) | Authentication and custom headers |
| [04-high-traffic.sh](./04-high-traffic.sh) | High volume traffic generation |
| [05-api-testing.sh](./05-api-testing.sh) | Real-world API testing scenarios |

## Quick Reference

### Basic Syntax

```bash
./load-lite --url <URL> --method <METHOD> --rps <RATE> --duration <SECONDS>
```

### All Available Flags

```bash
./load-lite \
  --url http://localhost:8080/api \    # Target URL (required)
  --method POST \                       # HTTP method: GET or POST (default: GET)
  --rps 100 \                           # Requests per second (default: 400)
  --duration 60 \                       # Duration in seconds (default: 10)
  --body '{"key":"value"}' \            # Request body for POST
  --headers 'Key:Value,Key2:Value2'     # Custom headers (comma-separated)
```

## Running Examples

1. **Make scripts executable:**
   ```bash
   chmod +x examples/*.sh
   ```

2. **Build the traffic simulator:**
   ```bash
   go build -o load-lite
   ```

3. **Run an example:**
   ```bash
   ./examples/01-basic-get.sh
   ```

## Common Patterns

### GET Request
```bash
./load-lite --url http://localhost:8080/api --rps 10 --duration 30
```

### POST with JSON
```bash
./load-lite \
  --url http://localhost:8080/api \
  --method POST \
  --body '{"name":"test"}' \
  --headers 'Content-Type:application/json'
```

### With Authentication
```bash
./load-lite \
  --url http://localhost:8080/api \
  --headers 'Authorization:Bearer your-token-here'
```

### Multiple Headers
```bash
./load-lite \
  --url http://localhost:8080/api \
  --headers 'Content-Type:application/json,Authorization:Bearer token,X-Custom:value'
```

## Tips

- **Start Low**: Begin with low RPS (5-10) and increase gradually
- **Monitor Target**: Watch your target system's metrics while testing
- **Use Appropriate Duration**: Short tests (10-30s) for quick checks, longer (60-300s) for stability
- **Headers Format**: Use `Key:Value` format, separate multiple with commas (no spaces around commas)
