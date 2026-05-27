package react

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent/loop"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// getToolDefinitions returns the current tool definitions, or nil when the
// agent has no registry.
func (a *ReActAgent) getToolDefinitions() []tool.Definition {
	if a.tools == nil {
		return nil
	}
	return a.tools.Definitions()
}

// executeToolCalls runs every requested tool call and appends each observation
// back into the session conversation.
func (a *ReActAgent) executeToolCalls(
	ctx context.Context,
	events chan<- agent.Event,
	session *runtime.Session,
	round int,
	calls []conversation.ToolCall,
	stats *agent.RunStats,
) bool {
	if a.tools == nil {
		stats.Finish()
		a.logRunStats(ctx, session.ID, stats)
		loop.Emit(ctx, events, agent.RunFailedEvent{
			SessionID: session.ID,
			Err:       errors.New("assistant requested tool call but no tools are registered"),
			Stats:     *stats,
		})
		return false
	}

	for _, call := range calls {
		stats.ToolCalls++
		if !loop.Emit(ctx, events, agent.ToolCallStartedEvent{
			SessionID: session.ID,
			Round:     round,
			Call:      call,
			Timeout:   a.toolTimeout,
		}) {
			return false
		}
		a.trace(
			ctx,
			"react.tool.started",
			fmt.Sprintf("Round %d: Tool %s started", round, call.Name),
			"session_id", session.ID,
			"round", round,
			"tool_call_id", call.ID,
			"tool_name", call.Name,
			"arguments_preview", loop.PreviewValue(call.Arguments),
			"timeout", a.toolTimeout.String(),
		)
		toolStartedAt := time.Now()
		result, err := a.executeToolCall(ctx, call)
		toolDuration := time.Since(toolStartedAt)
		if err != nil {
			stats.ToolFailures++
			observation := loop.ToolErrorObservation(err)
			a.trace(
				ctx,
				"react.tool.failed",
				fmt.Sprintf("Round %d: Tool %s failed", round, call.Name),
				"session_id", session.ID,
				"round", round,
				"tool_call_id", call.ID,
				"tool_name", call.Name,
				"duration", toolDuration.String(),
				"error", err,
			)
			result = tool.Result{
				CallID:  call.ID,
				Name:    call.Name,
				Content: observation,
			}
			if !loop.Emit(ctx, events, agent.ToolCallFailedEvent{
				SessionID:   session.ID,
				Round:       round,
				Call:        call,
				Observation: observation,
				Duration:    toolDuration,
			}) {
				return false
			}
		}

		result = loop.NormalizeToolResult(call, result)
		if err == nil {
			a.tools.MarkUsed(call.Name)
			if !loop.Emit(ctx, events, agent.ToolCallCompletedEvent{
				SessionID: session.ID,
				Round:     round,
				Result:    result,
				Duration:  toolDuration,
			}) {
				return false
			}
		}
		session.Conversation.Append(conversation.NewToolResultMessage(result.CallID, result.Content))
		a.trace(
			ctx,
			"react.tool.observation_appended",
			fmt.Sprintf("Round %d: Observe %s", round, result.Name),
			"session_id", session.ID,
			"round", round,
			"tool_call_id", result.CallID,
			"tool_name", result.Name,
			"result_chars", len(result.Content),
			"result_preview", loop.PreviewText(result.Content),
			"duration", toolDuration.String(),
		)
	}

	return true
}

// executeToolCall runs a single tool call with the configured timeout.
func (a *ReActAgent) executeToolCall(ctx context.Context, call conversation.ToolCall) (tool.Result, error) {
	if a.toolTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, a.toolTimeout)
		defer cancel()
	}

	result, err := a.tools.Execute(ctx, tool.Call{
		ID:        call.ID,
		Name:      call.Name,
		Arguments: call.Arguments,
	})
	if err != nil {
		return tool.Result{}, fmt.Errorf("execute tool call: %w", err)
	}
	return result, nil
}
