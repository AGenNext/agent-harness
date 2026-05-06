#!/bin/bash
# Deploy agent-harness to VPS

set -e

# Configuration
VPS_HOST="${VPS_HOST:-}"          # e.g., 192.168.1.100
VPS_USER="${VPS_USER:-root}"
VPS_PORT="${VPS_PORT:-22}"
DOMAIN="${DOMAIN:-}"               # e.g., harness.yourdomain.com

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Check required vars
check_vars() {
    if [[ -z "$VPS_HOST" ]]; then
        log_error "VPS_HOST not set"
        exit 1
    fi
}

# Deploy to VPS
deploy() {
    log_info "Building Docker image..."
    docker build -t harness:local .
    
    log_info "Stopping existing container..."
    docker stop harness 2>/dev/null || true
    docker rm harness 2>/dev/null || true
    
    log_info "Starting harness..."
    docker run -d \
        --name harness \
        -p 3000:3000 \
        -p 3022:3022 \
        -v /var/run/docker.sock:/var/run/docker.sock \
        -v harness_data:/data \
        -e HARNESS_SECRET="$HARNESS_SECRET" \
        --restart unless-stopped \
        harness:local
    
    log_info "Harness deployed at http://$VPS_HOST:3000"
}

# Docker Swarm deploy
deploy_swarm() {
    log_info "Deploying to Docker Swarm..."
    docker stack deploy -c docker-compose.yaml harness
    log_info "Harness deployed via swarm"
}

# Show status
status() {
    docker ps | grep harness || echo "Harness not running"
}

# Show logs
logs() {
    docker logs -f harness
}

# Help
help() {
    echo "Usage: $0 <command>"
    echo
    echo "Commands:"
    echo "  deploy      Build & deploy to local Docker"
    echo "  deploy-swarm  Deploy via Docker Swarm"
    echo "  status      Show container status"
    echo "  logs       Show container logs"
    echo
    echo "Environment:"
    echo "  VPS_HOST         VPS IP/hostname"
    echo "  VPS_USER         SSH user (default: root)"
    echo "  HARNESS_SECRET   Secret key"
    echo "  DOMAIN           Domain for SSL"
}

COMMAND="${1:-help}"
shift || true

case "$COMMAND" in
    deploy)
        check_vars
        deploy
        ;;
    deploy-swarm)
        deploy_swarm
        ;;
    status)
        status
        ;;
    logs)
        logs
        ;;
    *)
        help
        ;;
esac