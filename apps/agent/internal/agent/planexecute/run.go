package planexecute

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
)

func (a *PlanAndExecuteAgent) run(ctx context.Context, events chan<- agent.Event, input agent.RunInput) {
	defer close(events)

	runState := agent.NewRunState(input.SessionID, a.maxSteps)
	stats := runState.Stats
	complete := func(reason agent.TerminalReason) {
		runState.Finish()
		agent.Emit(ctx, events, agent.RunCompletedEvent{
			SessionID: input.SessionID,
			Reason:    reason,
			Stats:     *stats,
		})
	}

	fail := func(err error) {
		runState.Finish()
		agent.Emit(ctx, events, agent.RunFailedEvent{
			SessionID: input.SessionID,
			Err:       err,
			Stats:     *stats,
		})
	}

	if !agent.Emit(ctx, events, agent.RunStartedEvent{SessionID: input.SessionID, MaxRounds: runState.MaxRounds}) {
		return
	}
	input.Conversation.Append(conversation.NewUserMessage(input.UserPrompt))


	plan, err := a.createPlan(ctx, events, input, stats)
	if err != nil {
		fail(err)
		return
	}
	if len(plan.Steps) > a.maxSteps {
		plan.Steps = plan.Steps[:a.maxSteps]
	}
	if !agent.Emit(ctx, events, agent.PlanCreatedEvent{
		SessionID: input.SessionID,
		StepCount: len(plan.Steps),
		Summary:   planSummary(plan),
	}) {
		return
	}

	state := newExecutionState()
	results := []stepResult{}
	replans := 0
	for index := 0; index < len(plan.Steps); index++ {
		step := plan.Steps[index]
		runState.Round = index + 1
		stats.Rounds = index + 1
		if !agent.Emit(ctx, events, agent.StepStartedEvent{
			SessionID:   input.SessionID,
			StepID:      step.ID,
			Description: step.Description,
		}) {
			return
		}

		result, err := a.executeStep(ctx, events, input, step, state, stats)
		if err != nil {
			result = stepResult{StepID: step.ID, Status: "failed", Error: err.Error()}
			results = append(results, result)
			if !agent.Emit(ctx, events, agent.StepFailedEvent{SessionID: input.SessionID, StepID: step.ID, Err: err}) {
				return
			}
			if replans >= a.maxReplans {
				fail(err)
				return
			}
			replans++
			if !agent.Emit(ctx, events, agent.ReplanStartedEvent{SessionID: input.SessionID, StepID: step.ID}) {
				return
			}
			updated, err := a.replan(ctx, events, input, plan, results, state, stats)
			if err != nil {
				fail(err)
				return
			}
			plan = updated
			if !agent.Emit(ctx, events, agent.PlanUpdatedEvent{SessionID: input.SessionID, Summary: planSummary(plan)}) {
				return
			}
			index = len(results) - 1
			continue
		}

		results = append(results, result)
		state.apply(result)
		if !agent.Emit(ctx, events, agent.StepCompletedEvent{SessionID: input.SessionID, StepID: step.ID, Summary: result.Summary}) {
			return
		}
	}

	if !agent.Emit(ctx, events, agent.FinalizingEvent{SessionID: input.SessionID}) {
		return
	}
	answer, err := a.finalize(ctx, events, input, plan, results, state, stats)
	if err != nil {
		fail(err)
		return
	}
	input.Conversation.Append(answer)
	if !agent.Emit(ctx, events, agent.FinalAnswerEvent{SessionID: input.SessionID, Round: stats.Rounds, Message: answer}) {
		return
	}
	complete(agent.TerminalCompleted)
}
