package planexecute

import (
	"context"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
)

func (a *PlanAndExecuteAgent) executeStep(
	ctx context.Context,
	events chan<- agent.Event,
	input agent.RunInput,
	step planStep,
	state executionState,
	stats *agent.RunStats,
) (stepResult, error) {
	stepMessages := input.PromptSnapshot.BuildMessages(nil, []conversation.Message{
		conversation.NewUserMessage(executorPrompt(input.UserPrompt, step, state)),
	})
	toolCalls := []string{}
	observations := []string{}

	for round := 1; round <= a.maxStepRounds; round++ {
		response, err := a.callModel(ctx, events, input.SessionID, stats.Rounds, stepMessages, a.stepToolDefinitions(input, step), stats)
		if err != nil {
			return stepResult{}, err
		}

		stepMessages = append(stepMessages, response.Message)
		calls := agent.AssistantToolCalls(response.Message)
		if len(calls) == 0 {
			result, err := decodeStepResult(step.ID, conversation.Text(response.Message), toolCalls, observations)
			if err != nil {
				return stepResult{}, err
			}
			return result, nil
		}

		results, ok := agent.RunToolBatch(ctx, "plan_execute", events, input, stats.Rounds, calls, a.toolTimeout, stats)
		if !ok {
			return stepResult{}, ctx.Err()
		}

		for _, result := range results {
			toolCalls = append(toolCalls, result.Name)
			observations = append(observations, result.Content)
			stepMessages = append(stepMessages, conversation.NewToolResultMessage(result.CallID, result.Content))
		}
		if ctx.Err() != nil {
			return stepResult{}, ctx.Err()
		}
	}

	return stepResult{}, fmt.Errorf("step %q exceeded max executor rounds", step.ID)
}
