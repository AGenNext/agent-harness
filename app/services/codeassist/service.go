package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v50/github"
	"github.com/harness/harness/app/models"
	"github.com/harness/harness/app/store"
	"github.com/sashabaranov/go-openai"
)

// CodeAssistService generates code fixes from issues using AI
type CodeAssistService struct {
	aiClient   *openai.Client
	ghClient   *github.Client
	notifier  *Notifier
}

// NewCodeAssistService creates a new code-assist service
func NewCodeAssistService(apiKey string, ghToken string) *CodeAssistService {
	return &CodeAssistService{
		aiClient: openai.NewClient(apiKey),
		ghClient: github.NewClient(ghToken),
		notifier: NewNotifier(ghToken),
	}
}

// HandleIssue generates code fix for an issue
func (s *CodeAssistService) HandleIssue(ctx context.Context, owner, repo string, issueNum int) (*Result, error) {
	// Get issue from GitHub
	issue, _, err := s.ghClient.Issues.Get(ctx, owner, repo, issueNum)
	if err != nil {
		return nil, err
	}

	// Generate fix using AI
	resp, err := s.aiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatMessage{
			{Role: "system", Content: "You are a code assistant. Generate code fixes for GitHub issues. Only respond with code, no explanations."},
			{Role: "user", Content: fmt.Sprintf("Fix this issue:\n\nTitle: %s\n\nBody: %s", issue.GetTitle(), issue.GetBody())},
		},
		MaxTokens: 3000,
		Temperature: 0.7,
	})
	if err != nil {
		return nil, err
	}

	fix := resp.Choices[0].Message.Content
	branchName := fmt.Sprintf("fix/issue-%d", issueNum)

	// Create branch
	_, _, err = s.ghClient.Repositories.CreateBranch(ctx, owner, repo, branchName, "main")
	if err != nil {
		return nil, err
	}

	// Create PR
	pr, _, err := s.ghClient.PullRequests.Create(ctx, owner, repo, &github.NewPullRequest{
		Title:       github.String(fmt.Sprintf("Fix: %s", issue.GetTitle())),
		Head:        github.String(branchName),
		Base:        github.String("main"),
		Body:        github.String(fix),
		Draft:       github.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	// Notify human via GitHub comment
	s.notifier.IssueComment(ctx, owner, repo, issueNum, fmt.Sprintf(
		"🔧 I've created PR #%d to fix this issue. Please review!",
		*pr.Number,
	))

	return &Result{
		Fix:      fix,
		PRNumber:  *pr.Number,
		Status:   "ready",
	}, nil
}

// Handle processes HTTP requests
func (s *CodeAssistService) Handle(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return s.post(w, r)
	case "GET":
		return s.get(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}
}

func (s *CodeAssistService) post(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Owner string `json:"owner"`
		Repo  string `json:"repo"`
		Issue int    `json:"issue"`
	}
	if err := parseJSON(r, &req); err != nil {
		return err
	}

	// Get issue from store
	issue, err := store.GetIssue(r.Context(), req.Owner, req.Repo, req.Issue)
	if err != nil {
		return err
	}

	// Generate fix
	fix, err := s.HandleIssue(r.Context(), issue)
	if err != nil {
		return err
	}

	return jsonResponse(w, map[string]string{
		"fix":    fix,
		"status": "created",
	})
}

func (s *CodeAssistService) get(w http.ResponseWriter, r *http.Request) error {
	return jsonResponse(w, map[string]string{
		"agent":  "code-assist",
		"status": "ready",
	})
}