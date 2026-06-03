package planexecute

import (
	"context"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
)

func (a *PlanAndExecuteAgent) replan(
	ctx context.Context,
	events chan<- agent.Event,
	input agent.RunInput,
	plan executionPlan,
	results []stepResult,
	state executionState,
	stats *agent.RunStats,
) (executionPlan, error) {
	messages := input.PromptSnapshot.BuildMessages(nil, []conversation.Message{
		conversation.NewUserMessage(replannerPrompt(input.UserPrompt, plan, results, state)),
	})
	response, err := a.callModel(ctx, events, input.SessionID, stats.Rounds, messages, nil, stats)
	if err != nil {
		return executionPlan{}, err
	}

	var updated executionPlan
	if err := decodeJSONText(conversation.Text(response.Message), &updated); err != nil {
		return executionPlan{}, fmt.Errorf("decode replanner response: %w", err)
	}
	if err := validatePlan(updated); err != nil {
		return executionPlan{}, err
	}
	return updated, nil
}
