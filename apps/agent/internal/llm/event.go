package llm

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

func (TextDeltaEvent) isStreamEvent()       {}
func (StreamCompletedEvent) isStreamEvent() {}
func (StreamFailedEvent) isStreamEvent()    {}
