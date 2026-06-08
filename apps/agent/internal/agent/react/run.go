package react

import (
	"context"
	"errors"
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
		agent.Emit(ctx, events, agent.RunFailedEvent{
			SessionID: input.SessionID,
			Err:       err,
			Stats:     *stats,
		})
	}
	abort := func() {
		input.Conversation.Append(conversation.NewAssistantMessage("Run interrupted by user.", nil))
		runState.Finish()
		agent.EmitTerminal(events, agent.RunCompletedEvent{
			SessionID: input.SessionID,
			Reason:    agent.TerminalAborted,
			Stats:     *stats,
		})
	}
	complete := func(reason agent.TerminalReason) {
		runState.Finish()
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
	if ctx.Err() != nil {
		abort()
		return
	}

	sysTokens := input.PromptSnapshot.EstimateTokens()

	for round := 1; round <= runState.MaxRounds; round++ {
		runState.Round = round
		stats.Rounds = round

		maxTokens := 0
		if a.model != nil {
			maxTokens = a.model.MaxContextTokens()
		}

		maxHistoryTokens := maxTokens - sysTokens
		if ok := a.prepareContext(ctx, events, input, maxHistoryTokens, stats, round == 1); !ok {
			if ctx.Err() != nil {
				abort()
			}
			return
		}

		messages := input.PromptSnapshot.BuildMessages(nil, input.Conversation.Messages())
		currentTokens := conversation.EstimateTokens(messages)
		tools := input.ToolDefinitions()
		if !agent.Emit(ctx, events, agent.RoundStartedEvent{
			SessionID:            input.SessionID,
			Round:                round,
			MessageCount:         len(messages),
			ToolCount:            len(tools),
			CurrentContextTokens: currentTokens,
			MaxContextTokens:     maxTokens,
		}) {
			return
		}

		if !agent.Emit(ctx, events, agent.ModelCallStartedEvent{SessionID: input.SessionID, Round: round}) {
			return
		}
		llmStartedAt := time.Now()
		stats.LLMCalls++
		response, textStreamed, err := agent.GenerateWithStream(ctx, a.model, llm.GenerateRequest{
			SessionID: input.SessionID,
			Messages:  messages,
			Tools:     tools,
		}, events, round)
		llmDuration := time.Since(llmStartedAt)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				abort()
				return
			}
			fail(err)
			return
		}
		stats.TokenUsage.Add(response.Usage)

		input.Conversation.Append(response.Message)
		toolCalls := agent.AssistantToolCalls(response.Message)
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
			if !agent.Emit(ctx, events, agent.FinalAnswerEvent{
				SessionID: input.SessionID,
				Round:     round,
				Message:   assistantMessage,
			}) {
				return
			}
			complete(agent.TerminalCompleted)
			return
		}

		if ok := a.executeToolCalls(ctx, events, input, round, toolCalls, stats); !ok {
			if ctx.Err() != nil {
				abort()
			}
			return
		}
		if ctx.Err() != nil {
			abort()
			return
		}
	}

	complete(agent.TerminalMaxRounds)
}
