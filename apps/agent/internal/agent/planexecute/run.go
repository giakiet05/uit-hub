package planexecute

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent/loop"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
)

func (a *PlanAndExecuteAgent) run(ctx context.Context, events chan<- agent.Event, session *session.State, userPrompt string) {
	defer close(events)

	runState := agent.NewRunState(session.ID, a.maxSteps)
	stats := runState.Stats
	fail := func(err error) {
		runState.Finish()
		a.logRunStats(ctx, session.ID, stats)
		session.AddUsage(agent.UsageDeltaFromRunStats(stats))
		loop.Emit(ctx, events, agent.RunFailedEvent{SessionID: session.ID, Err: err, Stats: *stats})
	}
	complete := func() {
		runState.Finish()
		a.logRunStats(ctx, session.ID, stats)
		session.AddUsage(agent.UsageDeltaFromRunStats(stats))
		loop.Emit(ctx, events, agent.RunCompletedEvent{SessionID: session.ID, Reason: agent.TerminalCompleted, Stats: *stats})
	}

	if !loop.Emit(ctx, events, agent.RunStartedEvent{SessionID: session.ID, MaxRounds: runState.MaxRounds}) {
		return
	}
	session.Conversation.Append(conversation.NewUserMessage(userPrompt))
	a.trace(ctx, "plan_execute.run.started", "Plan-and-execute agent started", "session_id", session.ID, "max_steps", a.maxSteps)

	plan, err := a.createPlan(ctx, events, session, userPrompt, stats)
	if err != nil {
		fail(err)
		return
	}
	if len(plan.Steps) > a.maxSteps {
		plan.Steps = plan.Steps[:a.maxSteps]
	}
	if !loop.Emit(ctx, events, agent.PlanCreatedEvent{
		SessionID: session.ID,
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
		if !loop.Emit(ctx, events, agent.StepStartedEvent{
			SessionID:   session.ID,
			StepID:      step.ID,
			Description: step.Description,
		}) {
			return
		}

		result, err := a.executeStep(ctx, events, session, userPrompt, step, state, stats)
		if err != nil {
			result = stepResult{StepID: step.ID, Status: "failed", Error: err.Error()}
			results = append(results, result)
			if !loop.Emit(ctx, events, agent.StepFailedEvent{SessionID: session.ID, StepID: step.ID, Err: err}) {
				return
			}
			if replans >= a.maxReplans {
				fail(err)
				return
			}
			replans++
			if !loop.Emit(ctx, events, agent.ReplanStartedEvent{SessionID: session.ID, StepID: step.ID}) {
				return
			}
			updated, err := a.replan(ctx, events, session, userPrompt, plan, results, state, stats)
			if err != nil {
				fail(err)
				return
			}
			plan = updated
			if !loop.Emit(ctx, events, agent.PlanUpdatedEvent{SessionID: session.ID, Summary: planSummary(plan)}) {
				return
			}
			index = len(results) - 1
			continue
		}

		results = append(results, result)
		state.apply(result)
		if !loop.Emit(ctx, events, agent.StepCompletedEvent{SessionID: session.ID, StepID: step.ID, Summary: result.Summary}) {
			return
		}
	}

	if !loop.Emit(ctx, events, agent.FinalizingEvent{SessionID: session.ID}) {
		return
	}
	answer, err := a.finalize(ctx, events, session, userPrompt, plan, results, state, stats)
	if err != nil {
		fail(err)
		return
	}
	session.Conversation.Append(answer)
	if !loop.Emit(ctx, events, agent.FinalAnswerEvent{SessionID: session.ID, Round: stats.Rounds, Message: answer}) {
		return
	}
	a.trace(ctx, "plan_execute.run.completed", "Plan-and-execute agent completed", "session_id", session.ID, "steps", len(results))
	complete()
}
