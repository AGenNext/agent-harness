package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"
)

// =============================================
// Unified Agent Capability (all agents have same intelligence)
// =============================================

// Capability level - all agents use the same AI brain
type Capability struct {
	Name         string   `json:"name"`
	Model       string   `json:"model"`        // GPT-4, Claude, etc
	Temperature float64 `json:"temperature"` // 0.7
	MaxTokens   int     `json:"max_tokens"`  // 2000
	Skills     []string `json:"skills"`
	Tools      []string `json:"tools"`
}

// DefaultCapability - all agents get the same AI
var DefaultCapability = &Capability{
	Name:         "unified",
	Model:       "gpt-4",
	Temperature: 0.7,
	MaxTokens:   4000,
	Skills: []string{
		"reason",      // Chain of thought
		"code",       // Code generation
		"analyze",    // Analysis
		"plan",       // Planning
		"execute",    // Execution
		"reflect",    // Self-correction
	},
	Tools: []string{
		"bash",     // Shell commands
		"python",   // Python code
		"search",   // Web search
		"memory",   // Persistent memory
	},
}

// =============================================
// Universal Agent (all agents use same brain)
// =============================================

type UniversalAgent struct {
	Name         string
	Capabilities []*Capability
	LLM          *openai.Client
}

// NewUniversalAgent - any agent can have full capabilities
func NewUniversalAgent(name string) *UniversalAgent {
	return &UniversalAgent{
		Name: name,
		Capabilities: []*Capability{
			DefaultCapability,
			{CodeGenCapability()},
			{CodeReviewCapability()},
			{TestCapability()},
			{DeployCapability()},
		},
	}
}

// Full Reasoning Capability
func CodeGenCapability() Capability {
	return Capability{
		Name: "code-generation",
		Model: "gpt-4",
		Skills: []string{
			"understand", "reason", "plan", 
			"generate", "fix", "refactor", "optimize",
		},
		Tools: []string{"bash", "python", "editor"},
	}
}

// Full Review Capability  
func CodeReviewCapability() Capability {
	return Capability{
		Name: "code-review",
		Model: "gpt-4",
		Skills: []string{
			"understand", "reason", "analyze",
			"critique", "security", "performance",
		},
		Tools: []string{"analyzer", "scanner", "linter"},
	}
}

// Full Test Capability
func TestCapability() Capability {
	return Capability{
		Name: "testing",
		Model: "gpt-4",
		Skills: []string{
			"understand", "reason", "generate",
			"unit-test", "integration-test", "e2e-test",
		},
		Tools: []string{"pytest", "test-runner", "mock"},
	}
}

// Full Deploy Capability
func DeployCapability() Capability {
	return Capability{
		Name: "deployment",
		Model: "gpt-4",
		Skills: []string{
			"understand", "reason", "plan",
			"docker", "k8s", "terraform",
		},
		Tools: []string{"docker", "kubectl", "helm"},
	}
}

// =============================================
// Execute with Full Intelligence
// =============================================

func (a *UniversalAgent) Execute(ctx context.Context, task string) (string, error) {
	// Use full GPT-4 capabilities
	resp, err := a.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: a.Capabilities[0].Model,
		Messages: []openai.ChatMessage{
			{Role: "system", Content: fmt.Sprintf(`
You are a %s agent with full intelligence.
Capabilities: %v
Tools: %v

Use chain-of-thought reasoning, self-correction, and planning.
Always explain your reasoning before executing.
`, a.Name, a.Capabilities[0].Skills, a.Capabilities[0].Tools)},
			{Role: "user", Content: task},
		},
		MaxTokens: a.Capabilities[0].MaxTokens,
		Temperature: a.Capabilities[0].Temperature,
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

// =============================================
// All Agents Have Same Capability Set
// =============================================

// AgentDefinitions - all 4 agents are equal in intelligence
var AgentDefinitions = map[string]struct {
	Name      string
	Greeting string
	Skills   []string
}{
	"code-assist": {
		Name:      "Code Assist",
		Greeting: "I'm here to help with code! 🛠️",
		Skills:   []string{"fix", "generate", "docs", "refactor", "security"},
	},
	"code-review": {
		Name:      "Code Review",
		Greeting: "I'll review your code carefully! 🔍",
		Skills:   []string{"review", "analyze", "security", "performance", "suggest"},
	},
	"code-tester": {
		Name:      "Code Tester", 
		Greeting: "I'll test your code thoroughly! 🧪",
		Skills:   []string{"unit", "integration", "e2e", "fuzz", "property"},
	},
	"code-deploy": {
		Name:      "Code Deploy",
		Greeting: "I'll deploy your code safely! 🚀",
		Skills:   []string{"docker", "k8s", "rollback", "canary"},
	},
}

// EachAgentHasFullIntelligence returns true - all agents equal
func EachAgentHasFullIntelligence() bool {
	return true // All agents get GPT-4 with full capabilities
}