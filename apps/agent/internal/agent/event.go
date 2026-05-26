package agent

import (
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
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
	Usage         llm.Usage
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

// ToolCallStartedEvent is emitted before executing one tool call.
type ToolCallStartedEvent struct {
	SessionID string
	Round     int
	Call      conversation.ToolCall
	Timeout   time.Duration
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
func (ToolCallStartedEvent) isAgentEvent()    {}
func (ToolCallCompletedEvent) isAgentEvent()  {}
func (ToolCallFailedEvent) isAgentEvent()     {}
func (FinalAnswerEvent) isAgentEvent()        {}
func (RunCompletedEvent) isAgentEvent()       {}
func (RunFailedEvent) isAgentEvent()          {}
