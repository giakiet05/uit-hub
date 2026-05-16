package echo

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) Generate(ctx context.Context, request llm.GenerateRequest) (llm.GenerateResponse, error) {
	if err := ctx.Err(); err != nil {
		return llm.GenerateResponse{}, err
	}

	userText := lastUserText(request.Messages)
	return llm.GenerateResponse{
		Message: conversation.NewAssistantMessage("received: "+userText, nil),
	}, nil
}

func lastUserText(messages []conversation.Message) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if _, ok := messages[i].(conversation.UserMessage); ok {
			return conversation.Text(messages[i])
		}
	}
	return ""
}
