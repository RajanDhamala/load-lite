#!/bin/bash
# =============================================================================
# Example 02: POST Request with JSON Body
# =============================================================================
# Send POST requests with JSON payload. Useful for testing APIs that create
# or update resources.
# =============================================================================

# POST request with JSON body
./load-lite \
  --url http://localhost:8080/api/users \
  --method POST \
  --body '{"username":"testuser","email":"test@example.com"}' \
  --headers 'Content-Type:application/json' \
  --rps 5 \
  --duration 10

# =============================================================================
# Flags explained:
# --method POST     : Use HTTP POST instead of GET
# --body '...'      : The request body to send
# --headers '...'   : Set Content-Type header for JSON
# =============================================================================

# More POST examples:

# Create order
# ./load-lite \
#   --url http://localhost:8080/api/orders \
#   --method POST \
#   --body '{"product_id":123,"quantity":2,"price":29.99}' \
#   --headers 'Content-Type:application/json' \
#   --rps 10 \
#   --duration 30

# Submit form data (as JSON)
# ./load-lite \
#   --url http://localhost:8080/api/contact \
#   --method POST \
#   --body '{"name":"John Doe","message":"Hello!"}' \
#   --headers 'Content-Type:application/json' \
#   --rps 3 \
#   --duration 20

# Nested JSON object
# ./load-lite \
#   --url http://localhost:8080/api/events \
#   --method POST \
#   --body '{"event":"click","data":{"button":"submit","page":"/checkout"}}' \
#   --headers 'Content-Type:application/json' \
#   --rps 50 \
#   --duration 60
