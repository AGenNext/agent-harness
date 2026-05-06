# Policies

## Access Control

| Policy | Priority | Subject | Resource | Action | Condition |
|--------|----------|--------|----------|--------|-----------|
| admin-all | 100 | role:admin | * | * | - |
| dev-code-assist | 50 | role:developer | code-assist | read,execute | - |
| dev-deploy | 40 | role:developer | code-deploy | read,execute | env:dev,staging |
| viewer-read | 10 | role:viewer | * | read | - |
| default-deny | 0 | * | * | * | - |

## Security

```yaml
require_approval: true
allowed_envs: [dev, staging, prod]
allowed_skills: [fix, security, review, test, docker, k8s]
block_patterns:
  - "rm -rf /"
  - "DROP TABLE"
  - "curl | sh"
auto_timeout: 300
```

## Rate Limit

```yaml
window: 60s
max_requests: 60
burst: 10
```

## API

```bash
curl /api/policies
curl -X POST /api/policies/check -d '{"subject":"user:1","resource":"/code-assist","action":"execute"}'
curl -X POST /api/policies -d '{"name":"custom","effect":"allow",...}'
```