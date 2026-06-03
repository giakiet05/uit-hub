package localtool

import (
	"context"
	"errors"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// Calculator performs basic arithmetic for tool-calling tests.
type Calculator struct {
	tool.BaseTool
}

// NewCalculator creates a calculator tool.
func NewCalculator() *Calculator {
	return &Calculator{
		BaseTool: tool.NewBaseTool(
			tool.Definition{
				Name:        "calculator",
				Description: "Run a basic arithmetic operation on two numbers.",
				InputSchema: tool.ObjectSchema(
					map[string]any{
						"operation": tool.StringEnumProperty(
							"The arithmetic operation to run.",
							"add",
							"subtract",
							"multiply",
							"divide",
						),
						"a": tool.NumberProperty("The left operand."),
						"b": tool.NumberProperty("The right operand."),
					},
					"operation",
					"a",
					"b",
				),
			},
			tool.NewReadOnlyMetadata(true),
		),
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
