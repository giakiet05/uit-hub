package react

import (
	"context"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
)

const (
	defaultReactMicrocompactTriggerRatio = 0.70
	defaultReactSnipTriggerRatio         = 0.95
)

func (a *ReActAgent) prepareContext(
	ctx context.Context,
	events chan<- agent.Event,
	input agent.RunInput,
	maxHistoryTokens int,
	stats *agent.RunStats,
	allowContextCollapse bool,
) bool {
	if maxHistoryTokens <= 0 {
		return true
	}

	messages := input.Conversation.Messages()
	info := ""

	microThreshold := int(float64(maxHistoryTokens) * reactRatioOrDefault(
		a.compaction.MicrocompactTriggerRatio,
		defaultReactMicrocompactTriggerRatio,
	))
	currentTokens := conversation.EstimateTokens(messages)
	if !a.compaction.DisableMicrocompact && currentTokens > microThreshold {
		beforeTokens := currentTokens
		messages = conversation.Microcompact(
			messages,
			microThreshold,
			a.compaction.MicrocompactMinChars,
			a.compaction.MicrocompactKeepRecentToolResults,
		)
		currentTokens = conversation.EstimateTokens(messages)
		if currentTokens < beforeTokens {
			info += fmt.Sprintf("[Microcompact] Truncated tool results (%d -> %d tokens).\n", beforeTokens, currentTokens)
		}
	}

	if allowContextCollapse {
		collapseResult, err := a.collapseContext(ctx, input.SessionID, messages, maxHistoryTokens)
		if err != nil {
			info += fmt.Sprintf("[ContextCollapse] Failed: %v.\n", err)
		} else if collapseResult.Info != "" {
			messages = collapseResult.Messages
			info += collapseResult.Info
			stats.LLMCalls += collapseResult.LLMCalls
			stats.TokenUsage.Add(collapseResult.Usage)
		}
	}

	snipThreshold := int(float64(maxHistoryTokens) * reactRatioOrDefault(
		a.compaction.SnipTriggerRatio,
		defaultReactSnipTriggerRatio,
	))
	currentTokens = conversation.EstimateTokens(messages)
	if !a.compaction.DisableSnip && currentTokens > snipThreshold {
		beforeTokens := currentTokens
		messages = conversation.SnipCompact(messages, snipThreshold)
		currentTokens = conversation.EstimateTokens(messages)
		if currentTokens < beforeTokens {
			info += fmt.Sprintf("[SnipCompact] Removed old messages (%d -> %d tokens).\n", beforeTokens, currentTokens)
		}
	}

	input.Conversation.SetMessages(messages)
	if info == "" {
		return true
	}

	return agent.Emit(ctx, events, agent.CompactionTriggeredEvent{
		SessionID: input.SessionID,
		Info:      info,
	})
}

func reactRatioOrDefault(value float64, fallback float64) float64 {
	if value <= 0 || value > 1 {
		return fallback
	}
	return value
}
