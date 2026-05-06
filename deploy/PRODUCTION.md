# Production Hardening Guide

## Security

### 1. Secrets Management
```bash
# Never commit secrets! Use:
export $(grep -v '^#' .env | xargs)  # Load from .env

# Or use vault:
VAULT_ADDR=https://vault.example.com vault read secret/agent-harness
```

### 2. TLS/SSL
```yaml
# nginx.conf
server {
    listen 443 ssl http2;
    ssl_certificate /etc/ssl/certs/server.crt;
    ssl_certificate_key /etc/ssl/private/server.key;
    ssl_protocols TLSv1.2 TLSv1.3;
}
```

### 3. Rate Limiting
```yaml
# nginx.conf
limit_req_zone $binary_remote_addr zone=api:10m rate=100r/s;
limit_req zone=api burst=50;
```

### 4. Security Headers
```yaml
add_header X-Frame-Options DENY;
add_header X-Content-Type-Options nosniff;
add_header X-XSS-Protection "1; mode=block";
add_header Content-Security-Policy "default-src 'self'";
```

## Production Checklist

- [ ] Use PostgreSQL instead of SQLite
- [ ] Enable TLS/SSL
- [ ] Set up monitoring (Datadog, Prometheus)
- [ ] Configure backup strategy
- [ ] Set up log aggregation
- [ ] Configure auto-scaling
- [ ] Use secrets manager (Vault, AWS Secrets Manager)
- [ ] Set up WAF
- [ ] Configure audit logging

## Environment Variables

```bash
# Required for production:
HARNESS_SECRET=<generate-random-32-chars>
DATABASE_URL=postgres://user:pass@host:5432/db
REDIS_URL=redis://:pass@host:6379

# Optional:
DATADOG_API_KEY=
PROMETHEUS_AUTH_TOKEN=
```