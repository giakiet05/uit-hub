package agent

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

const maxToolIterations = 8

type Loop struct {
	provider  llm.Provider
	prompts   *prompt.Builder
	tools     *tool.Registry
	logger    *slog.Logger
	maxRounds int
}

func NewLoop(provider llm.Provider, prompts *prompt.Builder, logger *slog.Logger, registries ...*tool.Registry) *Loop {
	if logger == nil {
		logger = logging.NewNopLogger()
	}

	var registry *tool.Registry
	if len(registries) > 0 {
		registry = registries[0]
	}

	return &Loop{
		provider:  provider,
		prompts:   prompts,
		tools:     registry,
		logger:    logger,
		maxRounds: maxToolIterations,
	}
}

func (l *Loop) Run(ctx context.Context, session *runtime.Session, userPrompt string) (llm.Message, error) {
	session.Append(llm.NewTextMessage(llm.RoleUser, userPrompt))

	for round := 1; round <= l.maxRounds; round++ {
		messages := l.prompts.Messages(session.Messages())
		tools := l.toolDefinitions()
		l.logger.DebugContext(ctx, "Generating assistant response", "session_id", session.ID, "message_count", len(messages), "tool_count", len(tools), "round", round)
		response, err := l.provider.Generate(ctx, llm.GenerateRequest{
			SessionID: session.ID,
			Messages:  messages,
			Tools:     tools,
		})
		if err != nil {
			l.logger.DebugContext(ctx, "Provider generate failed", "session_id", session.ID, "round", round, "error", err)
			return llm.Message{}, err
		}

		session.Append(response.Message)
		l.logger.DebugContext(ctx, "Assistant response generated", "session_id", session.ID, "finish_reason", response.FinishReason, "tool_calls", len(response.Message.ToolCalls), "input_tokens", response.Usage.InputTokens, "output_tokens", response.Usage.OutputTokens)

		if len(response.Message.ToolCalls) == 0 {
			return response.Message, nil
		}

		if err := l.executeToolCalls(ctx, session, response.Message.ToolCalls); err != nil {
			return llm.Message{}, err
		}
	}

	return llm.Message{}, errors.New("agent loop reached maximum tool iterations")
}

func (l *Loop) toolDefinitions() []llm.ToolDefinition {
	if l.tools == nil {
		return nil
	}
	return l.tools.Definitions()
}

func (l *Loop) executeToolCalls(ctx context.Context, session *runtime.Session, calls []llm.ToolCall) error {
	if l.tools == nil {
		return errors.New("assistant requested tool call but no tools are registered")
	}

	for _, call := range calls {
		l.logger.DebugContext(ctx, "Executing tool call", "session_id", session.ID, "tool_call_id", call.ID, "tool_name", call.Name)
		result, err := l.tools.Execute(ctx, tool.Call{
			ID:        call.ID,
			Name:      call.Name,
			Arguments: call.Arguments,
		})
		if err != nil {
			l.logger.DebugContext(ctx, "Tool call failed", "session_id", session.ID, "tool_call_id", call.ID, "tool_name", call.Name, "error", err)
			result = tool.Result{
				CallID:  call.ID,
				Name:    call.Name,
				Content: fmt.Sprintf("tool error: %v", err),
			}
		}

		session.Append(llm.NewToolResultMessage(result.CallID, result.Name, result.Content))
		l.logger.DebugContext(ctx, "Tool call result appended", "session_id", session.ID, "tool_call_id", result.CallID, "tool_name", result.Name, "result_chars", len(result.Content))
	}

	return nil
}
