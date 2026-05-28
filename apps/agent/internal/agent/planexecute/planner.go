package planexecute

import (
	"context"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
)

func (a *PlanAndExecuteAgent) createPlan(
	ctx context.Context,
	events chan<- agent.Event,
	session *session.State,
	userPrompt string,
	stats *agent.RunStats,
) (executionPlan, error) {
	messages := session.BuildMessages(nil, []conversation.Message{
		conversation.NewUserMessage(plannerPrompt(userPrompt, session.ToolDefinitions())),
	})
	response, err := a.callModel(ctx, events, session.ID, 0, messages, nil, stats)
	if err != nil {
		return executionPlan{}, err
	}

	var plan executionPlan
	if err := decodeJSONText(conversation.Text(response.Message), &plan); err != nil {
		return executionPlan{}, fmt.Errorf("decode planner response: %w", err)
	}
	if err := validatePlan(plan); err != nil {
		return executionPlan{}, err
	}
	return plan, nil
}
