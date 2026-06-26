package orchestrator

import (
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent/react"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
)

// NewOrchestratorAgent creates a new Orchestrator agent.
// Currently it is built on top of the ReAct agent.
func NewOrchestratorAgent(model llm.Model, opts ...react.ReActOption) agent.Agent {
	return react.NewReActAgent(model, opts...)
}
