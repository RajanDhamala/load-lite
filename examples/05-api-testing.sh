#!/bin/bash
# =============================================================================
# Example 05: API Testing Scenarios
# =============================================================================
# Real-world API testing scenarios combining different flags for various
# use cases like REST APIs, microservices, and observability testing.
# =============================================================================

# -----------------------------------------------------------------------------
# Scenario 1: REST API CRUD Operations
# -----------------------------------------------------------------------------

# GET - List resources
echo "Testing GET /users endpoint..."
./load-lite \
  --url http://localhost:8080/api/users \
  --method GET \
  --rps 20 \
  --duration 15

# POST - Create resource
echo "Testing POST /users endpoint..."
./load-lite \
  --url http://localhost:8080/api/users \
  --method POST \
  --body '{"name":"Test User","email":"test@example.com"}' \
  --headers 'Content-Type:application/json' \
  --rps 10 \
  --duration 15

# -----------------------------------------------------------------------------
# Scenario 2: Authenticated API Calls
# -----------------------------------------------------------------------------

# ./load-lite \
#   --url http://localhost:8080/api/protected/data \
#   --method GET \
#   --headers 'Authorization:Bearer eyJhbGciOiJIUzI1NiIs...,Accept:application/json' \
#   --rps 30 \
#   --duration 30

# -----------------------------------------------------------------------------
# Scenario 3: Health Check Monitoring
# -----------------------------------------------------------------------------

# Continuous health check (useful for monitoring)
# ./load-lite \
#   --url http://localhost:8080/health \
#   --rps 1 \
#   --duration 300

# Liveness probe simulation
# ./load-lite \
#   --url http://localhost:8080/healthz \
#   --rps 2 \
#   --duration 60

# Readiness probe simulation
# ./load-lite \
#   --url http://localhost:8080/ready \
#   --rps 2 \
#   --duration 60

# -----------------------------------------------------------------------------
# Scenario 4: Microservices Communication
# -----------------------------------------------------------------------------

# Service A calling Service B
# ./load-lite \
#   --url http://service-b:8080/internal/api \
#   --headers 'X-Service-Name:service-a,X-Request-ID:$(uuidgen)' \
#   --rps 50 \
#   --duration 120

# Gateway API
# ./load-lite \
#   --url http://api-gateway:8080/v1/products \
#   --headers 'X-API-Version:v1,X-Client-ID:mobile-app' \
#   --rps 100 \
#   --duration 60

# -----------------------------------------------------------------------------
# Scenario 5: Webhook Simulation
# -----------------------------------------------------------------------------

# Simulate incoming webhooks
# ./load-lite \
#   --url http://localhost:8080/webhooks/payment \
#   --method POST \
#   --body '{"event":"payment.completed","data":{"id":"pay_123","amount":99.99}}' \
#   --headers 'Content-Type:application/json,X-Webhook-Secret:whsec_abc123' \
#   --rps 5 \
#   --duration 60

# -----------------------------------------------------------------------------
# Scenario 6: Observability Testing
# -----------------------------------------------------------------------------

# Generate traffic to trigger Prometheus alerts
# ./load-lite \
#   --url http://localhost:8080/api/slow-endpoint \
#   --rps 50 \
#   --duration 120

# Create traffic patterns for Grafana dashboards
# ./load-lite \
#   --url http://localhost:8080/api/metrics-test \
#   --rps 100 \
#   --duration 300
