// Package agent defines the runtime contract and concrete agent loops.
package agent

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// RunInput contains the session-owned state an agent loop needs for one user
// prompt.
type RunInput struct {
	SessionID       string
	UserPrompt      string
	Conversation    *conversation.Conversation
	PromptSnapshot  prompt.SystemPrompt
	Tools           *tool.ToolSet
	ResultBudgeter  *tool.ResultBudgeter
	ConcurrentTools bool
}

// ToolDefinitions returns all tools visible to the model for this run.
func (i RunInput) ToolDefinitions() []tool.Definition {
	if i.Tools == nil {
		return nil
	}
	return i.Tools.Definitions()
}

// ToolDefinition returns one visible tool definition by name.
func (i RunInput) ToolDefinition(name string) (tool.Definition, bool) {
	if i.Tools == nil {
		return tool.Definition{}, false
	}
	return i.Tools.Definition(name)
}

// Agent runs one user prompt and streams loop events until a terminal event is
// emitted and the channel closes.
type Agent interface {
	Run(ctx context.Context, input RunInput) <-chan Event
}
