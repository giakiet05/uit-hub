package llm

import (
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
	"github.com/giakiet05/uit-hub/apps/agent/internal/usage"
)

// GenerateRequest is the provider-neutral input for one LLM call.
type GenerateRequest struct {
	SessionID string
	Messages  []conversation.Message
	Tools     []tool.Definition
}

// GenerateResponse is the provider-neutral output from one LLM call.
type GenerateResponse struct {
	Message conversation.Message
	Usage   usage.TokenUsage
}
