#!/bin/bash
# Add a new agent to the team

set -e

NAME="${1:-}"
SKILLS="${2:-default}"

if [[ -z "$NAME" ]]; then
    echo "Usage: $0 <agent-name> [skills...]"
    echo
    echo "Examples:"
    echo "  $0 security-bot scan,audit,pen-test"
    echo "  $0 docs-bot write,update,review"
    echo "  $0 infra-bot provision,configure,teardown"
    echo
    echo "Available skills per agent type:"
    echo "  code-assist: fix,security,docs,refactor,optimize"
    echo "  code-review: review,security,performance,a11y,best-practices"
    echo "  code-tester: test,unit,integration,e2e,snapshot"
    echo "  code-deploy: docker,k8s,serverless,rollback,preview"
    echo
    exit 1
fi

# Default skills mapping
declare -A DEFAULT_SKILLS=(
    ["code-assist"]="fix,security,docs,refactor,optimize"
    ["code-review"]="review,security,performance,a11y,best-practices"
    ["code-tester"]="test,unit,integration,e2e,snapshot"
    ["code-deploy"]="docker,k8s,serverless,rollback,preview"
    ["security-bot"]="scan,audit,pen-test,fix"
    ["docs-bot"]="write,update,review,translate"
    ["infra-bot"]="provision,configure,teardown,scale"
)

SKILLS_DIR="app/services/${NAME}"
BOT_FILE="${SKILLS_DIR}/bot.go"
SERVICE_FILE="${SKILLS_DIR}/service.go"

# Get skills
if [[ "$SKILLS" == "default" ]] || [[ -z "$SKILLS" ]]; then
    SKILLS="${DEFAULT_SKILLS[$NAME]:-skill1,skill2,skill3}"
fi

echo "🤖 Creating agent: $NAME"
echo "   Skills: $SKILLS"
echo

# Create directory
mkdir -p "$SKILLS_DIR"

# Find available port (8081-8099)
for port in $(seq 8081 8099); do
    if ! lsof -i:$port >/dev/null 2>&1; then
        PORT=$port
        break
    fi
done
PORT=${PORT:-8081}

# Create bot.go
cat > "$BOT_FILE" << EOF
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/harness/harness/app/services"
)

// ${NAME^}Bot - ${NAME//-/ } agent
//
// Skills: $(echo $SKILLS | tr ',' '\n' | sed 's/^/   - /')
type ${NAME^}Bot struct {
	notifier *services.Notifier
	aiKey   string
	skills  []string
}

// New${NAME^}Bot creates a new ${NAME//-/ } bot
func New${NAME^}Bot() *${NAME^}Bot {
	bot := &${NAME^}Bot{
		notifier: services.NewNotifier(
			os.Getenv("GITHUB_TOKEN"),
			services.WithMattermost(os.Getenv("${NAME^^}_HOOK")),
			services.WithSlack(os.Getenv("SLACK_WEBHOOK_URL")),
			services.WithDiscord(os.Getenv("DISCORD_WEBHOOK_URL")),
		),
		aiKey: os.Getenv("OPENAI_API_KEY"),
		skills: []string{$(echo $SKILLS | sed 's/\([^,]*\)/"\1"/g' | tr ',' ',')},
	}
	return bot
}

// HandleMessage processes request
func (b *${NAME^}Bot) HandleMessage(ctx context.Context, owner, repo string) (*services.Result, error) {
	ghToken := os.Getenv("GITHUB_TOKEN")
	svc := services.New${NAME^}Service(b.aiKey, ghToken)
	result, err := svc.Handle(ctx, owner, repo)
	if err != nil {
		return nil, err
	}
	b.notifier.NotifyResult(ctx, result)
	return result, nil
}

// HTTPHandler exposes the bot via HTTP
func (b *${NAME^}Bot) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		b.handlePost(w, r)
	case "GET":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "🤖 ${NAME//-/ } bot ready")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (b *${NAME^}Bot) handlePost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Owner string \`json:"owner"\`
		Repo  string \`json:"repo"\`
	}
	if err := services.ParseJSON(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := b.HandleMessage(r.Context(), req.Owner, req.Repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	services.JSON(w, http.StatusOK, result)
}

func main() {
	bot := New${NAME^}Bot()
	http.HandleFunc("/${NAME}", bot.HTTPHandler)
	log.Println("🤖 ${NAME//-/ } bot running on :${PORT}")
	log.Fatal(http.ListenAndServe(":${PORT}", nil))
}
EOF

# Create service.go
cat > "$SERVICE_FILE" << EOF
package services

import (
	"context"
	"fmt"
	"github.com/google/go-github/v50/github"
	"github.com/sashabaranov/go-openai"
)

// ${NAME^}Service - ${NAME//-/ } service
type ${NAME^}Service struct {
	aiClient *openai.Client
	ghClient *github.Client
}

// New${NAME^}Service creates a new ${NAME//-/ } service
func New${NAME^}Service(apiKey, ghToken string) *${NAME^}Service {
	return &${NAME^}Service{
		aiClient: openai.NewClient(apiKey),
		ghClient: github.NewClient(ghToken),
	}
}

// Handle processes the request
func (s *${NAME^}Service) Handle(ctx context.Context, owner, repo string) (*Result, error) {
	// Implement your logic here
	return &Result{
		Owner:  owner,
		Repo:   repo,
		Status: "ready",
	}, nil
}
EOF

echo "✅ Created: $BOT_FILE"
echo "✅ Created: $SERVICE_FILE"
echo
echo "📝 Next steps:"
echo "   1. Implement the service logic in $SERVICE_FILE"
echo "   2. Set environment variables:"
echo "      export ${NAME^^}_HOOK=https://mattermost.../hooks/${NAME}"
echo "   3. Build: go build -o bin/${NAME} ./app/services/${NAME}"
echo "   4. Run: ./bin/${NAME}"
echo
echo "📋 Available skills:"
echo "   $(echo $SKILLS | tr ',' ' ')"