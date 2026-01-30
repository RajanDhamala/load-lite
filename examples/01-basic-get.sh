#!/bin/bash
# =============================================================================
# Example 01: Basic GET Request
# =============================================================================
# This is the simplest way to use load-lite. It sends GET requests to a URL
# at a specified rate for a given duration.
# =============================================================================

# Simple GET request - 10 requests per second for 10 seconds
./load-lite \
  --url http://localhost:8080/api \
  --rps 10 \
  --duration 10

# =============================================================================
# What this does:
# - Sends GET requests to http://localhost:8080/api
# - Rate: 10 requests per second
# - Duration: 10 seconds
# - Total requests: ~100
# =============================================================================

# More examples:

# Health check endpoint
# ./load-lite --url http://localhost:8080/health --rps 5 --duration 30

# API endpoint with path
# ./load-lite --url http://localhost:8080/api/v1/users --rps 20 --duration 60

# External service (be careful with rate!)
# ./load-lite --url https://api.example.com/status --rps 2 --duration 10
