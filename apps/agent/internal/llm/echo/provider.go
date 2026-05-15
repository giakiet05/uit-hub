package echo

import (
	"context"

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
		Message:      llm.NewTextMessage(llm.RoleAssistant, "received: "+userText),
		FinishReason: llm.FinishReasonStop,
	}, nil
}

func lastUserText(messages []llm.Message) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == llm.RoleUser {
			return messages[i].ContentText()
		}
	}
	return ""
}
