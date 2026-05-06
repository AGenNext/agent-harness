# Deploy agent-harness

## Quick Start

### 1. Docker (Local)

```bash
docker build -t harness .
docker run -d \
    --name harness \
    -p 3000:3000 \
    -p 3022:3022 \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v harness_data:/data \
    -e HARNESS_SECRET=your-secret \
    --restart unless-stopped \
    harness
```

### 2. Docker Compose

```bash
cp deploy/.env.example deploy/.env
# Edit .env with your settings

docker-compose up -d
```

### 3. VPS (systemd)

```bash
# Copy service file
sudo cp deploy/harness.service /etc/systemd/system/

# Copy env file
sudo cp deploy/.env /opt/harness/.env

# Start
sudo systemctl enable harness
sudo systemctl start harness

# Check status
sudo systemctl status harness
```

### 4. Cloud (Docker Hub / GHCR)

```bash
# Build & push
docker build -t ghcr.io/agennext/agent-harness:latest .
docker push ghcr.io/agennext/agent-harness:latest

# Pull & run on VPS
docker pull ghcr.io/agennext/agent-harness:latest
docker run -d ...
```

## Ports

| Port | Service |
|------|---------|
| 3000 | Web UI |
| 3022 | SSH/Git |

## Environment

| Variable | Required | Description |
|----------|----------|-------------|
| `HARNESS_SECRET` | Yes | Secret key |
| `HARNESS_DEBUG` | No | Debug mode |
| `DOCKER_HOST` | No | Docker socket |

## SSL (optional)

Use nginx or cloudflare for SSL termination.