package main

// =============================================
// Integration Guide: Adding Agents to Harness Kernel
// =============================================

/*
STEP 1: Add imports to cmd/gitness/wire.go
--------------------------------------

import (
    agentcontroller "github.com/harness/gitness/app/api/controller/agent"
    agentservice "github.com/harness/gitness/app/services/agent"
    "github.com/harness/gitness/app/services/agent/runtime"
)

STEP 2: Add to Initialize() function
--------------------------------

func Initialize(ib *Injector) error {
    // Existing code...
    
    // Add agent runtime
    ib.Register(&runtime.Manager{})

    // Add agent service
    agent := &agentservice.AgentService{
        Store:    ib.Store(),
        Runtime: ib.AgentRuntime(),
    }
    ib.Register(agent)

    // Add agent controller
    ib.Register(&agentcontroller.AgentController{
        AgentService: agent,
    })

    return nil
}

STEP 3: Add routes in router
-------------------------

// In app/router/router.go, add:

r.Get("/api/v1/agents", agentCtrl.ListAgents)
r.Get("/api/v1/agents/:name", agentCtrl.GetAgent)
r.Post("/api/v1/agents/:name/run", agentCtrl.RunAgent)
r.Get("/api/v1/runtimes", runtimeCtrl.List)

STEP 4: Environment variables
----------------------------

HARNESS_AGENTS_ENABLED=true
*/