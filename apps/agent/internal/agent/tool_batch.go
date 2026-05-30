package agent

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// RunToolBatch executes a list of tool calls, updating run stats and emitting
// events during the process. It implements Batch Orchestration (concurrent tool execution)
// by grouping calls into concurrent and serial batches based on tool metadata.
// It returns the tool results in submission order.
func RunToolBatch(
	ctx context.Context,
	logger *slog.Logger,
	agentType string,
	events chan<- Event,
	input RunInput,
	round int,
	calls []conversation.ToolCall,
	timeout time.Duration,
	stats *RunStats,
) ([]tool.Result, bool) {
	if input.Tools == nil {
		Emit(ctx, events, RunFailedEvent{
			SessionID: input.SessionID,
			Err:       errors.New("assistant requested tool call but no tools are registered"),
			Stats:     *stats,
		})
		return nil, false
	}

	batches := partitionToolCalls(calls, input.Tools)
	executor := tool.NewExecutor(input.Tools, timeout)
	var finalResults []tool.Result

	for i, b := range batches {
		if b.parallel && input.ConcurrentTools {
			if len(b.calls) > 1 {
				logger.DebugContext(ctx, "Executing tool batch concurrently", "batch_index", i, "tool_count", len(b.calls))
			} else {
				logger.DebugContext(ctx, "Executing single safe tool", "batch_index", i, "tool_name", b.calls[0].Name)
			}
			batchResults, ok := runParallelBatch(ctx, logger, agentType, events, input, round, b.calls, timeout, stats, executor)
			if !ok {
				return nil, false
			}
			finalResults = append(finalResults, batchResults...)
		} else {
			if len(b.calls) == 1 {
				logger.DebugContext(ctx, "Executing single unsafe tool", "batch_index", i, "tool_name", b.calls[0].Name)
			} else {
				logger.DebugContext(ctx, "Executing tool batch serially", "batch_index", i, "tool_count", len(b.calls), "safe", b.parallel)
			}
			batchResults, ok := runSerialBatch(ctx, logger, agentType, events, input, round, b.calls, timeout, stats, executor)
			if !ok {
				return nil, false
			}
			finalResults = append(finalResults, batchResults...)
		}
	}

	return finalResults, true
}

type batch struct {
	parallel bool
	calls    []conversation.ToolCall
}

func partitionToolCalls(calls []conversation.ToolCall, tools *tool.ToolSet) []batch {
	var batches []batch
	for _, call := range calls {
		toolImpl, ok := tools.Tool(call.Name)
		safe := false
		if ok {
			safe = toolImpl.Metadata().ConcurrencySafe
		}

		if safe && len(batches) > 0 && batches[len(batches)-1].parallel {
			batches[len(batches)-1].calls = append(batches[len(batches)-1].calls, call)
		} else {
			batches = append(batches, batch{parallel: safe, calls: []conversation.ToolCall{call}})
		}
	}
	return batches
}

func runSerialBatch(
	ctx context.Context,
	logger *slog.Logger,
	agentType string,
	events chan<- Event,
	input RunInput,
	round int,
	calls []conversation.ToolCall,
	timeout time.Duration,
	stats *RunStats,
	executor *tool.Executor,
) ([]tool.Result, bool) {
	results := make([]tool.Result, 0, len(calls))

	for _, call := range calls {
		result, isErr, ok := executeOneTool(ctx, logger, agentType, events, input, round, call, timeout, executor)
		if !ok {
			return nil, false
		}
		stats.ToolCalls++
		if isErr {
			stats.ToolFailures++
		}
		results = append(results, result)
	}
	return results, true
}

func runParallelBatch(
	ctx context.Context,
	logger *slog.Logger,
	agentType string,
	events chan<- Event,
	input RunInput,
	round int,
	calls []conversation.ToolCall,
	timeout time.Duration,
	stats *RunStats,
	executor *tool.Executor,
) ([]tool.Result, bool) {
	results := make([]tool.Result, len(calls))
	var wg sync.WaitGroup
	var mu sync.Mutex
	okFlag := true

	for i, call := range calls {
		wg.Add(1)
		go func(i int, call conversation.ToolCall) {
			defer wg.Done()
			
			// We check if another concurrent execution already failed context (e.g. timeout/cancel).
			if ctx.Err() != nil {
				return
			}
			
			result, isErr, ok := executeOneTool(ctx, logger, agentType, events, input, round, call, timeout, executor)
			
			mu.Lock()
			if !ok {
				okFlag = false
			} else {
				stats.ToolCalls++
				if isErr {
					stats.ToolFailures++
				}
				results[i] = result
			}
			mu.Unlock()
		}(i, call)
	}

	wg.Wait()
	if !okFlag || ctx.Err() != nil {
		return nil, false
	}
	return results, true
}

func executeOneTool(
	ctx context.Context,
	logger *slog.Logger,
	agentType string,
	events chan<- Event,
	input RunInput,
	round int,
	call conversation.ToolCall,
	timeout time.Duration,
	executor *tool.Executor,
) (tool.Result, bool, bool) {
	metadata, hasMetadata := executor.Metadata(call.Name)
	if hasMetadata && metadata.RequireApproval {
		responseChan := make(chan bool, 1)
		if !Emit(ctx, events, ToolPermissionRequestEvent{
			SessionID: input.SessionID,
			Round:     round,
			Call:      call,
			Response:  responseChan,
		}) {
			return tool.Result{}, false, false
		}

		var approved bool
		select {
		case <-ctx.Done():
			return tool.Result{}, false, false
		case approved = <-responseChan:
		}

		if !approved {
			observation := "User rejected this tool execution."
			result := tool.Result{
				CallID:  call.ID,
				Name:    call.Name,
				Content: observation,
			}
			if !Emit(ctx, events, ToolCallFailedEvent{
				SessionID:   input.SessionID,
				Round:       round,
				Call:        call,
				Observation: observation,
				Duration:    0,
			}) {
				return tool.Result{}, true, false
			}
			return result, true, true
		}
	}

	if !Emit(ctx, events, ToolCallStartedEvent{
		SessionID: input.SessionID,
		Round:     round,
		Call:      call,
		Timeout:   timeout,
	}) {
		return tool.Result{}, false, false
	}

	Trace(logger, ctx,
		agentType+".tool.started",
		fmt.Sprintf("Round %d: Tool %s started", round, call.Name),
		"session_id", input.SessionID,
		"round", round,
		"tool_call_id", call.ID,
		"tool_name", call.Name,
		"arguments_preview", PreviewValue(call.Arguments),
	)

	startedAt := time.Now()
	result, err := executor.Execute(ctx, tool.Call{
		ID:        call.ID,
		Name:      call.Name,
		Arguments: call.Arguments,
	})
	duration := time.Since(startedAt)

	if err != nil {
		observation := tool.ErrorObservation(err)
		result = tool.Result{
			CallID:  call.ID,
			Name:    call.Name,
			Content: observation,
		}
		if !Emit(ctx, events, ToolCallFailedEvent{
			SessionID:   input.SessionID,
			Round:       round,
			Call:        call,
			Observation: observation,
			Duration:    duration,
		}) {
			return tool.Result{}, true, false
		}
		return result, true, true
	} else {
		if !Emit(ctx, events, ToolCallCompletedEvent{
			SessionID: input.SessionID,
			Round:     round,
			Result:    result,
			Duration:  duration,
		}) {
			return tool.Result{}, false, false
		}
		return result, false, true
	}
}
