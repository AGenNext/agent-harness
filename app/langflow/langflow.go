package langflow

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/harness/harness/app/services"
)

// =============================================
// LangFlow Visual Components
// =============================================

// NodeType - all LangFlow node types
type NodeType string

const (
	NodeInput    NodeType = "input"     // Green - Start
	NodeLLM     NodeType = "llm"      // Purple - AI
	NodePrompt  NodeType = "prompt"    // Blue - Template
	NodeAgent  NodeType = "agent"    // Yellow - Agent
	NodeTool   NodeType = "tool"     // Orange - Tool
	NodeChain  NodeType = "chain"    // Red - Chain
	NodeRouter NodeType = "router"   // Pink - Decision
	NodeOutput NodeType = "output"   // Red - End
)

// NodeColor returns color for node type
func NodeColor(nt NodeType) string {
	colors := map[NodeType]string{
		NodeInput:   "#22c55e", // Green
		NodeLLM:    "#a855f7", // Purple
		NodePrompt: "#3b82f6", // Blue
		NodeAgent:  "#eab308", // Yellow
		NodeTool:   "#f97316", // Orange
		NodeChain:  "#ef4444", // Red
		NodeRouter: "#ec4899", // Pink
		NodeOutput: "#ef4444", // Red
	}
	return colors[nt]
}

// =============================================
// Visual Component
// =============================================

// Component is a visual node
type Component struct {
	ID     string     `json:"id"`
	Type   NodeType   `json:"type"`
	Name   string    `json:"name"`
	Label  string    `json:"label"`
	Color  string    `json:"color"`
	Position Position `json:"position"`
	Data   json.RawMessage `json:"data"`
	Inputs []string  `json:"inputs"`
	Outputs []string `json:"outputs"`
}

// Position on canvas
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// =============================================
// Visual Flow (Canvas)
// =============================================

// Flow is a visual workflow
type Flow struct {
	ID          string      `json:"id"`
	Name        string     `json:"name"`
	Description string    `json:"description"`
	Nodes       []*Component `json:"nodes"`
	Edges       []*VisualEdge `json:"edges"`
	StartNode   string     `json:"start_node"`
	Position    Position  `json:"position"`
}

// VisualEdge connects nodes visually
type VisualEdge struct {
	ID       string `json:"id"`
	FromNode string `json:"from_node"`
	ToNode  string `json:"to_node"`
	Label   string `json:"label"`
	Type    string `json:"type"` // straight, step, smoothstep
}

// =============================================
// Prebuilt Visual Flows
// =============================================

// NewVisualFlow creates a flow with positioning
func NewVisualFlow(name string) *Flow {
	return &Flow{
		ID:     name,
		Name:   name,
		Nodes:  []*Component{},
		Edges:  []*VisualEdge{},
	}
}

// AddNode adds a visual node
func (f *Flow) AddNode(c *Component) {
	f.Nodes = append(f.Nodes, c)
}

// Connect visual edge
func (f *Flow) Connect(from, to, label string) {
	f.Edges = append(f.Edges, &VisualEdge{
		ID:       fmt.Sprintf("%s-%s", from, to),
		FromNode: from,
		ToNode:   to,
		Label:   label,
		Type:    "smoothstep",
	})
}

// =============================================
// Prebuilt Visual Flows (Drag & Drop)
// =============================================

func PrebuiltVisualFlows() map[string]*Flow {
	return map[string]*Flow{
		"fix-issue-visual": fixIssueVisual(),
		"full-cicd-visual": fullCICDVisual(),
		"security-scan-visual": securityScanVisual(),
		"review-deploy-visual": reviewDeployVisual(),
	}
}

func fixIssueVisual() *Flow {
	f := NewVisualFlow("fix-issue-visual")
	f.Name = "🔧 Fix Issue"
	f.Description = "Take issue → Generate fix → Create PR"

	// Nodes with positions
	f.Nodes = []*Component{
		{ID: "input", Type: NodeInput, Name: "Issue", Label: "📥 Input Issue", Position: Position{X: 100, Y: 200}, Color: NodeColor(NodeInput)},
		{ID: "prompt", Type: NodePrompt, Name: "Fix Template", Position: Position{X: 250, Y: 200}, Color: NodeColor(NodePrompt)},
		{ID: "llm", Type: NodeLLM, Name: "GPT-4", Position: Position{X: 400, Y: 200}, Color: NodeColor(NodeLLM)},
		{ID: "agent", Type: NodeAgent, Name: "Code Assist", Position: Position{X: 550, Y: 200}, Color: NodeColor(NodeAgent)},
		{ID: "output", Type: NodeOutput, Name: "PR Created", Position: Position{X: 700, Y: 200}, Color: NodeColor(NodeOutput)},
	}

	// Edges
	f.Edges = []*VisualEdge{
		{FromNode: "input", ToNode: "prompt", Label: "issue data"},
		{FromNode: "prompt", ToNode: "llm", Label: "prompt"},
		{FromNode: "llm", ToNode: "agent", Label: "response"},
		{FromNode: "agent", ToNode: "output", Label: "PR"},
	}

	f.StartNode = "input"
	return f
}

func fullCICDVisual() *Flow {
	f := NewVisualFlow("full-cicd-visual")
	f.Name = "🚀 Full CI/CD"
	f.Description = "Fix → Review → Test → Deploy"

	f.Nodes = []*Component{
		{ID: "issue", Type: NodeInput, Name: "New Issue", Position: Position{X: 50, Y: 150}, Color: NodeColor(NodeInput)},
		{ID: "fix", Type: NodeAgent, Name: "Code Assist", Position: Position{X: 200, Y: 150}, Color: "#eab308"},
		{ID: "review", Type: NodeAgent, Name: "Code Review", Position: Position{X: 350, Y: 50}, Color: "#eab308"},
		{ID: "test", Type: NodeAgent, Name: "Code Tester", Position: Position{X: 350, Y: 250}, Color: "#eab308"},
		{ID: "decision", Type: NodeRouter, Name: "Tests Pass?", Position: Position{X: 500, Y: 150}, Color: "#ec4899"},
		{ID: "deploy", Type: NodeAgent, Name: "Code Deploy", Position: Position{X: 650, Y: 150}, Color: "#eab308"},
		{ID: "output", Type: NodeOutput, Name: "Done", Position: Position{X: 800, Y: 150}, Color: "#ef4444"},
	}

	f.Edges = []*VisualEdge{
		{FromNode: "issue", ToNode: "fix", Label: "issue"},
		{FromNode: "fix", ToNode: "review", Label: "code"},
		{FromNode: "fix", ToNode: "test", Label: "code"},
		{FromNode: "review", ToNode: "decision", Label: "review"},
		{FromNode: "test", ToNode: "decision", Label: "tests"},
		{FromNode: "decision", ToNode: "deploy", Label: "yes"},
		{FromNode: "deploy", ToNode: "output", Label: "deployed"},
	}

	f.StartNode = "issue"
	return f
}

func securityScanVisual() *Flow {
	f := NewVisualFlow("security-scan-visual")
	f.Name = "🔒 Security Scan"
	f.Description = "Scan → Fix → Verify"

	f.Nodes = []*Component{
		{ID: "scan", Type: NodeInput, Name: "Code Input", Position: Position{X: 100, Y: 200}, Color: NodeColor(NodeInput)},
		{ID: "analyze", Type: NodeLLM, Name: "Security AI", Position: Position{X: 250, Y: 200}, Color: NodeColor(NodeLLM)},
		{ID: "findings", Type: NodeRouter, Name: "Issues Found?", Position: Position{X: 400, Y: 200}, Color: NodeColor(NodeRouter)},
		{ID: "fix", Type: NodeAgent, Name: "Code Assist", Position: Position{X: 550, Y: 200}, Color: NodeColor(NodeAgent)},
		{ID: "verify", Type: NodeAgent, Name: "Code Tester", Position: Position{X: 700, Y: 200}, Color: NodeColor(NodeAgent)},
		{ID: "output", Type: NodeOutput, Name: "Report", Position: Position{X: 850, Y: 200}, Color: NodeColor(NodeOutput)},
	}

	f.Edges = []*VisualEdge{
		{FromNode: "scan", ToNode: "analyze", Label: "code"},
		{FromNode: "analyze", ToNode: "findings", Label: "analysis"},
		{FromNode: "findings", ToNode: "fix", Label: "yes"},
		{FromNode: "findings", ToNode: "output", Label: "no"},
		{FromNode: "fix", ToNode: "verify", Label: "fixes"},
		{FromNode: "verify", ToNode: "output", Label: "verified"},
	}

	f.StartNode = "scan"
	return f
}

func reviewDeployVisual() *Flow {
	f := NewVisualFlow("review-deploy-visual")
	f.Name = "📋 Review → Deploy"
	f.Description = "Review PR and deploy"

	f.Nodes = []*Component{
		{ID: "pr", Type: NodeInput, Name: "PR Input", Position: Position{X: 100, Y: 200}, Color: NodeColor(NodeInput)},
		{ID: "review", Type: NodeAgent, Name: "Code Review", Position: Position{X: 250, Y: 200}, Color: NodeColor(NodeAgent)},
		{ID: "approve", Type: NodeRouter, Name: "Approved?", Position: Position{X: 400, Y: 200}, Color: NodeColor(NodeRouter)},
		{ID: "deploy", Type: NodeAgent, Name: "Code Deploy", Position: Position{X: 550, Y: 200}, Color: NodeColor(NodeAgent)},
		{ID: "output", Type: NodeOutput, Name: "Status", Position: Position{X: 700, Y: 200}, Color: NodeColor(NodeOutput)},
	}

	f.Edges = []*VisualEdge{
		{FromNode: "pr", ToNode: "review", Label: "PR"},
		{FromNode: "review", ToNode: "approve", Label: "review"},
		{FromNode: "approve", ToNode: "deploy", Label: "approve"},
		{FromNode: "approve", ToNode: "output", Label: "reject"},
		{FromNode: "deploy", ToNode: "output", Label: "deployed"},
	}

	f.StartNode = "pr"
	return f
}

// =============================================
// LangFlow Engine
// =============================================

type Engine struct {
	flows     map[string]*Flow
	notifier  *services.Notifier
}

func NewEngine(notifier *services.Notifier) *Engine {
	e := &Engine{
		flows: make(map[string]*Flow),
		notifier: notifier,
	}
	// Register prebuilt
	for name, f := range PrebuiltVisualFlows() {
		e.flows[name] = f
	}
	return e
}

// ListFlows lists all available flows
func (e *Engine) ListFlows() []*Flow {
	list := make([]*Flow, 0, len(e.flows))
	for _, f := range e.flows {
		list = append(list, f)
	}
	return list
}

// GetFlow retrieves a flow
func (e *Engine) GetFlow(name string) *Flow {
	return e.flows[name]
}

// Execute runs a flow
func (e *Engine) Execute(ctx context.Context, flowName string, input map[string]interface{}) (map[string]interface{}, error) {
	flow := e.flows[flowName]
	if flow == nil {
		return nil, fmt.Errorf("flow not found: %s", flowName)
	}

	// Execute nodes in order
	results := make(map[string]interface{})
	current := flow.StartNode

	for current != "" {
		node := findNode(flow, current)
		if node == nil {
			break
		}

		result, err := executeNode(ctx, node, results)
		if err != nil {
			results[current] = map[string]interface{}{"error": err.Error()}
		} else {
			results[current] = result
		}

		current = nextEdge(flow, current)
	}

	return results, nil
}

func findNode(flow *Flow, id string) *Component {
	for _, n := range flow.Nodes {
		if n.ID == id {
			return n
		}
	}
	return nil
}

func nextEdge(flow *Flow, fromID string) string {
	for _, edge := range flow.Edges {
		if edge.FromNode == fromID {
			return edge.ToNode
		}
	}
	return ""
}

func executeNode(ctx context.Context, node *Component, results map[string]interface{}) (interface{}, error) {
	switch node.Type {
	case NodeInput:
		return results["input"], nil

	case NodeAgent:
		return "agent executed", nil

	case NodeLLM:
		return "llm response", nil

	case NodeOutput:
		return results, nil

	default:
		return nil, nil
	}
}

// ToJSON exports for LangFlow visual editor
func (f *Flow) ToJSON() (string, error) {
	b, _ := json.MarshalIndent(f, "", "  ")
	return string(b), nil
}

// FromJSON imports from LangFlow
func FromJSON(data string) (*Flow, error) {
	var flow Flow
	err := json.Unmarshal([]byte(data), &flow)
	return &flow, err
}