package agent

import (
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
	"github.com/giakiet05/uit-hub/apps/agent/internal/usage"
)

// Event is the closed set of events emitted by an agent loop.
type Event interface {
	isAgentEvent()
}

// TerminalReason explains why an agent loop stopped.
type TerminalReason string

const (
	// TerminalCompleted means the model produced a final answer.
	TerminalCompleted TerminalReason = "completed"
	// TerminalMaxRounds means the loop stopped at the configured round limit.
	TerminalMaxRounds TerminalReason = "max_rounds"
	// TerminalAborted means the loop stopped because the context was canceled.
	TerminalAborted TerminalReason = "aborted"
)

// RunStartedEvent is emitted when a user prompt starts an agent run.
type RunStartedEvent struct {
	SessionID string
	MaxRounds int
}

// RoundStartedEvent is emitted before one model/tool round begins.
type RoundStartedEvent struct {
	SessionID    string
	Round        int
	MessageCount int
	ToolCount    int
}

// ModelCallStartedEvent is emitted immediately before calling the LLM provider.
type ModelCallStartedEvent struct {
	SessionID string
	Round     int
}

// ModelCallCompletedEvent is emitted after the LLM provider returns.
type ModelCallCompletedEvent struct {
	SessionID     string
	Round         int
	Usage         usage.TokenUsage
	Duration      time.Duration
	ToolCallNames []string
	MessageText   string
	TextStreamed  bool
}

// ModelTextDeltaEvent is emitted for visible assistant text streamed by the
// model before the final response is available.
type ModelTextDeltaEvent struct {
	SessionID string
	Round     int
	Delta     string
}

// PlanCreatedEvent is emitted when a plan-and-execute run creates its initial
// structured plan.
type PlanCreatedEvent struct {
	SessionID string
	StepCount int
	Summary   string
}

// StepStartedEvent is emitted before executing one planned step.
type StepStartedEvent struct {
	SessionID   string
	StepID      string
	Description string
}

// StepCompletedEvent is emitted after one planned step finishes successfully.
type StepCompletedEvent struct {
	SessionID string
	StepID    string
	Summary   string
}

// StepFailedEvent is emitted when one planned step fails before replanning.
type StepFailedEvent struct {
	SessionID string
	StepID    string
	Err       error
}

// ReplanStartedEvent is emitted before asking the model to repair a failed
// plan or step.
type ReplanStartedEvent struct {
	SessionID string
	StepID    string
}

// PlanUpdatedEvent is emitted after replanning changes the remaining work.
type PlanUpdatedEvent struct {
	SessionID string
	Summary   string
}

// FinalizingEvent is emitted before synthesizing the final answer from step
// results.
type FinalizingEvent struct {
	SessionID string
}

// ToolCallStartedEvent is emitted before executing one tool call.
type ToolCallStartedEvent struct {
	SessionID string
	Round     int
	Call      conversation.ToolCall
	Timeout   time.Duration
}

// ToolPermissionRequestEvent is emitted when a tool requires user approval before execution.
type ToolPermissionRequestEvent struct {
	SessionID string
	Round     int
	Call      conversation.ToolCall
	Response  chan bool
}

// ToolCallCompletedEvent is emitted after a tool call succeeds.
type ToolCallCompletedEvent struct {
	SessionID string
	Round     int
	Result    tool.Result
	Duration  time.Duration
}

// ToolCallFailedEvent is emitted after a tool call fails and is converted into
// a sanitized model observation.
type ToolCallFailedEvent struct {
	SessionID   string
	Round       int
	Call        conversation.ToolCall
	Observation string
	Duration    time.Duration
}

// FinalAnswerEvent is emitted when the model produces the final assistant
// answer for a run.
type FinalAnswerEvent struct {
	SessionID string
	Round     int
	Message   conversation.AssistantMessage
}

// RunCompletedEvent is emitted for non-error terminal states.
type RunCompletedEvent struct {
	SessionID string
	Reason    TerminalReason
	Stats     RunStats
}

// RunFailedEvent is emitted when the loop cannot continue because of an error.
type RunFailedEvent struct {
	SessionID string
	Err       error
	Stats     RunStats
}

func (RunStartedEvent) isAgentEvent()         {}
func (RoundStartedEvent) isAgentEvent()       {}
func (ModelCallStartedEvent) isAgentEvent()   {}
func (ModelCallCompletedEvent) isAgentEvent() {}
func (ModelTextDeltaEvent) isAgentEvent()     {}
func (PlanCreatedEvent) isAgentEvent()        {}
func (StepStartedEvent) isAgentEvent()        {}
func (StepCompletedEvent) isAgentEvent()      {}
func (StepFailedEvent) isAgentEvent()         {}
func (ReplanStartedEvent) isAgentEvent()      {}
func (PlanUpdatedEvent) isAgentEvent()        {}
func (FinalizingEvent) isAgentEvent()         {}
func (ToolCallStartedEvent) isAgentEvent()    {}
func (ToolPermissionRequestEvent) isAgentEvent() {}
func (ToolCallCompletedEvent) isAgentEvent()  {}
func (ToolCallFailedEvent) isAgentEvent()     {}
func (FinalAnswerEvent) isAgentEvent()        {}
func (RunCompletedEvent) isAgentEvent()       {}
func (RunFailedEvent) isAgentEvent()          {}
