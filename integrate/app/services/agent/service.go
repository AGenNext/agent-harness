// Agent service for Harness kernel
package agent

import (
	"context"
	"time"

	"github.com/harness/gitness/app/store"
)

type Agent struct {
	Name    string   `json:"name"`
	Type   string   `json:"type"`
	Status string   `json:"status"`
	Skills []string `json:"skills"`
}

type Result struct {
	Name      string        `json:"name"`
	Status   string        `json:"status"`
	Output   string        `json:"output"`
	Duration time.Duration `json:"duration"`
}

type Service struct {
	store store.Store
}

func NewService(store store.Store) *Service {
	return &Service{store: store}
}

func (s *Service) List(ctx context.Context) ([]*Agent, error) {
	return []*Agent{
		{Name: "code-assist", Type: "agent", Skills: []string{"fix", "security", "docs"}},
		{Name: "code-review", Type: "agent", Skills: []string{"review", "security"}},
		{Name: "code-tester", Type: "agent", Skills: []string{"test", "unit"}},
		{Name: "code-deploy", Type: "agent", Skills: []string{"deploy", "rollback"}},
	}, nil
}

func (s *Service) Run(ctx context.Context, name string, input interface{}) (*Result, error) {
	return &Result{Name: name, Status: "success", Duration: time.Second}, nil
}