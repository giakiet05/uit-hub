// Package localtool contains temporary native tools used to exercise the agent
// loop before MCP adapters are introduced.
package localtool

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// Echo returns text exactly as provided by the model.
type Echo struct{}

// NewEcho creates an echo tool.
func NewEcho() *Echo {
	return &Echo{}
}

// Definition describes the echo tool schema.
func (t *Echo) Definition() tool.Definition {
	return tool.Definition{
		Name:        "echo",
		Description: "Return the provided text exactly as received.",
		InputSchema: tool.ObjectSchema(
			map[string]any{
				"text": tool.StringProperty("The text to echo."),
			},
			"text",
		),
	}
}

// Execute validates the text argument and returns it unchanged.
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
