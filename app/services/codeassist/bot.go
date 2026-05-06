package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/harness/harness/app/services"
)

// CodeAssistBot - creates code fixes from issues
// 
// Commands:
//   - "fix issue #123" - creates fix for issue
//   - "fix: <description>" - creates issue & fix
//
// Mattermost: @code-assist
type CodeAssistBot struct {
	notifier *services.Notifier
	aiKey    string
	
	// Skills available
	skills []string
}

// NewCodeAssistBot creates a new code-assist bot
func NewCodeAssistBot() *CodeAssistBot {
	bot := &CodeAssistBot{
		notifier: services.NewNotifier(
			os.Getenv("GITHUB_TOKEN"),
			services.WithMattermost(os.Getenv("CODE_ASSIST_HOOK")),
			services.WithSlack(os.Getenv("SLACK_WEBHOOK_URL")),
			services.WithDiscord(os.Getenv("DISCORD_WEBHOOK_URL")),
		),
		aiKey: os.Getenv("OPENAI_API_KEY"),
		skills: []string{
			"fix",        // Generate code fixes
			"security",  // Fix security vulnerabilities
			"docs",      // Generate documentation
			"refactor",  // Refactor code
			"optimize",  // Optimize performance
		},
	}
	return bot
}

// HandleMessage processes Mattermost message
func (b *CodeAssistBot) HandleMessage(ctx context.Context, owner, repo string, issueNum int) (*services.Result, error) {
	ghToken := os.Getenv("GITHUB_TOKEN")
	assist := services.NewCodeAssistService(b.aiKey, ghToken)

	// Generate fix
	result, err := assist.HandleIssue(ctx, owner, repo, issueNum)
	if err != nil {
		return nil, err
	}

	// Notify via Mattermost
	b.notifier.NotifyResult(ctx, result)

	return result, nil
}

// HTTPHandler exposes the bot via HTTP
func (b *CodeAssistBot) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		b.handlePost(w, r)
	case "GET":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "🤖 code-assist bot ready")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (b *CodeAssistBot) handlePost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Owner  string `json:"owner"`
		Repo  string `json:"repo"`
		Issue int    `json:"issue"`
		Text  string `json:"text"`
	}

	if err := services.ParseJSON(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := b.HandleMessage(r.Context(), req.Owner, req.Repo, req.Issue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	services.JSON(w, http.StatusOK, result)
}

func main() {
	bot := NewCodeAssistBot()
	http.HandleFunc("/code-assist", bot.HTTPHandler)
	log.Println("🤖 code-assist bot running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}