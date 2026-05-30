package react

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
)

// decisionSummary returns the assistant's visible tool-use summary for debug logging.
func decisionSummary(message conversation.Message) string {
	text := strings.TrimSpace(conversation.Text(message))
	if text == "" {
		return "empty"
	}
	return agent.PreviewText(text)
}

func (a *ReActAgent) run(ctx context.Context, events chan<- agent.Event, input agent.RunInput) {
	defer close(events)

	runState := agent.NewRunState(input.SessionID, a.maxRounds)
	stats := runState.Stats
	fail := func(err error) {
		runState.Finish()
		agent.LogRunStats(a.logger, "react", ctx, input.SessionID, stats)
		agent.Emit(ctx, events, agent.RunFailedEvent{
			SessionID: input.SessionID,
			Err:       err,
			Stats:     *stats,
		})
	}
	complete := func(reason agent.TerminalReason) {
		runState.Finish()
		agent.LogRunStats(a.logger, "react", ctx, input.SessionID, stats)
		agent.Emit(ctx, events, agent.RunCompletedEvent{
			SessionID: input.SessionID,
			Reason:    reason,
			Stats:     *stats,
		})
	}

	if !agent.Emit(ctx, events, agent.RunStartedEvent{SessionID: input.SessionID, MaxRounds: runState.MaxRounds}) {
		return
	}

	input.Conversation.Append(conversation.NewUserMessage(input.UserPrompt))
	agent.Trace(a.logger, ctx, "react.run.started", "Agent started", "session_id", input.SessionID, "max_rounds", a.maxRounds)

	for round := 1; round <= runState.MaxRounds; round++ {
		runState.Round = round
		stats.Rounds = round
		messages := input.BuildAgentMessages(nil)
		agent.WritePromptDebugFile(ctx, a.logger, messages)
		tools := input.ToolDefinitions()
		if !agent.Emit(ctx, events, agent.RoundStartedEvent{
			SessionID:    input.SessionID,
			Round:        round,
			MessageCount: len(messages),
			ToolCount:    len(tools),
		}) {
			return
		}
		agent.Trace(a.logger, ctx,
			"react.round.generating",
			fmt.Sprintf("Round %d: Think", round),
			"session_id", input.SessionID,
			"message_count", len(messages),
			"tool_count", len(tools),
			"round", round,
		)

		if !agent.Emit(ctx, events, agent.ModelCallStartedEvent{SessionID: input.SessionID, Round: round}) {
			return
		}
		llmStartedAt := time.Now()
		stats.LLMCalls++
		response, textStreamed, err := a.generate(ctx, events, round, llm.GenerateRequest{
			SessionID: input.SessionID,
			Messages:  messages,
			Tools:     tools,
		})
		llmDuration := time.Since(llmStartedAt)
		if err != nil {
			agent.Trace(a.logger, ctx,
				"react.round.generate_failed",
				"Assistant generation failed",
				"session_id", input.SessionID,
				"round", round,
				"duration", llmDuration.String(),
				"error", err,
			)
			fail(err)
			return
		}
		stats.TokenUsage.Add(response.Usage)

		input.Conversation.Append(response.Message)
		toolCalls := agent.AssistantToolCalls(response.Message)
		agent.Trace(a.logger, ctx,
			"react.round.generated",
			fmt.Sprintf("Round %d: model response generated", round),
			"session_id", input.SessionID,
			"round", round,
			"tool_calls", len(toolCalls),
			"tool_names", strings.Join(agent.ToolCallNames(toolCalls), ","),
			"input_tokens", response.Usage.InputTokens,
			"output_tokens", response.Usage.OutputTokens,
			"duration", llmDuration.String(),
		)
		if !agent.Emit(ctx, events, agent.ModelCallCompletedEvent{
			SessionID:     input.SessionID,
			Round:         round,
			Usage:         response.Usage,
			Duration:      llmDuration,
			ToolCallNames: agent.ToolCallNames(toolCalls),
			MessageText:   conversation.Text(response.Message),
			TextStreamed:  textStreamed,
		}) {
			return
		}

		if len(toolCalls) == 0 {
			assistantMessage, ok := response.Message.(conversation.AssistantMessage)
			if !ok {
				fail(errors.New("llm response is not an assistant message"))
				return
			}
			agent.Trace(a.logger, ctx,
				"react.round.answer",
				fmt.Sprintf("Round %d: Answer", round),
				"session_id", input.SessionID,
				"round", round,
				"answer_chars", len(conversation.Text(response.Message)),
			)
			if !agent.Emit(ctx, events, agent.FinalAnswerEvent{
				SessionID: input.SessionID,
				Round:     round,
				Message:   assistantMessage,
			}) {
				return
			}
			agent.Trace(a.logger, ctx, "react.run.completed", "Agent completed", "session_id", input.SessionID, "round", round)
			complete(agent.TerminalCompleted)
			return
		}

		agent.Trace(a.logger, ctx,
			"react.round.act",
			fmt.Sprintf("Round %d: Act", round),
			"session_id", input.SessionID,
			"round", round,
			"decision_summary", decisionSummary(response.Message),
			"tool_calls", len(toolCalls),
			"tool_names", strings.Join(agent.ToolCallNames(toolCalls), ","),
		)

		if ok := a.executeToolCalls(ctx, events, input, round, toolCalls, stats); !ok {
			return
		}
	}

	agent.Trace(a.logger, ctx,
		"react.run.max_rounds_reached",
		"Agent reached maximum tool rounds",
		"session_id", input.SessionID,
		"max_rounds", a.maxRounds,
	)
	complete(agent.TerminalMaxRounds)
}
