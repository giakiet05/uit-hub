package eval

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
)

// Recorder consumes agent events and builds eval traces and run summaries.
type Recorder struct {
	result Result
	trace  Trace
}

// NewRecorder creates a recorder for one eval case.
func NewRecorder(testCase Case) *Recorder {
	startedAt := time.Now().UTC()
	return &Recorder{
		result: Result{
			CaseID:    testCase.ID,
			Title:     testCase.Title,
			StartedAt: startedAt,
		},
		trace: Trace{
			Events: []TraceEvent{},
		},
	}
}

// Consume drains the agent event stream.
func (r *Recorder) Consume(ctx context.Context, events <-chan agent.Event) {
	for {
		select {
		case <-ctx.Done():
			r.result.Error = ctx.Err().Error()
			return
		case event, ok := <-events:
			if !ok {
				return
			}
			r.Record(event)
		}
	}
}

// Record appends one agent event to the eval result and trace.
func (r *Recorder) Record(event agent.Event) {
	switch typed := event.(type) {
	case agent.RunStartedEvent:
		r.addTrace("run started: max_rounds=%d", typed.MaxRounds)
	case agent.RoundStartedEvent:
		// We do not overwrite RoundCount here for multi-turn to work organically.
		r.addTrace("round %d: messages=%d tools=%d", typed.Round, typed.MessageCount, typed.ToolCount)
	case agent.ModelCallStartedEvent:
		r.result.LLMCalls++
		r.addTrace("round %d: model call started", typed.Round)
	case agent.ModelCallCompletedEvent:
		r.result.TokenUsage.Add(typed.Usage)
		if len(typed.ToolCallNames) == 0 {
			r.addTrace("round %d: model completed without tool calls", typed.Round)
			return
		}
		r.addTrace("round %d: model requested tools: %s", typed.Round, strings.Join(typed.ToolCallNames, ", "))
	case agent.PlanCreatedEvent:
		r.addTrace("plan created: steps=%d summary=%s", typed.StepCount, typed.Summary)
	case agent.StepStartedEvent:
		r.addTrace("step started: %s -> %s", typed.StepID, typed.Description)
	case agent.StepCompletedEvent:
		r.addTrace("step completed: %s -> %s", typed.StepID, preview(typed.Summary))
	case agent.StepFailedEvent:
		r.addTrace("step failed: %s -> %v", typed.StepID, typed.Err)
	case agent.ReplanStartedEvent:
		r.addTrace("replan started after step: %s", typed.StepID)
	case agent.PlanUpdatedEvent:
		r.addTrace("plan updated: %s", typed.Summary)
	case agent.FinalizingEvent:
		r.addTrace("finalizing")
	case agent.ToolCallStartedEvent:
		r.result.ToolCalls++
		r.result.ToolCallNames = append(r.result.ToolCallNames, typed.Call.Name)
		r.addTrace("round %d: tool started: %s", typed.Round, typed.Call.Name)
	case agent.ToolCallCompletedEvent:
		r.addTrace("round %d: tool completed: %s -> %s", typed.Round, typed.Result.Name, preview(typed.Result.Content))
	case agent.ToolCallFailedEvent:
		r.result.ToolFailures++
		r.addTrace("round %d: tool failed: %s -> %s", typed.Round, typed.Call.Name, typed.Observation)
	case agent.FinalAnswerEvent:
		r.result.FinalAnswer = conversation.Text(typed.Message)
		r.addTrace("round %d: final answer: %s", typed.Round, preview(r.result.FinalAnswer))
	case agent.RunCompletedEvent:
		r.result.TerminalReason = typed.Reason
		r.result.RoundCount += typed.Stats.Rounds // Accumulated per run
		r.addTrace("run completed: %s", typed.Reason)
	case agent.RunFailedEvent:
		r.result.Error = typed.Err.Error()
		r.result.RoundCount += typed.Stats.Rounds
		r.addTrace("run failed: %s", typed.Err.Error())
	}
}

// Finish applies judge results and records the final timestamp.
func (r *Recorder) Finish(checks []CheckResult) Result {
	finishedAt := time.Now().UTC()
	r.result.FinishedAt = finishedAt
	r.result.Duration = finishedAt.Sub(r.result.StartedAt).String()
	r.result.CheckResults = checks
	r.result.Passed = checksPassed(checks)
	return r.result
}

// Trace returns the recorded event trace.
func (r *Recorder) Trace() Trace {
	return r.trace
}

func (r *Recorder) addTrace(format string, args ...any) {
	r.trace.Events = append(r.trace.Events, TraceEvent{
		Time: time.Now().UTC(),
		Text: fmt.Sprintf(format, args...),
	})
}

func checksPassed(checks []CheckResult) bool {
	for _, check := range checks {
		if !check.Passed {
			return false
		}
	}
	return true
}

func preview(text string) string {
	text = strings.TrimSpace(strings.ReplaceAll(text, "\n", " "))
	const max = 240
	if len(text) <= max {
		return text
	}
	return text[:max] + "..."
}
