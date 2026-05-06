package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Unified listener - one tool for GitHub, Slack, Email, Jira, Linear, Discord
type Event struct {
	Source string `json:"source"`
	Type   string `json:"type"`
	Input  map[string]interface{} `json:"input"`
}

type Handler struct {
	routes map[string]func(ctx context.Context, e *Event) (string, error)
}

func New() *Handler {
	h := &Handler{routes: make(map[string]func(ctx context.Context, e *Event) (string, error))}
	h.registerAll()
	return h
}

func (h *Handler) registerAll() {
	h.routes["github.issue"] = func(ctx context.Context, e *Event) (string, error) { return "code-assist:fix", nil }
	h.routes["github.pr"] = func(ctx context.Context, e *Event) (string, error) { return "code-review:review", nil }
	h.routes["slack.command"] = func(ctx context.Context, e *Event) (string, error) { return "process command", nil }
	h.routes["email"] = func(ctx context.Context, e *Event) (string, error) { return "code-assist:fix", nil }
	h.routes["jira"] = func(ctx context.Context, e *Event) (string, error) { return "process issue", nil }
	h.routes["linear"] = func(ctx context.Context, e *Event) (string, error) { return "process issue", nil }
	h.routes["discord"] = func(ctx context.Context, e *Event) (string, error) { return "process message", nil }
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) error {
	source := r.URL.Query().Get("source")
	var input map[string]interface{}
	json.NewDecoder(r.Body).Decode(&input)

	fn := h.routes[source]
	if fn == nil {
		fn = h.routes["github.issue"]
	}

	result, _ := fn(r.Context(), &Event{Source: source, Input: input})
	return json.NewEncoder(w).Encode(map[string]string{"result": result})
}

/*
# AnyListener - One Tool for All Triggers

POST /listeners/handle?source=github
POST /listeners/handle?source=slack  
POST /listeners/handle?source=email
POST /listeners/handle?source=jira

# Slack: /agent fix, /agent review
# Email: Subject with "fix" triggers code-assist
# GitHub: Issue opened → code-assist
# Jira/Linear: Issue created → task
*/