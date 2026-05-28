package app

import (
	"fmt"
	"log/slog"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent/planexecute"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent/react"
	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
)

// newAgent creates the configured agent loop implementation.
func newAgent(
	cfg config.Config,
	provider llm.Provider,
	logger *slog.Logger,
) (agent.Agent, error) {
	switch cfg.Agent.Type {
	case "", agent.TypeReAct:
		return react.NewReActAgent(
			provider,
			react.WithLogger(logger),
			react.WithMaxRounds(cfg.Agent.MaxRounds),
			react.WithToolTimeout(cfg.Agent.ToolTimeout),
		), nil
	case agent.TypePlanAndExecute:
		return planexecute.NewPlanAndExecuteAgent(
			provider,
			planexecute.WithPlanLogger(logger),
			planexecute.WithPlanMaxSteps(cfg.Agent.MaxRounds),
			planexecute.WithPlanToolTimeout(cfg.Agent.ToolTimeout),
		), nil
	default:
		return nil, fmt.Errorf("unsupported agent type %q", cfg.Agent.Type)
	}
}
