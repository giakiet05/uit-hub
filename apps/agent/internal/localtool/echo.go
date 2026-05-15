package localtool

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

type Echo struct{}

func NewEcho() *Echo {
	return &Echo{}
}

func (t *Echo) Definition() llm.ToolDefinition {
	return llm.ToolDefinition{
		Name:        "echo",
		Description: "Return the provided text exactly as received.",
		InputSchema: map[string]any{
			"type":                 "object",
			"additionalProperties": false,
			"properties": map[string]any{
				"text": map[string]any{
					"type":        "string",
					"description": "The text to echo.",
				},
			},
			"required": []string{"text"},
		},
	}
}

func (t *Echo) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	if err := ctx.Err(); err != nil {
		return tool.Result{}, err
	}

	text, err := stringArg(call.Arguments, "text")
	if err != nil {
		return tool.Result{}, err
	}

	return tool.Result{
		CallID:  call.ID,
		Name:    call.Name,
		Content: text,
	}, nil
}
