# Docker Hub Push

## GitHub Actions (Automatic)

Set these secrets in your repo:
- `DOCKER_USERNAME` - Docker Hub username
- `DOCKER_PASSWORD` - Docker Hub password or access token

The workflow will auto-push on push to `main`.

## Manual Push

```bash
# Build
docker build -t agent-harness .

# Tag for Docker Hub
docker tag agent-harness:latest $DOCKER_USERNAME/agent-harness:latest
docker tag agent-harness:latest $DOCKER_USERNAME/agent-harness:v1.0.0

# Login
docker login

# Push
docker push $DOCKER_USERNAME/agent-harness:latest
docker push $DOCKER_USERNAME/agent-harness:v1.0.0
```

## GHCR (GitHub Container Registry)

```bash
# Tag
docker tag agent-harness:latest ghcr.io/$USER/agent-harness:latest

# Login (use GitHub token as password)
echo $GITHUB_TOKEN | docker login ghcr.io -U $USER --password-stdin

# Push
docker push ghcr.io/$USER/agent-harness:latest
```