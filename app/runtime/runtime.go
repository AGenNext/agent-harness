package runtime

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// =============================================
// Runtime Types
// =============================================

type RuntimeType string

const (
	RuntimeDocker   RuntimeType = "docker"
	RuntimeK8s     RuntimeType = "k8s"
	RuntimeLambda  RuntimeType = "lambda"
	RuntimeVM      RuntimeType = "vm"
)

// Runtime config per agent
type Config struct {
	Type    RuntimeType `json:"type"`    // docker, k8s, lambda, vm
	Image  string    `json:"image"`   // docker image
	Memory string    `json:"memory"`  // 1Gi, 2Gi
	Cpu    string    `json:"cpu"`    // 1000m
	Timeout int       `json:"timeout"` // seconds
}

// Agent runtime mapping
var AgentRuntimes = map[string]*Config{
	"code-assist": {
		Type: RuntimeDocker, Image: "openautonomyx/code-assist:latest",
		Memory: "1Gi", Cpu: "1000m", Timeout: 300,
	},
	"code-review": {
		Type: RuntimeDocker, Image: "openautonomyx/code-review:latest",
		Memory: "2Gi", Cpu: "2000m", Timeout: 600,
	},
	"code-tester": {
		Type: RuntimeDocker, Image: "openautonomyx/code-tester:latest",
		Memory: "4Gi", Cpu: "4000m", Timeout: 900,
	},
	"code-deploy": {
		Type: RuntimeK8s, Image: "openautonomyx/code-deploy:latest",
		Memory: "2Gi", Cpu: "2000m", Timeout: 600,
	},
}

// =============================================
// Runtime Manager
// =============================================

type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

// ExecuteAgent runs agent in its runtime
func (m *Manager) ExecuteAgent(ctx context.Context, agent, input string) (string, error) {
	config, ok := AgentRuntimes[agent]
	if !ok {
		return "", fmt.Errorf("unknown agent: %s", agent)
	}

	switch config.Type {
	case RuntimeDocker:
		return m.runDocker(ctx, config, input)
	case RuntimeK8s:
		return m.runK8s(ctx, config, input)
	case RuntimeLambda:
		return m.runLambda(ctx, config, input)
	default:
		return m.runDocker(ctx, config, input)
	}
}

// run Docker
func (m *Manager) runDocker(ctx context.Context, cfg *Config, input string) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"--network", "host",
		"-e", "GITHUB_TOKEN="+os.Getenv("GITHUB_TOKEN"),
		"-e", "OPENAI_API_KEY="+os.Getenv("OPENAI_API_KEY"),
		"--memory", cfg.Memory,
		"--cpu-count", cfg.Cpu,
		cfg.Image, input,
	)
	return cmd.Output()
}

// run Kubernetes
func (m *Manager) runK8s(ctx context.Context, cfg *Config, input string) (string, error) {
	jobName := fmt.Sprintf("agent-%d", time.Now().UnixNano())
	cmd := exec.CommandContext(ctx, "kubectl", "run", jobName,
		"--image="+cfg.Image,
		"--env=GITHUB_TOKEN="+os.Getenv("GITHUB_TOKEN"),
		"--restart=Never",
		"--", "sh", "-c", input,
	)
	_, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Wait for completion
	cmd = exec.CommandContext(ctx, "kubectl", "wait", "--for=condition=complete", "job/"+jobName, "--timeout="+fmt.Sprintf("%ds", cfg.Timeout))
	_, err = cmd.Output()
	return "", err
}

// run Lambda
func (m *Manager) runLambda(ctx context.Context, cfg *Config, input string) (string, error) {
	function := os.Getenv("LAMBDA_FUNCTION_" + agent)
	cmd := exec.CommandContext(ctx, "aws", "lambda", "invoke",
		"--function-name", function,
		"--payload", fmt.Sprintf(`{"input":"%s"}`, input),
		"/tmp/out.json",
	)
	return cmd.Output()
}

// ListAgentRuntimes shows runtime per agent
func (m *Manager) ListAgentRuntimes() map[string]*Config {
	return AgentRuntimes
}

// UpdateRuntime changes agent runtime
func (m *Manager) UpdateRuntime(agent string, cfg *Config) {
	AgentRuntimes[agent] = cfg
}

// =============================================
// API
// =============================================

/*
# List runtimes
GET /api/runtimes

# Execute agent
POST /api/runtimes/execute
{
  "agent": "code-assist",
  "input": "fix issue 123"
}

# Update runtime
PUT /api/runtimes/code-assist
{
  "type": "docker",
  "image": "custom-image:latest",
  "memory": "2Gi"
}
*/