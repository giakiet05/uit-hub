package react

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/usage"
)

const (
	defaultContextCollapseTriggerRatio     = 0.80
	defaultContextCollapseTargetRatio      = 0.65
	defaultContextCollapseKeepRecentTurns  = 3
	defaultContextCollapseMaxSegments      = 3
	defaultContextCollapseMaxConcurrency   = 2
	defaultContextCollapseMinSegmentTokens = 1500
)

// ContextCollapseOptions configures LLM-backed collapse of old conversation
// turns before the main round model call.
type ContextCollapseOptions struct {
	Enabled          bool
	TriggerRatio     float64
	TargetRatio      float64
	KeepRecentTurns  int
	MaxSegments      int
	MaxConcurrency   int
	MinSegmentTokens int
}

type contextCollapseResult struct {
	Messages []conversation.Message
	Info     string
	Usage    usage.TokenUsage
	LLMCalls int
}

type collapseJobResult struct {
	index   int
	summary string
	usage   usage.TokenUsage
	err     error
}

func (a *ReActAgent) collapseContext(
	ctx context.Context,
	sessionID string,
	messages []conversation.Message,
	maxHistoryTokens int,
) (contextCollapseResult, error) {
	result := contextCollapseResult{Messages: messages}
	if !a.contextCollapse.Enabled || a.model == nil || maxHistoryTokens <= 0 {
		return result, nil
	}

	options := normalizeContextCollapseOptions(a.contextCollapse)
	currentTokens := conversation.EstimateTokens(messages)
	triggerTokens := int(float64(maxHistoryTokens) * options.TriggerRatio)
	if currentTokens <= triggerTokens {
		return result, nil
	}

	targetTokens := int(float64(maxHistoryTokens) * options.TargetRatio)
	segments := conversation.BuildCollapseSegments(messages, conversation.CollapseSegmentOptions{
		KeepRecentTurns:     options.KeepRecentTurns,
		MaxSegments:         options.MaxSegments,
		MinSegmentTokens:    options.MinSegmentTokens,
		TargetTokensToCover: currentTokens - targetTokens,
	})
	if len(segments) == 0 {
		return result, nil
	}

	summaries, collapseUsage, err := a.collapseSegments(ctx, sessionID, messages, segments, options.MaxConcurrency)
	if err != nil {
		return result, err
	}

	collapsed := conversation.ReplaceCollapseSegments(messages, segments, summaries)
	afterTokens := conversation.EstimateTokens(collapsed)
	result.Messages = collapsed
	result.Usage = collapseUsage
	result.LLMCalls = len(segments)
	result.Info = fmt.Sprintf(
		"[ContextCollapse] Summarized %d old segment(s) (%d -> %d tokens).\n",
		len(segments),
		currentTokens,
		afterTokens,
	)
	return result, nil
}

func normalizeContextCollapseOptions(options ContextCollapseOptions) ContextCollapseOptions {
	options.TriggerRatio = collapseRatioOrDefault(options.TriggerRatio, defaultContextCollapseTriggerRatio)
	options.TargetRatio = collapseRatioOrDefault(options.TargetRatio, defaultContextCollapseTargetRatio)
	if options.TargetRatio >= options.TriggerRatio {
		options.TargetRatio = defaultContextCollapseTargetRatio
	}
	if options.KeepRecentTurns < 0 {
		options.KeepRecentTurns = defaultContextCollapseKeepRecentTurns
	}
	if options.KeepRecentTurns == 0 {
		options.KeepRecentTurns = defaultContextCollapseKeepRecentTurns
	}
	if options.MaxSegments <= 0 {
		options.MaxSegments = defaultContextCollapseMaxSegments
	}
	if options.MaxConcurrency <= 0 {
		options.MaxConcurrency = defaultContextCollapseMaxConcurrency
	}
	if options.MinSegmentTokens <= 0 {
		options.MinSegmentTokens = defaultContextCollapseMinSegmentTokens
	}
	return options
}

func collapseRatioOrDefault(value float64, fallback float64) float64 {
	if value <= 0 || value > 1 {
		return fallback
	}
	return value
}

func (a *ReActAgent) collapseSegments(
	ctx context.Context,
	sessionID string,
	messages []conversation.Message,
	segments []conversation.CollapseSegment,
	maxConcurrency int,
) ([]string, usage.TokenUsage, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan collapseJobResult, len(segments))
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for index, segment := range segments {
		wg.Add(1)
		go func(index int, segment conversation.CollapseSegment) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				results <- collapseJobResult{index: index, err: ctx.Err()}
				return
			}

			summary, tokenUsage, err := a.summarizeSegment(ctx, sessionID, index, messages[segment.Start:segment.End])
			if err != nil {
				cancel()
			}
			results <- collapseJobResult{
				index:   index,
				summary: summary,
				usage:   tokenUsage,
				err:     err,
			}
		}(index, segment)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	summaries := make([]string, len(segments))
	var totalUsage usage.TokenUsage
	var joinedErr error
	for result := range results {
		if result.err != nil {
			joinedErr = errors.Join(joinedErr, result.err)
			continue
		}
		summaries[result.index] = result.summary
		totalUsage.Add(result.usage)
	}
	if joinedErr != nil {
		return nil, usage.TokenUsage{}, joinedErr
	}
	return summaries, totalUsage, nil
}

func (a *ReActAgent) summarizeSegment(
	ctx context.Context,
	sessionID string,
	segmentIndex int,
	messages []conversation.Message,
) (string, usage.TokenUsage, error) {
	response, err := a.model.Generate(ctx, llm.GenerateRequest{
		SessionID: sessionID,
		Messages: []conversation.Message{
			conversation.NewSystemMessage(contextCollapseSystemPrompt),
			conversation.NewUserMessage(formatCollapseSegment(segmentIndex, messages)),
		},
	})
	if err != nil {
		return "", usage.TokenUsage{}, err
	}

	summary := strings.TrimSpace(conversation.Text(response.Message))
	if summary == "" {
		return "", response.Usage, errors.New("context collapse produced an empty summary")
	}
	return summary, response.Usage, nil
}

const contextCollapseSystemPrompt = `You compress old conversation history for an agent.

Preserve durable user intent, confirmed facts, important tool observations, decisions, and unfinished tasks.
Do not invent facts.
Do not include unnecessary wording.
Return only the compact summary in Vietnamese unless the segment is clearly in another language.`

func formatCollapseSegment(segmentIndex int, messages []conversation.Message) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Compress old conversation segment %d:\n\n", segmentIndex+1))
	for _, message := range messages {
		builder.WriteString("Role: ")
		builder.WriteString(string(conversation.RoleOf(message)))
		builder.WriteString("\n")
		if text := strings.TrimSpace(conversation.Text(message)); text != "" {
			builder.WriteString("Text:\n")
			builder.WriteString(text)
			builder.WriteString("\n")
		}
		if assistant, ok := message.(conversation.AssistantMessage); ok && len(assistant.ToolCalls) > 0 {
			builder.WriteString("Tool calls:\n")
			for _, call := range assistant.ToolCalls {
				builder.WriteString("- ")
				builder.WriteString(call.Name)
				if call.ID != "" {
					builder.WriteString(" id=")
					builder.WriteString(call.ID)
				}
				builder.WriteString("\n")
			}
		}
		if toolResult, ok := message.(conversation.ToolResultMessage); ok && toolResult.ToolCallID != "" {
			builder.WriteString("Tool call ID: ")
			builder.WriteString(toolResult.ToolCallID)
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}
	return builder.String()
}
