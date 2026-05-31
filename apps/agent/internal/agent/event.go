package agent

import (
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/eventbus"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
	"github.com/giakiet05/uit-hub/apps/agent/internal/usage"
)

// Event is the closed set of events emitted by an agent loop.
type Event interface {
	isAgentEvent()
	Topic() eventbus.Topic
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
func (RunStartedEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (RoundStartedEvent) isAgentEvent()         {}
func (RoundStartedEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (ModelCallStartedEvent) isAgentEvent()         {}
func (ModelCallStartedEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (ModelCallCompletedEvent) isAgentEvent()         {}
func (ModelCallCompletedEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (ModelTextDeltaEvent) isAgentEvent()         {}
func (ModelTextDeltaEvent) Topic() eventbus.Topic { return eventbus.TopicStream }

func (PlanCreatedEvent) isAgentEvent()         {}
func (PlanCreatedEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (StepStartedEvent) isAgentEvent()         {}
func (StepStartedEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (StepCompletedEvent) isAgentEvent()         {}
func (StepCompletedEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (StepFailedEvent) isAgentEvent()         {}
func (StepFailedEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (ReplanStartedEvent) isAgentEvent()         {}
func (ReplanStartedEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (PlanUpdatedEvent) isAgentEvent()         {}
func (PlanUpdatedEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (FinalizingEvent) isAgentEvent()         {}
func (FinalizingEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (ToolCallStartedEvent) isAgentEvent()         {}
func (ToolCallStartedEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (ToolPermissionRequestEvent) isAgentEvent()         {}
func (ToolPermissionRequestEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (ToolCallCompletedEvent) isAgentEvent()         {}
func (ToolCallCompletedEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (ToolCallFailedEvent) isAgentEvent()         {}
func (ToolCallFailedEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (FinalAnswerEvent) isAgentEvent()         {}
func (FinalAnswerEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (RunCompletedEvent) isAgentEvent()         {}
func (RunCompletedEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }

func (RunFailedEvent) isAgentEvent()         {}
func (RunFailedEvent) Topic() eventbus.Topic { return eventbus.TopicLifecycle }
