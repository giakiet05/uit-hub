package llm

import (
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// GenerateRequest is the provider-neutral input for one LLM call.
type GenerateRequest struct {
	SessionID string
	Model     string
	Messages  []conversation.Message
	Tools     []tool.Definition
}

// GenerateResponse is the provider-neutral output from one LLM call.
type GenerateResponse struct {
	Message conversation.Message
	Usage   Usage
}

// StreamEvent is the closed set of provider-neutral events from one streaming
// LLM call.
type StreamEvent interface {
	isStreamEvent()
}

// TextDeltaEvent carries a visible assistant text delta.
type TextDeltaEvent struct {
	Delta string
}

// StreamCompletedEvent carries the final mapped LLM response.
type StreamCompletedEvent struct {
	Response GenerateResponse
}

// StreamFailedEvent carries a streaming provider error.
type StreamFailedEvent struct {
	Err error
}

// Usage contains token counters reported by the provider.
type Usage struct {
	InputTokens  int
	OutputTokens int
}

func (TextDeltaEvent) isStreamEvent()       {}
func (StreamCompletedEvent) isStreamEvent() {}
func (StreamFailedEvent) isStreamEvent()    {}
