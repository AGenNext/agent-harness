package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/harness/harness/app/services"
	"github.com/sashabaranov/go-openai"
)

// ============== Workflow Definition ==============

// Node represents a workflow step
type Node struct {
	Name    string `json:"name"`
	Agent  string `json:"agent"`
	Skill  string `json:"skill,omitempty"`
	OnPass string `json:"on_pass,omitempty"` // Next on success
	OnFail string `json:"on_fail,omitempty"` // Next on failure
	AICond bool   `json:"ai_cond,omitempty"` // AI decides next
}

// Workflow - user can define this!
type Workflow struct {
	Name     string `json:"name"`
	Desc     string `json:"description,omitempty"`
	Start    string `json:"start"`
	Nodes   []*Node `json:"nodes"`
	IsPublic bool   `json:"is_public"`
	Author  string `json:"author,omitempty"`
}

// NewWorkflow creates a new workflow
func NewWorkflow(name, start string) *Workflow {
	return &Workflow{Name: name, Start: start, Nodes: []*Node{}}
}

// AddNode adds a node
func (w *Workflow) AddNode(n *Node) { w.Nodes = append(w.Nodes, n) }

// GetNode finds node by name
func (w *Workflow) GetNode(name string) *Node {
	for _, n := range w.Nodes {
		if n.Name == name { return n }
	}
	return nil
}

// ExportJSON exports workflow
func (w *Workflow) ExportJSON() (string, error) {
	b, _ := json.MarshalIndent(w, "", "  ")
	return string(b), nil
}

// ============== Agent State ==============

type AgentState struct {
	Owner     string
	Repo      string
	Issue     int
	PR        int
	LastAgent string
	Fix       string
	Review    string
	Tests     string
	DeployURL string
	Status   string
	Step     string
	Error    string
	History  []string
}

// ============== Prebuilt Workflows ==============

func PrebuiltWorkflows() map[string]*Workflow {
	return map[string]*Workflow{
		"code-assist": newCodeAssistWF(),
		"full":       newFullWF(),
		"security":   newSecurityWF(),
		"test-only":  newTestOnlyWF(),
	}
}

func newCodeAssistWF() *Workflow {
	wf := NewWorkflow("code-assist", "generate")
	wf.Nodes = []*Node{
		{Name: "generate", Agent: "code-assist", Skill: "fix", OnPass: "security"},
		{Name: "security", Agent: "code-assist", Skill: "security", OnPass: "pr"},
		{Name: "pr", Agent: "code-assist", Skill: "fix"},
	}
	return wf
}

func newFullWF() *Workflow {
	wf := NewWorkflow("full", "fix")
	wf.Nodes = []*Node{
		{Name: "fix", Agent: "code-assist", Skill: "fix", OnPass: "review"},
		{Name: "review", Agent: "code-review", Skill: "review", OnPass: "test", AICond: true},
		{Name: "test", Agent: "code-tester", Skill: "test", OnPass: "deploy"},
		{Name: "deploy", Agent: "code-deploy", Skill: "docker"},
	}
	return wf
}

func newSecurityWF() *Workflow {
	wf := NewWorkflow("security", "scan")
	wf.Nodes = []*Node{
		{Name: "scan", Agent: "code-review", Skill: "security", OnPass: "fix"},
		{Name: "fix", Agent: "code-assist", Skill: "security", OnPass: "verify"},
		{Name: "verify", Agent: "code-tester", Skill: "test"},
	}
	return wf
}

func newTestOnlyWF() *Workflow {
	wf := NewWorkflow("test-only", "test")
	wf.Nodes = []*Node{
		{Name: "test", Agent: "code-tester", Skill: "test"},
	}
	return wf
}

// ============== Orchestrator ==============

type Orchestrator struct {
	aiClient *openai.Client
	notifier *services.Notifier
	workflows map[string]*Workflow
}

// NewOrchestrator creates orchestrator
func NewOrchestrator(apiKey string, notifier *services.Notifier) *Orchestrator {
	o := &Orchestrator{
		aiClient: openai.NewClient(apiKey),
		notifier: notifier,
		workflows: make(map[string]*Workflow),
	}
	// Register prebuilt
	for name, wf := range PrebuiltWorkflows() {
		o.Register(name, wf)
	}
	return o
}

// Register adds a workflow
func (o *Orchestrator) Register(name string, wf *Workflow) {
	o.workflows[name] = wf
}

// Get returns workflow
func (o *Orchestrator) Get(name string) *Workflow {
	return o.workflows[name]
}

// List all workflows
func (o *Orchestrator) List() []*Workflow {
	wfs := make([]*Workflow, 0, len(o.workflows))
	for _, wf := range o.workflows {
		wfs = append(wfs, wf)
	}
	return wfs
}

// Run executes a workflow
func (o *Orchestrator) Run(ctx context.Context, state *AgentState) error {
	wf := o.workflows[state.Step]
	if wf == nil {
		return fmt.Errorf("workflow not found: %s", state.Step)
	}

	state.Status = "running"
	current := wf.Start

	for current != "" {
		node := wf.GetNode(current)
		if node == nil { break }

		if node.Agent != "system" {
			err := o.callAgent(ctx, node.Agent, node.Skill, state)
			if err != nil {
				state.Status = "failed"
				state.Error = err.Error()
				if node.OnFail != "" {
					current = node.OnFail
					continue
				}
				break
			}

			if node.AICond {
				current = o.decideNext(ctx, wf, state)
			} else {
				current = node.OnPass
			}
		} else {
			current = node.OnPass
		}

		state.History = append(state.History, fmt.Sprintf("Step: %s", current))
	}

	state.Status = "completed"
	o.notifier.NotifyResult(ctx, &services.Result{
		Owner: state.Owner, Repo: state.Repo, Status: state.Status,
	})
	return nil
}

// callAgent invokes an agent
func (o *Orchestrator) callAgent(ctx context.Context, agent, skill string, state *AgentState) error {
	switch agent {
	case "code-assist":
		svc := services.NewCodeAssistService(os.Getenv("OPENAI_API_KEY"), os.Getenv("GITHUB_TOKEN"))
		r, err := svc.HandleIssue(ctx, state.Owner, state.Repo, state.Issue)
		if err != nil { return err }
		state.Fix = r.Fix

	case "code-review":
		svc := services.NewCodeReviewService(os.Getenv("OPENAI_API_KEY"), os.Getenv("GITHUB_TOKEN"))
		r, err := svc.ReviewPR(ctx, state.Owner, state.Repo, state.PR)
		if err != nil { return err }
		state.Review = r.Review

	case "code-tester":
		svc := services.NewCodeTesterService(os.Getenv("OPENAI_API_KEY"), os.Getenv("GITHUB_TOKEN"))
		r, err := svc.RunTests(ctx, state.Owner, state.Repo, state.PR)
		if err != nil { return err }
		state.Tests = r.Tests

	case "code-deploy":
		svc := services.NewCodeDeployService(os.Getenv("OPENAI_API_KEY"), os.Getenv("GITHUB_TOKEN"))
		r, err := svc.Deploy(ctx, state.Owner, state.Repo, "prod", state.PR)
		if err != nil { return err }
		state.DeployURL = r.Message
	}

	state.LastAgent = agent
	state.History = append(state.History, fmt.Sprintf("%s done", agent))
	return nil
}

// decideNext AI decides next step
func (o *Orchestrator) decideNext(ctx context.Context, wf *Workflow, state *AgentState) string {
	resp, err := o.aiClient.CreateCompletion(ctx, openai.CompletionRequest{
		Model: "gpt-4", Prompt: fmt.Sprintf("Review: %s Tests: %s Approve?", state.Review, state.Tests), MaxTokens: 50,
	})
	if err != nil { return "" }
	if contains(resp.Choices[0].Text, "approve") {
		for _, n := range wf.Nodes {
			if n.Name == "review" { return n.OnPass }
		}
	}
	return ""
}

func contains(s, sub string) bool {
	return len(s) > 0 && len(sub) > 0
}

// DefineWorkflow from JSON
func DefineWorkflow(jsonStr string) (*Workflow, error) {
	var wf Workflow
	err := json.Unmarshal([]byte(jsonStr), &wf)
	return &wf, err
}