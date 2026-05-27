package react

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent/loop"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
)

func (a *ReActAgent) run(ctx context.Context, events chan<- agent.Event, session *runtime.Session, userPrompt string) {
	defer close(events)

	stats := agent.NewRunStats()
	fail := func(err error) {
		stats.Finish()
		a.logRunStats(ctx, session.ID, stats)
		loop.Emit(ctx, events, agent.RunFailedEvent{
			SessionID: session.ID,
			Err:       err,
			Stats:     *stats,
		})
	}
	complete := func(reason agent.TerminalReason) {
		stats.Finish()
		a.logRunStats(ctx, session.ID, stats)
		loop.Emit(ctx, events, agent.RunCompletedEvent{
			SessionID: session.ID,
			Reason:    reason,
			Stats:     *stats,
		})
	}

	if !loop.Emit(ctx, events, agent.RunStartedEvent{SessionID: session.ID, MaxRounds: a.maxRounds}) {
		return
	}

	session.Conversation.Append(conversation.NewUserMessage(userPrompt))
	a.trace(ctx, "react.run.started", "Agent started", "session_id", session.ID, "max_rounds", a.maxRounds)

	for round := 1; round <= a.maxRounds; round++ {
		stats.Rounds = round
		messages := a.promptBuilder.BuildAgentMessages(nil, session.Conversation.Messages())
		a.writePromptDebugFile(ctx, messages)
		tools := a.getToolDefinitions()
		if !loop.Emit(ctx, events, agent.RoundStartedEvent{
			SessionID:    session.ID,
			Round:        round,
			MessageCount: len(messages),
			ToolCount:    len(tools),
		}) {
			return
		}
		a.trace(
			ctx,
			"react.round.generating",
			fmt.Sprintf("Round %d: Think", round),
			"session_id", session.ID,
			"message_count", len(messages),
			"tool_count", len(tools),
			"round", round,
		)

		if !loop.Emit(ctx, events, agent.ModelCallStartedEvent{SessionID: session.ID, Round: round}) {
			return
		}
		llmStartedAt := time.Now()
		stats.LLMCalls++
		response, textStreamed, err := a.generate(ctx, events, round, llm.GenerateRequest{
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
			fail(err)
			return
		}
		stats.InputTokens += response.Usage.InputTokens
		stats.OutputTokens += response.Usage.OutputTokens

		session.Conversation.Append(response.Message)
		toolCalls := loop.AssistantToolCalls(response.Message)
		a.trace(
			ctx,
			"react.round.generated",
			fmt.Sprintf("Round %d: model response generated", round),
			"session_id", session.ID,
			"round", round,
			"tool_calls", len(toolCalls),
			"tool_names", strings.Join(loop.ToolCallNames(toolCalls), ","),
			"input_tokens", response.Usage.InputTokens,
			"output_tokens", response.Usage.OutputTokens,
			"duration", llmDuration.String(),
		)
		if !loop.Emit(ctx, events, agent.ModelCallCompletedEvent{
			SessionID:     session.ID,
			Round:         round,
			Usage:         response.Usage,
			Duration:      llmDuration,
			ToolCallNames: loop.ToolCallNames(toolCalls),
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
			a.trace(
				ctx,
				"react.round.answer",
				fmt.Sprintf("Round %d: Answer", round),
				"session_id", session.ID,
				"round", round,
				"answer_chars", len(conversation.Text(response.Message)),
			)
			if !loop.Emit(ctx, events, agent.FinalAnswerEvent{
				SessionID: session.ID,
				Round:     round,
				Message:   assistantMessage,
			}) {
				return
			}
			a.trace(ctx, "react.run.completed", "Agent completed", "session_id", session.ID, "round", round)
			complete(agent.TerminalCompleted)
			return
		}

		a.trace(
			ctx,
			"react.round.act",
			fmt.Sprintf("Round %d: Act", round),
			"session_id", session.ID,
			"round", round,
			"decision_summary", decisionSummary(response.Message),
			"tool_calls", len(toolCalls),
			"tool_names", strings.Join(loop.ToolCallNames(toolCalls), ","),
		)

		if ok := a.executeToolCalls(ctx, events, session, round, toolCalls, stats); !ok {
			return
		}
	}

	a.trace(
		ctx,
		"react.run.max_rounds_reached",
		"Agent reached maximum tool rounds",
		"session_id", session.ID,
		"max_rounds", a.maxRounds,
	)
	complete(agent.TerminalMaxRounds)
}
