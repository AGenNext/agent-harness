# Agent Configuration

## LLM Providers

| Provider | Env Variable | Model | Config |
|----------|-------------|-------|--------|
| `openai` | `OPENAI_API_KEY` | GPT-4, GPT-3.5 | Model via `OPENAI_MODEL` |
| `anthropic` | `ANTHROPIC_API_KEY` | Claude 3.5 Sonnet | Model via `ANTHROPIC_MODEL` |
| `local` | `LOCAL_LLM_URL` | Any Ollama model | Model via `LOCAL_MODEL` |
| `cohere` | `COHERE_API_KEY` | Command R | - |
| `google` | `GOOGLE_API_KEY` | Gemini Pro | - |
| `azure` | `AZURE_OPENAI_KEY` | GPT-4 Azure | Endpoint via `AZURE_ENDPOINT` |

---

## Notification Channels

| Channel | Env Variable | Example |
|---------|-------------|---------|
| **Slack** | `SLACK_WEBHOOK_URL` | `https://hooks.slack.com/services/...` |
| **Discord** | `DISCORD_WEBHOOK_URL` | `https://discord.com/api/webhooks/...` |
| **Mattermost** | `CODE_ASSIST_HOOK` | `https://mattermost.../hooks/code-assist` |
| **WhatsApp** | `WHATSAPP_PHONE` | `+1234567890` |
| **Email** | `EMAIL_TO` | `team@example.com` |
| **Webhook** | `GENERIC_WEBHOOK_URL` | `https://myapp.com/hook` |

---

## Agent Ports

| Agent | Port | Env Hook |
|-------|------|---------|
| code-assist | 8081 | `CODE_ASSIST_HOOK` |
| code-review | 8082 | `CODE_REVIEW_HOOK` |
| code-tester | 8083 | `CODE_TESTER_HOOK` |
| code-deploy | 8084 | `CODE_DEPLOY_HOOK` |

---

## Agent Skills

### code-assist
```
fix, security, docs, refactor, optimize
```

### code-review
```
review, security, performance, accessibility, best-practices
```

### code-tester
```
test, unit, integration, e2e, snapshot
```

### code-deploy
```
docker, k8s, serverless, rollback, preview
```

---

## Example .env

```bash
# ===== LLM =====
LLM_PROVIDER=anthropic
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
LOCAL_LLM_URL=http://localhost:11434

# ===== GitHub =====
GITHUB_TOKEN=ghp_...

# ===== Slack =====
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/...

# ===== Discord =====
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/...

# ===== Mattermost =====
CODE_ASSIST_HOOK=https://mattermost.example.com/hooks/code-assist
CODE_REVIEW_HOOK=https://mattermost.example.com/hooks/code-review
CODE_TESTER_HOOK=https://mattermost.example.com/hooks/code-tester
CODE_DEPLOY_HOOK=https://mattermost.example.com/hooks/code-deploy

# ===== WhatsApp =====
WHATSAPP_PHONE=+1234567890

# ===== Email =====
EMAIL_TO=team@example.com
SMTP_HOST=smtp.example.com
SMTP_PORT=587

# ===== Generic =====
GENERIC_WEBHOOK_URL=https://myapp.com/hook

# ===== Optional =====
AUTO_MERGE_ENABLED=true
DEBUG=false
```

---

## Quick Start

```bash
# 1. Copy config
cp .env.example .env

# 2. Edit with your keys
nano .env

# 3. Start all agents
docker-compose up -d

# 4. Check status
curl localhost:3000/health
```