package jira

import (
	"context"
	"encoding/json"
	"net/http"
)

type Config struct {
	URL   string
	User  string
	Token string
}

type Event struct {
	WebhookEvent string `json:"webhookEvent"`
	Issue      Issue  `json:"issue"`
}

type Issue struct {
	Key     string `json:"key"`
	Type    string `json:"type"`
	Status string `json:"status"`
}

type Client struct {
	config *Config
	agent  interface{ Execute(ctx context.Context, agent, input string) (string, error) }
}

func NewClient(cfg *Config, agent interface{ Execute(ctx context.Context, agent, input string) (string, error) }) *Client {
	return &Client{config: cfg, agent: agent}
}

func (c *Client) HandleWebhook(w http.ResponseWriter, r *http.Request) error {
	var event Event
	json.NewDecoder(r.Body).Decode(&event)

	switch event.WebhookEvent {
	case "jira:issue_created":
		// Trigger code-assist
	case "jira:issue_updated":
		// Check status
	}

	return nil
}