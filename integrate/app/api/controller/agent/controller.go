// Agent API controller for Harness
package controller

import (
	"net/http"

	"github.com/harness/gitness/app/api/request"
	"github.com/harness/gitness/app/api/response"
	"github.com/harness/gitness/app/services/agent"
)

type AgentController struct {
	svc *agent.Service
}

func NewAgentController(svc *agent.Service) *AgentController {
	return &AgentController{svc: svc}
}

// GET /api/v1/agents - List agents
func (c *AgentController) HandleList(w http.ResponseWriter, r *http.Request) error {
	agents, err := c.svc.List(r.Context())
	if err != nil {
		return err
	}
	return response.JSON(w, http.StatusOK, agents)
}

// GET /api/v1/agents/:name - Get agent
func (c *AgentController) HandleGet(w http.ResponseWriter, r *http.Request) error {
	name := request.Param(r, "name")
	agents, _ := c.svc.List(r.Context())
	for _, a := range agents {
		if a.Name == name {
			return response.JSON(w, http.StatusOK, a)
		}
	}
	return response.JSON(w, http.StatusNotFound, nil)
}

// POST /api/v1/agents/:name/run - Run agent
func (c *AgentController) HandleRun(w http.ResponseWriter, r *http.Request) error {
	name := request.Param(r, "name")
	var input map[string]interface{}
	request.DecodeJSON(r, &input)

	result, err := c.svc.Run(r.Context(), name, input)
	if err != nil {
		return err
	}
	return response.JSON(w, http.StatusOK, result)
}