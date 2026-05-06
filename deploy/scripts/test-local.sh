#!/bin/bash
# Local test without Docker

echo "🧪 Agent-Harness Local E2E Test"
echo "========================================="

# Start harness in background
echo "[1/5] Starting Harness kernel..."
cd /workspace/project
gitness --port 3000 --docker-host /var/run/docker.sock &
HARNESS_PID=$!
sleep 5

# Test health
echo "[2/5] Testing health..."
curl -sf http://localhost:3000/health && echo " ✅" || echo " ❌"

# Test agents API  
echo "[3/5] Testing agents..."
curl -sf http://localhost:3000/api/v1/agents | head -c 100 && echo " ✅" || echo " ❌"

# Test workflows
echo "[4/5] Testing workflows..."
curl -sf http://localhost:3000/api/v1/workflows | head -c 100 && echo " ✅" || echo " ❌"

# Run test script
echo "[5/5] Running full tests..."
bash /workspace/project/deploy/scripts/test-e2e.sh

echo "========================================="
echo "Done! Check results above"

# Cleanup
kill $HARNESS_PID 2>/dev/null
