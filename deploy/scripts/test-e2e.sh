#!/bin/bash
# E2E Test for Agent-Harness

set -e

HARNESS_URL="${HARNESS_URL:-http://localhost:3000}"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Test Harness health
test_health() {
    log_info "Testing Harness kernel..."
    curl -sf "${HARNESS_URL}/health" && log_info "✅ Healthy" || log_error "❌ Failed"
}

# Test agents API
test_agents() {
    log_info "Testing agents API..."
    curl -sf "${HARNESS_URL}/api/v1/agents" && log_info "✅ Agents API" || log_error "❌ Failed"
}

# Test ports
test_ports() {
    log_info "Testing agent ports..."
    for port in 8081 8082 8083 8084; do
        curl -sf "http://localhost:${port}" 2>/dev/null && echo "Port ${port} ✅" || echo "Port ${port} ❌"
    done
}

# Test runtimes
test_runtimes() {
    log_info "Testing runtimes API..."
    curl -sf "${HARNESS_URL}/api/v1/runtimes" && log_info "✅ Runtimes" || log_error "❌ Failed"
}

# Main
main() {
    echo "========================================="
    echo "🧪 Agent-Harness E2E"
    echo "========================================="
    test_health
    test_agents
    test_ports
    test_runtimes
    echo "========================================="
    echo "✅ Done"
}

main "$@"