package localtool

import (
	"context"
	"errors"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// Calculator performs basic arithmetic for tool-calling tests.
type Calculator struct{}

// NewCalculator creates a calculator tool.
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Definition describes the calculator tool schema.
func (t *Calculator) Definition() tool.Definition {
	return tool.Definition{
		Name:        "calculator",
		Description: "Run a basic arithmetic operation on two numbers.",
		InputSchema: map[string]any{
			"type":                 "object",
			"additionalProperties": false,
			"properties": map[string]any{
				"operation": map[string]any{
					"type":        "string",
					"description": "The arithmetic operation to run.",
					"enum":        []string{"add", "subtract", "multiply", "divide"},
				},
				"a": map[string]any{
					"type":        "number",
					"description": "The left operand.",
				},
				"b": map[string]any{
					"type":        "number",
					"description": "The right operand.",
				},
			},
			"required": []string{"operation", "a", "b"},
		},
	}
}

// Execute runs the requested arithmetic operation.
func (t *Calculator) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	if err := ctx.Err(); err != nil {
		return tool.Result{}, err
	}

	operation, err := stringArg(call.Arguments, "operation")
	if err != nil {
		return tool.Result{}, err
	}
	a, err := numberArg(call.Arguments, "a")
	if err != nil {
		return tool.Result{}, err
	}
	b, err := numberArg(call.Arguments, "b")
	if err != nil {
		return tool.Result{}, err
	}

	var value float64
	switch operation {
	case "add":
		value = a + b
	case "subtract":
		value = a - b
	case "multiply":
		value = a * b
	case "divide":
		if b == 0 {
			return tool.Result{}, errors.New("division by zero")
		}
		value = a / b
	default:
		return tool.Result{}, fmt.Errorf("unsupported operation %q", operation)
	}

	return tool.Result{
		CallID:  call.ID,
		Name:    call.Name,
		Content: fmt.Sprintf("%g", value),
	}, nil
}
