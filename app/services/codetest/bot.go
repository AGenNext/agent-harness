package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/harness/harness/app/services"
)

// CodeTesterBot - runs tests on PRs
//
// Commands:
//   - "test PR #123" - runs all tests
//   - "unit PR #123" - runs unit tests
//   - "integration PR #123" - runs integration tests
//   - "e2e PR #123" - runs e2e tests
//   - "snapshot PR #123" - runs snapshot tests
//
// Mattermost: @code-tester
type CodeTesterBot struct {
	notifier *services.Notifier
	aiKey    string
	skills  []string
}

// NewCodeTesterBot creates a new code-tester bot
func NewCodeTesterBot() *CodeTesterBot {
	bot := &CodeTesterBot{
		notifier: services.NewNotifier(
			os.Getenv("GITHUB_TOKEN"),
			services.WithMattermost(os.Getenv("CODE_TESTER_HOOK")),
			services.WithSlack(os.Getenv("SLACK_WEBHOOK_URL")),
			services.WithDiscord(os.Getenv("DISCORD_WEBHOOK_URL")),
		),
		aiKey: os.Getenv("OPENAI_API_KEY"),
		skills: []string{
			"test",       // Run all tests
			"unit",      // Unit tests
			"integration", // Integration tests
			"e2e",       // End-to-end tests
			"snapshot",  // Snapshot tests
		},
	}
	return bot
}

// HandleMessage processes test request
func (b *CodeTesterBot) HandleMessage(ctx context.Context, owner, repo string, prNum int) (*services.Result, error) {
	ghToken := os.Getenv("GITHUB_TOKEN")
	tester := services.NewCodeTesterService(b.aiKey, ghToken)

	result, err := tester.RunTests(ctx, owner, repo, prNum)
	if err != nil {
		return nil, err
	}

	b.notifier.NotifyResult(ctx, result)
	return result, nil
}

// HTTPHandler exposes the bot via HTTP
func (b *CodeTesterBot) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		b.handlePost(w, r)
	case "GET":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "🧪 code-tester bot ready")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (b *CodeTesterBot) handlePost(w http.ResponseWriter, r *http.Request) {
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
	bot := NewCodeTesterBot()
	http.HandleFunc("/code-tester", bot.HTTPHandler)
	log.Println("🧪 code-tester bot running on :8083")
	log.Fatal(http.ListenAndServe(":8083", nil))
}