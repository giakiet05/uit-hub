package react

import (
	"context"
	"errors"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
)

func (a *ReActAgent) generate(ctx context.Context, events chan<- agent.Event, round int, request llm.GenerateRequest) (llm.GenerateResponse, bool, error) {
	streamer, ok := a.provider.(llm.StreamProvider)
	if !ok {
		response, err := a.provider.Generate(ctx, request)
		return response, false, err
	}

	var response llm.GenerateResponse
	textStreamed := false
	for event := range streamer.Stream(ctx, request) {
		switch typed := event.(type) {
		case llm.TextDeltaEvent:
			textStreamed = true
			if !agent.Emit(ctx, events, agent.ModelTextDeltaEvent{
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
