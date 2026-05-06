package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/harness/harness/app/services"
)

// CodeDeployBot - deploys to environments
//
// Commands:
//   - "deploy to prod" - deploys to production
//   - "deploy PR #123" - deploys PR preview
//   - "docker deploy" - Docker deployment
//   - "k8s deploy" - Kubernetes deployment
//   - "serverless deploy" - Lambda/CloudFunctions
//   - "rollback" - rollback deployment
//
// Mattermost: @code-deploy
type CodeDeployBot struct {
	notifier *services.Notifier
	aiKey    string
	skills  []string
}

// NewCodeDeployBot creates a new code-deploy bot
func NewCodeDeployBot() *CodeDeployBot {
	bot := &CodeDeployBot{
		notifier: services.NewNotifier(
			os.Getenv("GITHUB_TOKEN"),
			services.WithMattermost(os.Getenv("CODE_DEPLOY_HOOK")),
			services.WithSlack(os.Getenv("SLACK_WEBHOOK_URL")),
			services.WithDiscord(os.Getenv("DISCORD_WEBHOOK_URL")),
		),
		aiKey: os.Getenv("OPENAI_API_KEY"),
		skills: []string{
			"docker",     // Docker deployment
			"k8s",      // Kubernetes
			"serverless", // Lambda/CloudFunctions
			"rollback", // Rollback
			"preview",  // Preview deployment
		},
	}
	return bot
}

// HandleMessage processes deploy request
func (b *CodeDeployBot) HandleMessage(ctx context.Context, owner, repo, env string, prNum int) (*services.Result, error) {
	ghToken := os.Getenv("GITHUB_TOKEN")
	deployer := services.NewCodeDeployService(b.aiKey, ghToken)

	result, err := deployer.Deploy(ctx, owner, repo, env, prNum)
	if err != nil {
		return nil, err
	}

	b.notifier.NotifyResult(ctx, result)
	return result, nil
}

// HTTPHandler exposes the bot via HTTP
func (b *CodeDeployBot) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		b.handlePost(w, r)
	case "GET":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "🚀 code-deploy bot ready")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (b *CodeDeployBot) handlePost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Owner string `json:"owner"`
		Repo  string `json:"repo"`
		Env   string `json:"env"`
		PR   int    `json:"pr"`
	}

	if err := services.ParseJSON(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := b.HandleMessage(r.Context(), req.Owner, req.Repo, req.Env, req.PR)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	services.JSON(w, http.StatusOK, result)
}

func main() {
	bot := NewCodeDeployBot()
	http.HandleFunc("/code-deploy", bot.HTTPHandler)
	log.Println("🚀 code-deploy bot running on :8084")
	log.Fatal(http.ListenAndServe(":8084", nil))
}