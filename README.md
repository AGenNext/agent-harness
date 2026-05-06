# Code Commit GitHub Extension

**10 ways to run**: 4 platforms × multiple options

## 4 Agents

| Agent | Purpose |
|------|---------|
| code-assist | Write code fixes |
| code-review | Review code |
| code-tester | Run tests |
| code-deploy | Deploy |

---

## 10 Options (4 agents × platforms)

### Codespaces (Cloud/Sandbox)

| Agent | Command |
|------|---------|
| code-assist | `python src/code_assist/server.py` |
| code-review | `python src/code_review/server.py` |
| code-tester | `python src/code_tester/server.py` |
| code-deploy | `python src/code_deploy/server.py` |

### CLI (Terminal)

| Agent | Command |
|------|---------|
| code-assist | `code-commit assist "fix login bug"` |
| code-review | `code-commit review` |
| code-tester | `code-commit test` |
| code-deploy | `code-commit deploy` |

### Desktop (App)

**One app** = GitHub + VS Code + Docker + **4 AI Agents**

| Feature | Tool |
|---------|------|
| Repo/PR/Issues | GitHub API |
| Code editor | Monaco (VS Code) |
| Containers | Docker Engine |
| Write code | code-assist |
| Review code | code-review |
| Run tests | code-tester |
| Deploy | code-deploy

```bash
# Build
code-commit.exe
```

---

## Quick Start

### 1. Choose your option above

### 2. Add secrets:

| Secret | Required | Description |
|--------|----------|-------------|
| `GITHUB_TOKEN` | Yes | GitHub PAT (repo scope) |
| `OPENAI_API_KEY` | No | OpenAI key |
| `ANTHROPIC_API_KEY` | No | Anthropic key (Claude) |

### 3. Configure webhook in GitHub:

- **Payload URL**: `https://your-server.com/webhook`
- **Events**: Issues, Issue comments, Pull requests

---

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `GITHUB_TOKEN` | Yes | GitHub Personal Access Token |
| `OPENAI_API_KEY` | No | OpenAI API key |
| `ANTHROPIC_API_KEY` | No | Anthropic API key |
| `AUTO_MERGE_ENABLED` | No | Auto-merge PRs (default: true) |

---

## Triggering Fixes

### Option 1: Add Label
Add label `auto-fix` or `fix-me` to an issue.

### Option 2: Comment
Comment `/fix` on an issue.

---

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /webhook` | GitHub webhook receiver |
| `GET /health` | Health check |
| `POST /fix` | Code assist fix endpoint |

---

## License

MIT