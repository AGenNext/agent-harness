#!/bin/bash
# Deploy to Harbor

HARBOR_URL="${HARBOR_URL:-harbor.autonomyx.io}"
HARBOR_USER="${HARBOR_USER:-admin}"
HARBOR_PASS="${HARBOR_PASS:-}"

if [ -z "$HARBOR_PASS" ]; then
  echo "Usage: HARBOR_URL=harbor.example.com HARBOR_USER=user HARBOR_PASS=pass ./deploy-harbor.sh"
  exit 1
fi

echo "$HARBOR_PASS" | podman login "$HARBOR_URL" -u "$HARBOR_USER" --password-stdin
podman build -t "$HARBOR_URL/agent-harness:latest" .
podman push "$HARBOR_URL/agent-harness:latest"
