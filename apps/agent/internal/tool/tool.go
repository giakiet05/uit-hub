package tool

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
)

type Tool interface {
	Definition() llm.ToolDefinition
	Execute(ctx context.Context, call Call) (Result, error)
}

type Call struct {
	ID        string
	Name      string
	Arguments map[string]any
}

type Result struct {
	CallID  string
	Name    string
	Content string
}
