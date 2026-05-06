package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/harness/harness/app/services"
)

// CodeReviewBot - reviews code from PRs
//
// Commands:
//   - "review PR #123" - reviews PR
//   - "security PR #123" - security scan
//   - "perf PR #123" - performance review
//   - "lgtm PR #123" - approves PR
//   - "request changes PR #123" - requests changes
//
// Mattermost: @code-review
type CodeReviewBot struct {
	notifier *services.Notifier
	aiKey    string
	skills  []string
}

// NewCodeReviewBot creates a new code-review bot
func NewCodeReviewBot() *CodeReviewBot {
	bot := &CodeReviewBot{
		notifier: services.NewNotifier(
			os.Getenv("GITHUB_TOKEN"),
			services.WithMattermost(os.Getenv("CODE_REVIEW_HOOK")),
			services.WithSlack(os.Getenv("SLACK_WEBHOOK_URL")),
			services.WithDiscord(os.Getenv("DISCORD_WEBHOOK_URL")),
		),
		aiKey: os.Getenv("OPENAI_API_KEY"),
		skills: []string{
			"review",      // General code review
			"security",   // Security scan
			"performance", // Performance analysis
			"accessibility", // a11y check
			"best-practices", // Best practices
		},
	}
	return bot
}

// HandleMessage processes review request
func (b *CodeReviewBot) HandleMessage(ctx context.Context, owner, repo string, prNum int) (*services.Result, error) {
	ghToken := os.Getenv("GITHUB_TOKEN")
	review := services.NewCodeReviewService(b.aiKey, ghToken)

	result, err := review.ReviewPR(ctx, owner, repo, prNum)
	if err != nil {
		return nil, err
	}

	b.notifier.NotifyResult(ctx, result)
	return result, nil
}

// HTTPHandler exposes the bot via HTTP
func (b *CodeReviewBot) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		b.handlePost(w, r)
	case "GET":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "🔍 code-review bot ready")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (b *CodeReviewBot) handlePost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Owner string `json:"owner"`
		Repo  string `json:"repo"`
		PR    int    `json:"pr"`
	}

	if err := services.ParseJSON(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := b.HandleMessage(r.Context(), req.Owner, req.Repo, req.PR)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	services.JSON(w, http.StatusOK, result)
}

func main() {
	bot := NewCodeReviewBot()
	http.HandleFunc("/code-review", bot.HTTPHandler)
	log.Println("🔍 code-review bot running on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}