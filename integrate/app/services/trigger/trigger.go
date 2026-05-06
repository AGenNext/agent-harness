package trigger

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// =============================================
// External Alert → Agent Trigger
// =============================================

type TriggerEvent struct {
	Source    string `json:"source"`
	EventType string `json:"event_type"`
	Agent    string `json:"agent"`
	Action   string `json:"action"`
	Input    map[string]interface{} `json:"input"`
}

type AgentExecutor interface {
	Execute(ctx context.Context, agent, input string) (string, error)
}

type Handler struct {
	rules    map[string]*Rule
	executor AgentExecutor
}

type Rule struct {
	Source     string `json:"source"`
	EventType string `json:"event_type"`
	Agent     string `json:"agent"`
	Action    string `json:"action"`
}

func NewHandler(ex AgentExecutor) *Handler {
	h := &Handler{executor: ex, rules: make(map[string]*Rule)}
	h.registerDefaults()
	return h
}

func (h *Handler) registerDefaults() {
	h.rules["github.issue.opened"] = &Rule{Source: "github", EventType: "issue.opened", Agent: "code-assist", Action: "fix"}
	h.rules["github.pr.opened"] = &Rule{Source: "github", EventType: "pr.opened", Agent: "code-review", Action: "review"}
	h.rules["datadog.alert"] = &Rule{Source: "datadog", EventType: "alert", Agent: "code-assist", Action: "fix"}
	h.rules["sentry.error"] = &Rule{Source: "sentry", EventType: "error", Agent: "code-assist", Action: "fix"}
	h.rules["pagerduty.incident"] = &Rule{Source: "pagerduty", EventType: "incident", Agent: "code-deploy", Action: "rollback"}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) error {
	source := r.URL.Query().Get("source")
	eventType := r.URL.Query().Get("event")

	key := source + "." + eventType
	rule, ok := h.rules[key]
	if !ok {
		return fmt.Errorf("no rule: %s", key)
	}

	var input map[string]interface{}
	json.NewDecoder(r.Body).Decode(&input)

	result, _ := h.executor.Execute(context.Background(), rule.Agent, fmt.Sprintf("%v", input))

	return json.NewEncoder(w).Encode(map[string]string{
		"agent": rule.Agent, "result": result,
	})
}

/*
# Endpoints:
POST /triggers/github?source=github&event=issue.opened
POST /triggers/datadog?source=datadog&event=alert
POST /triggers/sentry?source=sentry&event=error

# Default triggers:
github.issue.opened → code-assist fix
github.pr.opened → code-review review
datadog.alert → code-assist fix
sentry.error → code-assist fix
pagerduty.incident → code-deploy rollback
*/