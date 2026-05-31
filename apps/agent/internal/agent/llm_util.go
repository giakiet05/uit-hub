package agent

import (
	"context"
	"errors"

	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
)

// GenerateWithStream calls the LLM provider and emits streaming events to the UI
// if the provider supports it. It acts as a bridge between the neutral LLM layer
// and the agent's reactive event loop.
func GenerateWithStream(
	ctx context.Context,
	provider llm.Provider,
	request llm.GenerateRequest,
	events chan<- Event,
	round int,
) (llm.GenerateResponse, bool, error) {
	streamer, ok := provider.(llm.StreamProvider)
	if !ok {
		// Fallback to synchronous generation if provider doesn't support streaming
		response, err := provider.Generate(ctx, request)
		return response, false, err
	}

	var response llm.GenerateResponse
	textStreamed := false
	for event := range streamer.Stream(ctx, request) {
		switch typed := event.(type) {
		case llm.TextDeltaEvent:
			textStreamed = true
			if !Emit(ctx, events, ModelTextDeltaEvent{
				SessionID: request.SessionID,
				Round:     round,
				Delta:     typed.Delta,
			}) {
				return llm.GenerateResponse{}, textStreamed, ctx.Err()
			}
		case llm.StreamCompletedEvent:
			response = typed.Response
			return response, textStreamed, nil
		case llm.StreamFailedEvent:
			return llm.GenerateResponse{}, textStreamed, typed.Err
		}
	}

	return llm.GenerateResponse{}, textStreamed, errors.New("llm stream closed without completion")
}
