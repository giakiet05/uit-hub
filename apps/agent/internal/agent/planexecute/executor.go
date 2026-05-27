package planexecute

import (
	"context"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent/loop"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
)

func (a *PlanAndExecuteAgent) executeStep(
	ctx context.Context,
	events chan<- agent.Event,
	session *runtime.Session,
	userPrompt string,
	step planStep,
	state executionState,
	stats *agent.RunStats,
) (stepResult, error) {
	stepMessages := a.promptBuilder.BuildAgentMessages(nil, []conversation.Message{
		conversation.NewUserMessage(executorPrompt(userPrompt, step, state)),
	})
	toolCalls := []string{}
	observations := []string{}

	for round := 1; round <= a.maxStepRounds; round++ {
		response, err := a.callModel(ctx, events, session.ID, stats.Rounds, stepMessages, a.stepToolDefinitions(step), stats)
		if err != nil {
			return stepResult{}, err
		}

		stepMessages = append(stepMessages, response.Message)
		calls := loop.AssistantToolCalls(response.Message)
		if len(calls) == 0 {
			result, err := decodeStepResult(step.ID, conversation.Text(response.Message), toolCalls, observations)
			if err != nil {
				return stepResult{}, err
			}
			return result, nil
		}

		for _, call := range calls {
			result, err := a.executeToolCall(ctx, events, session.ID, stats.Rounds, call, stats)
			if err != nil {
				observations = append(observations, loop.ToolErrorObservation(err))
				stepMessages = append(stepMessages, conversation.NewToolResultMessage(call.ID, loop.ToolErrorObservation(err)))
				continue
			}
			toolCalls = append(toolCalls, call.Name)
			observations = append(observations, result.Content)
			stepMessages = append(stepMessages, conversation.NewToolResultMessage(result.CallID, result.Content))
		}
	}

	return stepResult{}, fmt.Errorf("step %q exceeded max executor rounds", step.ID)
}
