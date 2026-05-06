package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-github/v50/github"
)

// Notifier handles communication with humans via multiple channels
type Notifier struct {
	ghClient    *github.Client
	slackURL    string
	discordURL string
	mattermost string
	whatsApp   string
	email     string
	webhookURL string
}

// NewNotifier creates a new notifier
func NewNotifier(ghToken string, opts ...NotifierOption) *Notifier {
	n := &Notifier{
		ghClient: github.NewClient(ghToken),
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

// NotifierOption configures notifier
type NotifierOption func(*Notifier)

// WithSlack sets Slack webhook
func WithSlack(url string) NotifierOption {
	return func(n *Notifier) { n.slackURL = url }
}

// WithDiscord sets Discord webhook
func WithDiscord(url string) NotifierOption {
	return func(n *Notifier) { n.discordURL = url }
}

// WithMattermost sets Mattermost webhook
func WithMattermost(url string) NotifierOption {
	return func(n *Notifier) { n.mattermost = url }
}

// WithWhatsApp sets WhatsApp Business API
func WithWhatsApp(addr string) NotifierOption {
	return func(n *Notifier) { n.whatsApp = addr }
}

// WithEmail sets email recipients
func WithEmail(addr string) NotifierOption {
	return func(n *Notifier) { n.email = addr }
}

// WithWebhook sets generic webhook
func WithWebhook(url string) NotifierOption {
	return func(n *Notifier) { n.webhookURL = url }
}

// IssueComment comments on GitHub issue
func (n *Notifier) IssueComment(ctx context.Context, owner, repo string, issueNum int, message string) error {
	if n.ghClient == nil {
		return nil
	}
	_, _, err := n.ghClient.Issues.CreateComment(ctx, owner, repo, issueNum, &github.IssueComment{
		Body: github.String(message),
	})
	return err
}

// PRComment comments on Pull Request
func (n *Notifier) PRComment(ctx context.Context, owner, repo string, prNum int, message string) error {
	if n.ghClient == nil {
		return nil
	}
	_, _, err := n.ghClient.PullRequests.CreateComment(ctx, owner, repo, prNum, &github.PullRequestComment{
		Body: github.String(message),
	})
	return err
}

// SlackNotify sends notification to Slack
func (n *Notifier) SlackNotify(ctx context.Context, message string) error {
	if n.slackURL == "" {
		return nil
	}
	payload := map[string]interface{}{
		"text":            message,
		"allowed_mentions": map[string]string{"parse": ""},
	}
	return n.sendJSON(n.slackURL, payload)
}

// DiscordNotify sends notification to Discord
func (n *Notifier) DiscordNotify(ctx context.Context, message string) error {
	if n.discordURL == "" {
		return nil
	}
	payload := map[string]interface{}{
		"content":            message,
		"allowed_mentions": map[string][]string{"parse": {}},
	}
	return n.sendJSON(n.discordURL, payload)
}

// MattermostNotify sends notification to Mattermost
func (n *Notifier) MattermostNotify(ctx context.Context, message string) error {
	if n.mattermost == "" {
		return nil
	}
	payload := map[string]interface{}{
		"message": message,
	}
	return n.sendJSON(n.mattermost, payload)
}

// WhatsAppNotify sends notification via WhatsApp Business API
func (n *Notifier) WhatsAppNotify(ctx context.Context, message string) error {
	if n.whatsApp == "" {
		return nil
	}
	// WhatsApp Business API payload
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":               n.whatsApp,
		"type":             "text",
		"text": map[string]string{
			"body": message,
		},
	}
	return n.sendJSON(n.whatsApp, payload)
}

// EmailNotify sends email notification
func (n *Notifier) EmailNotify(ctx context.Context, subject, body string) error {
	if n.email == "" {
		return nil
	}
	// For simple SMTP, would use net/smtp
	// This is a placeholder - implement with actual SMTP server
	payload := map[string]string{
		"to":      n.email,
		"subject": subject,
		"body":    body,
	}
	jsonStr, _ := json.Marshal(payload)
	fmt.Println(string(jsonStr)) // Print for now
	return nil
}

// WebhookNotify sends to generic webhook
func (n *Notifier) WebhookNotify(ctx context.Context, data map[string]interface{}) error {
	if n.webhookURL == "" {
		return nil
	}
	return n.sendJSON(n.webhookURL, data)
}

// NotifyResult sends result notification to all channels
func (n *Notifier) NotifyResult(ctx context.Context, result *Result) error {
	msg := fmt.Sprintf("✅ %s", result.Message())

	// Send to all configured channels
	_ = n.IssueComment(ctx, result.Owner, result.Repo, result.Issue, msg)
	_ = n.SlackNotify(ctx, msg)
	_ = n.DiscordNotify(ctx, msg)
	_ = n.MattermostNotify(ctx, msg)
	_ = n.WhatsAppNotify(ctx, msg)
	_ = n.EmailNotify(ctx, "Agent Update", msg)
	_ = n.WebhookNotify(ctx, map[string]interface{}{
		"status":  result.Status,
		"message": result.Message(),
		"pr":     result.PR,
	})

	return nil
}

// sendJSON sends JSON payload to URL
func (n *Notifier) sendJSON(url string, payload interface{}) error {
	jsonStr, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status: %d", resp.StatusCode)
	}
	return nil
}