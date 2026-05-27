package planexecute

import (
	"context"
	"errors"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent/loop"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

func (a *PlanAndExecuteAgent) executeToolCall(
	ctx context.Context,
	events chan<- agent.Event,
	sessionID string,
	round int,
	call conversation.ToolCall,
	stats *agent.RunStats,
) (tool.Result, error) {
	if a.tools == nil {
		return tool.Result{}, errors.New("no tools are registered")
	}
	stats.ToolCalls++
	if !loop.Emit(ctx, events, agent.ToolCallStartedEvent{SessionID: sessionID, Round: round, Call: call, Timeout: a.toolTimeout}) {
		return tool.Result{}, ctx.Err()
	}

	toolCtx := ctx
	if a.toolTimeout > 0 {
		var cancel context.CancelFunc
		toolCtx, cancel = context.WithTimeout(ctx, a.toolTimeout)
		defer cancel()
	}
	startedAt := time.Now()
	result, err := a.tools.Execute(toolCtx, tool.Call{
		ID:        call.ID,
		Name:      call.Name,
		Arguments: call.Arguments,
	})
	duration := time.Since(startedAt)
	if err != nil {
		stats.ToolFailures++
		observation := loop.ToolErrorObservation(err)
		loop.Emit(ctx, events, agent.ToolCallFailedEvent{
			SessionID:   sessionID,
			Round:       round,
			Call:        call,
			Observation: observation,
			Duration:    duration,
		})
		return tool.Result{}, err
	}

	result = loop.NormalizeToolResult(call, result)
	a.tools.MarkUsed(call.Name)
	if !loop.Emit(ctx, events, agent.ToolCallCompletedEvent{SessionID: sessionID, Round: round, Result: result, Duration: duration}) {
		return tool.Result{}, ctx.Err()
	}
	return result, nil
}

func (a *PlanAndExecuteAgent) getToolDefinitions() []tool.Definition {
	if a.tools == nil {
		return nil
	}
	return a.tools.Definitions()
}

func (a *PlanAndExecuteAgent) stepToolDefinitions(step planStep) []tool.Definition {
	if a.tools == nil {
		return nil
	}

	names := []string{
		"load_mcp_tool",
		"calculator",
		"read_file",
		"write_file",
		"memory_read",
		"memory_write",
		"memory_list",
	}
	names = append(names, step.SuggestedTools...)

	seen := map[string]struct{}{}
	definitions := []tool.Definition{}
	for _, name := range names {
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}

		definition, exists := a.tools.Definition(name)
		if !exists {
			continue
		}
		definitions = append(definitions, definition)
	}

	if len(definitions) == 0 {
		return a.tools.Definitions()
	}
	return definitions
}
