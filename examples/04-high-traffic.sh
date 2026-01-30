#!/bin/bash
# =============================================================================
# Example 04: High Traffic Simulation
# =============================================================================
# Generate high volume traffic for testing system behavior under load.
# Useful for observability demos, auto-scaling triggers, and capacity testing.
# =============================================================================

# WARNING: High RPS can impact your system. Start low and increase gradually.

# Medium load - 100 RPS for 1 minute
./load-lite \
  --url http://localhost:8080/api \
  --rps 100 \
  --duration 60

# =============================================================================
# RPS Guidelines:
# - Low:    1-50 RPS    (basic testing, smoke tests)
# - Medium: 50-200 RPS  (performance baseline)
# - High:   200-500 RPS (stress testing, scaling triggers)
# =============================================================================

# Gradual load increase (run these sequentially)
# echo "Phase 1: Warm up - 50 RPS"
# ./load-lite --url http://localhost:8080/api --rps 50 --duration 30

# echo "Phase 2: Medium load - 150 RPS"
# ./load-lite --url http://localhost:8080/api --rps 150 --duration 60

# echo "Phase 3: High load - 300 RPS"
# ./load-lite --url http://localhost:8080/api --rps 300 --duration 60

# Sustained load for auto-scaling tests (5 minutes)
# ./load-lite \
#   --url http://localhost:8080/api \
#   --rps 200 \
#   --duration 300

# Quick burst - high RPS for short duration
# ./load-lite \
#   --url http://localhost:8080/api \
#   --rps 500 \
#   --duration 10

# Long-running test for stability (10 minutes at moderate load)
# ./load-lite \
#   --url http://localhost:8080/api \
#   --rps 75 \
#   --duration 600
