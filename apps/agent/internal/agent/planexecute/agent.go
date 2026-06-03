package planexecute

import (
	"context"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
)

const (
	defaultMaxPlanSteps       = 12
	defaultMaxStepRounds      = 4
	defaultMaxReplans         = 2
	planAndExecuteToolTimeout = 15 * time.Second
)

// PlanAndExecuteAgent implements a single-agent planner/executor/finalizer
// agent. The runtime owns orchestration: every step result returns to this agent
// before the next module is called.
type PlanAndExecuteAgent struct {
	model         llm.Model
	maxSteps      int
	maxStepRounds int
	maxReplans    int
	toolTimeout   time.Duration
}

// PlanAndExecuteOption configures optional PlanAndExecuteAgent settings.
type PlanAndExecuteOption func(*PlanAndExecuteAgent)

func NewPlanAndExecuteAgent(model llm.Model, opts ...PlanAndExecuteOption) *PlanAndExecuteAgent {
	agent := &PlanAndExecuteAgent{
		model:         model,
		maxSteps:      defaultMaxPlanSteps,
		maxStepRounds: defaultMaxStepRounds,
		maxReplans:    defaultMaxReplans,
		toolTimeout:   planAndExecuteToolTimeout,
	}

	for _, opt := range opts {
		opt(agent)
	}

	return agent
}



// WithPlanMaxSteps overrides the maximum number of planned steps.
func WithPlanMaxSteps(maxSteps int) PlanAndExecuteOption {
	return func(agent *PlanAndExecuteAgent) {
		if maxSteps > 0 {
			agent.maxSteps = maxSteps
		}
	}
}

// WithPlanToolTimeout overrides per-tool execution timeout.
func WithPlanToolTimeout(timeout time.Duration) PlanAndExecuteOption {
	return func(agent *PlanAndExecuteAgent) {
		if timeout > 0 {
			agent.toolTimeout = timeout
		}
	}
}

// Run appends the user prompt, creates a plan, executes planned steps, then
// synthesizes the final assistant answer.
func (a *PlanAndExecuteAgent) Run(ctx context.Context, input agent.RunInput) <-chan agent.Event {
	events := make(chan agent.Event)
	go a.run(ctx, events, input)
	return events
}
