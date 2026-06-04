package planexecute

import (
	"context"
	"errors"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
)

func (a *PlanAndExecuteAgent) finalize(
	ctx context.Context,
	events chan<- agent.Event,
	input agent.RunInput,
	plan executionPlan,
	results []stepResult,
	state executionState,
	stats *agent.RunStats,
) (conversation.AssistantMessage, error) {
	messages := input.PromptSnapshot.BuildMessages(nil, []conversation.Message{
		conversation.NewUserMessage(finalizerPrompt(input.UserPrompt, plan, results, state)),
	})
	response, err := a.callModel(ctx, events, input.SessionID, stats.Rounds, messages, nil, stats)
	if err != nil {
		return conversation.AssistantMessage{}, err
	}

	answer := strings.TrimSpace(conversation.Text(response.Message))
	if answer == "" {
		return conversation.AssistantMessage{}, errors.New("finalizer returned empty answer")
	}
	return conversation.NewAssistantMessage(answer, nil), nil
}
