package localtool

import (
	"context"
	"errors"
	"log/slog"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

type SubAgentTool struct {
	tool.BaseTool
	agent         agent.Agent
	tools         *tool.ToolSet
	promptBuilder func() prompt.SystemPrompt
}

// NewSubAgentTool wraps an agent into a tool.
func NewSubAgentTool(
	definition tool.Definition,
	agentInst agent.Agent,
	tools *tool.ToolSet,
	promptBuilder func() prompt.SystemPrompt,
) tool.Tool {
	return &SubAgentTool{
		BaseTool: tool.NewBaseTool(definition, tool.Metadata{
			ReadOnly:          false,
			Destructive:       false,
			ConcurrencySafe:   true,
			RequireApproval:   false,
			FinishOnInterrupt: false,
		}),
		agent:         agentInst,
		tools:         tools,
		promptBuilder: promptBuilder,
	}
}

func (t *SubAgentTool) Execute(ctx context.Context, call tool.Call) (tool.Result, error) {
	task, _ := call.Arguments["task"].(string)
	if task == "" {
		return tool.Result{}, errors.New("missing task argument")
	}

	promptSnapshot := t.promptBuilder()

	conv := conversation.NewConversation()
	events := t.agent.Run(ctx, agent.RunInput{
		SessionID:      call.ID,
		UserPrompt:     task,
		Conversation:   &conv,
		PromptSnapshot: promptSnapshot,
		Tools:          t.tools,
	})

	var result tool.Result
	var hasResult bool

	for ev := range events {
		switch e := ev.(type) {
		case agent.FinalAnswerEvent:
			result = tool.Result{
				CallID:  call.ID,
				Name:    call.Name,
				Content: conversation.Text(e.Message),
			}
			hasResult = true
		case agent.RunFailedEvent:
			slog.ErrorContext(ctx, "Sub-agent failed", "subagent", call.Name, "error", e.Err)
			return tool.Result{}, e.Err
		case agent.RunCompletedEvent:
			slog.InfoContext(ctx, "Sub-agent run completed",
				"subagent", call.Name,
				"llm_calls", e.Stats.LLMCalls,
				"input_tokens", e.Stats.TokenUsage.InputTokens,
				"output_tokens", e.Stats.TokenUsage.OutputTokens,
			)
		}
	}
	
	if !hasResult {
		return tool.Result{}, errors.New("sub-agent ended without final answer")
	}
	return result, nil
}
