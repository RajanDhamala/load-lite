#!/bin/bash
# =============================================================================
# Example 03: Custom Headers & Authentication
# =============================================================================
# Add custom HTTP headers including authentication tokens, API keys,
# and other custom headers required by your API.
# =============================================================================

# Bearer token authentication
./load-lite \
  --url http://localhost:8080/api/protected \
  --headers 'Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...' \
  --rps 10 \
  --duration 20

# =============================================================================
# Headers format:
# - Single header:   'Key:Value'
# - Multiple headers: 'Key1:Value1,Key2:Value2,Key3:Value3'
# - No spaces around the colon
# =============================================================================

# API Key authentication
# ./load-lite \
#   --url http://api.example.com/data \
#   --headers 'X-API-Key:your-api-key-here' \
#   --rps 5 \
#   --duration 30

# Multiple custom headers
# ./load-lite \
#   --url http://localhost:8080/api/data \
#   --headers 'Authorization:Bearer token123,X-Request-ID:req-001,X-Client-Version:1.0.0' \
#   --rps 10 \
#   --duration 15

# POST with authentication and JSON
# ./load-lite \
#   --url http://localhost:8080/api/items \
#   --method POST \
#   --body '{"name":"New Item","price":19.99}' \
#   --headers 'Content-Type:application/json,Authorization:Bearer mytoken,X-Tenant-ID:tenant123' \
#   --rps 20 \
#   --duration 60

# Basic Auth style (base64 encoded credentials)
# ./load-lite \
#   --url http://localhost:8080/api/secure \
#   --headers 'Authorization:Basic dXNlcm5hbWU6cGFzc3dvcmQ=' \
#   --rps 5 \
#   --duration 10

# Custom tracing headers (useful for observability)
# ./load-lite \
#   --url http://localhost:8080/api/trace \
#   --headers 'X-Trace-ID:trace-12345,X-Span-ID:span-67890,X-Parent-ID:parent-11111' \
#   --rps 25 \
#   --duration 30
