package llm

import (
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

type GenerateRequest struct {
	SessionID string
	Model     string
	Messages  []conversation.Message
	Tools     []tool.Definition
}

type GenerateResponse struct {
	Message conversation.Message
	Usage   Usage
}

type Usage struct {
	InputTokens  int
	OutputTokens int
}
