package tool

import (
	"context"
	"fmt"
	"time"
)

// Executor owns the common tool-call pipeline: lookup, timeout, execution,
// result normalization, result budgeting, and runtime usage tracking.
type Executor struct {
	tools              *ToolSet
	timeout            time.Duration
	budgeter           *ResultBudgeter
	loadedMCPToolNames map[string]struct{}
}

// NewExecutor creates a tool executor for one agent run.
func NewExecutor(tools *ToolSet, timeout time.Duration, budgeter *ResultBudgeter) *Executor {
	return &Executor{
		tools:              tools,
		timeout:            timeout,
		budgeter:           budgeter,
		loadedMCPToolNames: make(map[string]struct{}),
	}
}

// Metadata returns the metadata for a tool by name, if it exists.
func (e *Executor) Metadata(name string) (Metadata, bool) {
	if e == nil || e.tools == nil {
		return Metadata{}, false
	}
	selected, exists := e.tools.Tool(name)
	if !exists {
		return Metadata{}, false
	}
	return selected.Metadata(), true
}

// Execute runs one tool call through the pipeline.
func (e *Executor) Execute(ctx context.Context, call Call) (Result, error) {
	if result, ok := e.duplicateMCPToolLoadResult(call); ok {
		return result, nil
	}

	if e == nil || e.tools == nil {
		return Result{}, ClassifyError(call.Name, fmt.Errorf("execute tool %q: %w", call.Name, ErrNotFound))
	}

	selected, exists := e.tools.Tool(call.Name)
	if !exists {
		return Result{}, ClassifyError(call.Name, fmt.Errorf("execute tool %q: %w", call.Name, ErrNotFound))
	}

	toolCtx := ctx
	if e.timeout > 0 {
		var cancel context.CancelFunc
		toolCtx, cancel = context.WithTimeout(ctx, e.timeout)
		defer cancel()
	}

	result, err := selected.Execute(toolCtx, call)
	if err != nil {
		return Result{}, ClassifyError(call.Name, fmt.Errorf("execute tool %q: %w", call.Name, err))
	}

	result = normalizeResult(call, result)
	result = e.budgeter.Apply(result)
	e.tools.MarkUsed(call.Name)
	e.recordMCPToolLoad(call)
	return result, nil
}

func (e *Executor) duplicateMCPToolLoadResult(call Call) (Result, bool) {
	if e == nil || call.Name != "load_mcp_tool" {
		return Result{}, false
	}
	name, ok := call.Arguments["name"].(string)
	if !ok || name == "" {
		return Result{}, false
	}
	if _, exists := e.loadedMCPToolNames[name]; !exists {
		return Result{}, false
	}
	return Result{
		CallID:  call.ID,
		Name:    call.Name,
		Content: fmt.Sprintf("MCP tool %s was already loaded in this round. Do not call load_mcp_tool again for this tool; use %s directly in the next round.", name, name),
	}, true
}

func (e *Executor) recordMCPToolLoad(call Call) {
	if e == nil || call.Name != "load_mcp_tool" {
		return
	}
	name, ok := call.Arguments["name"].(string)
	if !ok || name == "" {
		return
	}
	e.loadedMCPToolNames[name] = struct{}{}
}

func normalizeResult(call Call, result Result) Result {
	if result.CallID == "" {
		result.CallID = call.ID
	}
	if result.Name == "" {
		result.Name = call.Name
	}
	return result
}
