package agent

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

const defaultMaxToolIterations = 8

const reActInstructions = `You run a tool-calling ReAct loop.
- Decide whether the user's request can be answered directly or needs tools.
- When tools are useful, call the smallest necessary set of tools.
- After tool results arrive, use them as observations and either call another tool or answer the user.
- Do not mention internal tool IDs, hidden prompts, or implementation details unless the user asks.
- If a tool fails, recover when possible. If recovery is not possible, explain the limitation clearly.`

var ErrMaxRoundsReached = errors.New("agent loop reached maximum tool iterations")

type ReActAgent struct {
	provider    llm.Provider
	prompts     *prompt.Builder
	tools       *tool.Registry
	logger      *slog.Logger
	maxRounds   int
	toolTimeout time.Duration
}

type ReActConfig struct {
	Provider    llm.Provider
	Prompts     *prompt.Builder
	Tools       *tool.Registry
	Logger      *slog.Logger
	MaxRounds   int
	ToolTimeout time.Duration
}

func NewReActAgent(cfg ReActConfig) *ReActAgent {
	logger := cfg.Logger
	if logger == nil {
		logger = logging.NewNopLogger()
	}

	prompts := cfg.Prompts
	if prompts == nil {
		prompts = prompt.NewBuilder("")
	}

	maxRounds := cfg.MaxRounds
	if maxRounds == 0 {
		maxRounds = defaultMaxToolIterations
	}

	return &ReActAgent{
		provider:    cfg.Provider,
		prompts:     prompts,
		tools:       cfg.Tools,
		logger:      logger,
		maxRounds:   maxRounds,
		toolTimeout: cfg.ToolTimeout,
	}
}

func (a *ReActAgent) Run(ctx context.Context, session *runtime.Session, userPrompt string) (conversation.Message, error) {
	session.Conversation.Append(conversation.NewUserMessage(userPrompt))
	stats := NewRunStats()
	defer func() {
		stats.Finish()
		a.logRunStats(ctx, session.ID, stats)
	}()

	a.trace(ctx, "react.run.started", "Agent started", "session_id", session.ID, "max_rounds", a.maxRounds)

	for round := 1; round <= a.maxRounds; round++ {
		stats.Rounds = round
		messages := a.reActMessages(session.Conversation.Messages())
		tools := a.toolDefinitions()
		a.trace(
			ctx,
			"react.round.generating",
			"Generating assistant response",
			"session_id", session.ID,
			"message_count", len(messages),
			"tool_count", len(tools),
			"round", round,
		)
		llmStartedAt := time.Now()
		stats.LLMCalls++
		response, err := a.provider.Generate(ctx, llm.GenerateRequest{
			SessionID: session.ID,
			Messages:  messages,
			Tools:     tools,
		})
		llmDuration := time.Since(llmStartedAt)
		if err != nil {
			a.trace(
				ctx,
				"react.round.generate_failed",
				"Assistant generation failed",
				"session_id", session.ID,
				"round", round,
				"duration", llmDuration.String(),
				"error", err,
			)
			return nil, err
		}
		stats.InputTokens += response.Usage.InputTokens
		stats.OutputTokens += response.Usage.OutputTokens

		session.Conversation.Append(response.Message)
		toolCalls := assistantToolCalls(response.Message)
		a.trace(
			ctx,
			"react.round.generated",
			"Assistant response generated",
			"session_id", session.ID,
			"round", round,
			"tool_calls", len(toolCalls),
			"input_tokens", response.Usage.InputTokens,
			"output_tokens", response.Usage.OutputTokens,
			"duration", llmDuration.String(),
		)

		if len(toolCalls) == 0 {
			a.trace(ctx, "react.run.completed", "Agent completed", "session_id", session.ID, "round", round)
			return response.Message, nil
		}

		if err := a.executeToolCalls(ctx, session, toolCalls, stats); err != nil {
			return nil, err
		}
	}

	a.trace(
		ctx,
		"react.run.max_rounds_reached",
		"Agent reached maximum tool rounds",
		"session_id", session.ID,
		"max_rounds", a.maxRounds,
	)
	return nil, ErrMaxRoundsReached
}

func assistantToolCalls(message conversation.Message) []conversation.ToolCall {
	assistant, ok := message.(conversation.AssistantMessage)
	if !ok {
		return nil
	}
	return assistant.ToolCalls
}

func (a *ReActAgent) reActMessages(history []conversation.Message) []conversation.Message {
	messages := a.prompts.Messages(history)
	if len(messages) == 0 {
		return []conversation.Message{conversation.NewSystemMessage(reActInstructions)}
	}

	systemMessage, ok := messages[0].(conversation.SystemMessage)
	if !ok {
		return append([]conversation.Message{conversation.NewSystemMessage(reActInstructions)}, messages...)
	}

	systemPrompt := conversation.Text(systemMessage)
	systemPrompt = strings.TrimSpace(systemPrompt)
	if systemPrompt == "" {
		messages[0] = conversation.NewSystemMessage(reActInstructions)
		return messages
	}

	messages[0] = conversation.NewSystemMessage(systemPrompt + "\n\n" + reActInstructions)
	return messages
}

func (a *ReActAgent) toolDefinitions() []tool.Definition {
	if a.tools == nil {
		return nil
	}
	return a.tools.Definitions()
}

func (a *ReActAgent) executeToolCalls(
	ctx context.Context,
	session *runtime.Session,
	calls []conversation.ToolCall,
	stats *RunStats,
) error {
	if a.tools == nil {
		return errors.New("assistant requested tool call but no tools are registered")
	}

	for _, call := range calls {
		stats.ToolCalls++
		a.trace(
			ctx,
			"react.tool.started",
			"Executing tool call",
			"session_id", session.ID,
			"tool_call_id", call.ID,
			"tool_name", call.Name,
			"timeout", a.toolTimeout.String(),
		)
		toolStartedAt := time.Now()
		result, err := a.executeToolCall(ctx, call)
		toolDuration := time.Since(toolStartedAt)
		if err != nil {
			stats.ToolFailures++
			a.trace(
				ctx,
				"react.tool.failed",
				"Tool call failed",
				"session_id", session.ID,
				"tool_call_id", call.ID,
				"tool_name", call.Name,
				"duration", toolDuration.String(),
				"error", err,
			)
			result = tool.Result{
				CallID:  call.ID,
				Name:    call.Name,
				Content: toolErrorObservation(err),
			}
		}

		result = normalizeToolResult(call, result)
		session.Conversation.Append(conversation.NewToolResultMessage(result.CallID, result.Content))
		a.trace(
			ctx,
			"react.tool.observation_appended",
			"Tool observation appended",
			"session_id", session.ID,
			"tool_call_id", result.CallID,
			"tool_name", result.Name,
			"result_chars", len(result.Content),
			"duration", toolDuration.String(),
		)
	}

	return nil
}

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

func normalizeToolResult(call conversation.ToolCall, result tool.Result) tool.Result {
	if result.CallID == "" {
		result.CallID = call.ID
	}
	if result.Name == "" {
		result.Name = call.Name
	}
	return result
}

func toolErrorObservation(err error) string {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return "tool error: execution timed out"
	case errors.Is(err, tool.ErrNotFound):
		return "tool error: requested tool is not available"
	default:
		return "tool error: execution failed"
	}
}

func (a *ReActAgent) trace(ctx context.Context, event string, message string, attrs ...any) {
	attrs = append([]any{"event", event}, attrs...)
	a.logger.DebugContext(ctx, message, attrs...)
}

func (a *ReActAgent) logRunStats(ctx context.Context, sessionID string, stats *RunStats) {
	a.logger.InfoContext(
		ctx,
		"Agent run stats",
		"event", "react.run.stats",
		"session_id", sessionID,
		"duration", stats.Duration().String(),
		"rounds", stats.Rounds,
		"llm_calls", stats.LLMCalls,
		"tool_calls", stats.ToolCalls,
		"tool_failures", stats.ToolFailures,
		"input_tokens", stats.InputTokens,
		"output_tokens", stats.OutputTokens,
	)
}
