package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
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
const logPreviewLimit = 160

var ErrMaxRoundsReached = errors.New("agent loop reached maximum tool iterations")

// ReActAgent implements the baseline reason-act-observe loop for tool-calling
// agents.
type ReActAgent struct {
	provider    llm.Provider
	prompts     *prompt.Builder
	tools       *tool.Registry
	logger      *slog.Logger
	maxRounds   int
	toolTimeout time.Duration
}

// ReActConfig groups the dependencies and runtime limits required to build a
// ReActAgent.
type ReActConfig struct {
	Provider    llm.Provider
	Prompts     *prompt.Builder
	Tools       *tool.Registry
	Logger      *slog.Logger
	MaxRounds   int
	ToolTimeout time.Duration
}

// NewReActAgent constructs a ReActAgent and fills optional dependencies with
// conservative defaults.
func NewReActAgent(cfg ReActConfig) *ReActAgent {
	logger := cfg.Logger
	if logger == nil {
		logger = logging.NewNopLogger()
	}

	prompts := cfg.Prompts
	if prompts == nil {
		prompts = prompt.NewBuilder(prompt.SystemPrompt{})
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

// Run appends the user prompt to the session and loops until the model returns
// a final assistant message or the configured round limit is reached.
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
		messages := a.prompts.BuildAgentMessages(TypeReAct.String(), session.Conversation.Messages())
		a.writePromptDebugFile(ctx, messages)
		tools := a.getToolDefinitions()
		a.trace(
			ctx,
			"react.round.generating",
			fmt.Sprintf("Round %d: Think", round),
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
		toolCalls := getAssistantToolCalls(response.Message)
		a.trace(
			ctx,
			"react.round.generated",
			fmt.Sprintf("Round %d: model response generated", round),
			"session_id", session.ID,
			"round", round,
			"tool_calls", len(toolCalls),
			"tool_names", strings.Join(toolCallNames(toolCalls), ","),
			"input_tokens", response.Usage.InputTokens,
			"output_tokens", response.Usage.OutputTokens,
			"duration", llmDuration.String(),
		)

		if len(toolCalls) == 0 {
			a.trace(
				ctx,
				"react.round.answer",
				fmt.Sprintf("Round %d: Answer", round),
				"session_id", session.ID,
				"round", round,
				"answer_chars", len(conversation.Text(response.Message)),
			)
			a.trace(ctx, "react.run.completed", "Agent completed", "session_id", session.ID, "round", round)
			return response.Message, nil
		}

		a.trace(
			ctx,
			"react.round.act",
			fmt.Sprintf("Round %d: Act", round),
			"session_id", session.ID,
			"round", round,
			"decision_summary", decisionSummary(response.Message),
			"tool_calls", len(toolCalls),
			"tool_names", strings.Join(toolCallNames(toolCalls), ","),
		)

		if err := a.executeToolCalls(ctx, session, round, toolCalls, stats); err != nil {
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

// writePromptDebugFile writes the latest LLM message stack for local prompt
// inspection without exposing it through the TUI log stream.
func (a *ReActAgent) writePromptDebugFile(ctx context.Context, messages []conversation.Message) {
	if err := os.MkdirAll("tmp", 0o755); err != nil {
		a.trace(ctx, "react.prompt_debug.write_failed", "Prompt debug directory creation failed", "error", err)
		return
	}
	if err := os.WriteFile(filepath.Join("tmp", "prompt.md"), []byte(renderPromptDebug(messages)), 0o600); err != nil {
		a.trace(ctx, "react.prompt_debug.write_failed", "Prompt debug file write failed", "error", err)
	}
}

func renderPromptDebug(messages []conversation.Message) string {
	parts := make([]string, 0, len(messages))
	for _, message := range messages {
		role := conversation.RoleOf(message)
		text := conversation.Text(message)
		if text == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("## %s\n\n%s", role, text))
	}
	return strings.Join(parts, "\n\n")
}

// getAssistantToolCalls extracts tool calls from an assistant message.
func getAssistantToolCalls(message conversation.Message) []conversation.ToolCall {
	assistant, ok := message.(conversation.AssistantMessage)
	if !ok {
		return nil
	}
	return assistant.ToolCalls
}

// getToolDefinitions returns the current tool schemas, or nil when the agent has
// no registry.
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
	session *runtime.Session,
	round int,
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
			fmt.Sprintf("Round %d: Tool %s started", round, call.Name),
			"session_id", session.ID,
			"round", round,
			"tool_call_id", call.ID,
			"tool_name", call.Name,
			"arguments_preview", previewValue(call.Arguments),
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
				Content: toolErrorObservation(err),
			}
		}

		result = normalizeToolResult(call, result)
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
			"result_preview", previewText(result.Content),
			"duration", toolDuration.String(),
		)
	}

	return nil
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

// normalizeToolResult fills missing result metadata from the original tool call.
func normalizeToolResult(call conversation.ToolCall, result tool.Result) tool.Result {
	if result.CallID == "" {
		result.CallID = call.ID
	}
	if result.Name == "" {
		result.Name = call.Name
	}
	return result
}

// toolErrorObservation converts internal tool failures into sanitized
// observations that can be safely shown back to the model.
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

// toolCallNames returns tool names in call order for compact logging.
func toolCallNames(calls []conversation.ToolCall) []string {
	names := make([]string, 0, len(calls))
	for _, call := range calls {
		names = append(names, call.Name)
	}
	return names
}

// decisionSummary returns the assistant's visible tool-use summary for debug
// logging.
func decisionSummary(message conversation.Message) string {
	text := strings.TrimSpace(conversation.Text(message))
	if text == "" {
		return "empty"
	}
	return previewText(text)
}

// previewValue renders a structured value as a bounded single-line log preview.
func previewValue(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return previewText(fmt.Sprintf("%v", value))
	}
	return previewText(string(data))
}

// previewText returns a bounded single-line preview suitable for structured
// logs.
func previewText(text string) string {
	text = strings.ReplaceAll(text, "\n", "\\n")
	runes := []rune(text)
	if len(runes) <= logPreviewLimit {
		return text
	}
	return string(runes[:logPreviewLimit]) + "..."
}

// trace writes a debug event with a stable event key.
func (a *ReActAgent) trace(ctx context.Context, event string, message string, attrs ...any) {
	attrs = append([]any{"event", event}, attrs...)
	a.logger.DebugContext(ctx, message, attrs...)
}

// logRunStats writes the summary counters for one completed or failed run.
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
