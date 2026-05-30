package planexecute

import (
	"context"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

func (a *PlanAndExecuteAgent) callModel(
	ctx context.Context,
	events chan<- agent.Event,
	sessionID string,
	round int,
	messages []conversation.Message,
	tools []tool.Definition,
	stats *agent.RunStats,
) (llm.GenerateResponse, error) {
	if !agent.Emit(ctx, events, agent.ModelCallStartedEvent{SessionID: sessionID, Round: round}) {
		return llm.GenerateResponse{}, ctx.Err()
	}
	startedAt := time.Now()
	stats.LLMCalls++
	response, err := a.provider.Generate(ctx, llm.GenerateRequest{
		SessionID: sessionID,
		Messages:  messages,
		Tools:     tools,
	})
	duration := time.Since(startedAt)
	if err != nil {
		return llm.GenerateResponse{}, err
	}
	stats.TokenUsage.Add(response.Usage)

	toolCalls := agent.AssistantToolCalls(response.Message)
	if !agent.Emit(ctx, events, agent.ModelCallCompletedEvent{
		SessionID:     sessionID,
		Round:         round,
		Usage:         response.Usage,
		Duration:      duration,
		ToolCallNames: agent.ToolCallNames(toolCalls),
		MessageText:   conversation.Text(response.Message),
	}) {
		return llm.GenerateResponse{}, ctx.Err()
	}
	return response, nil
}
