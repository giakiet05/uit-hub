package prompt

import (
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
)

const (
	dynamicContextHeader  = "=== DYNAMIC CONTEXT ==="
	uncachedContextHeader = "=== UNCACHED CONTEXT ==="
)

// StaticPart is prompt text that should stay stable for every runtime session.
type StaticPart string

// DynamicPart is session-scoped prompt text loaded at session start.
type DynamicPart string

// UncachedPart is volatile prompt text built for one LLM call.
type UncachedPart string

// SystemPrompt contains prompt parts stable for one runtime session.
type SystemPrompt struct {
	StaticParts  []StaticPart
	DynamicParts []DynamicPart
}

// Render converts the prompt into the final system prompt text.
func (p SystemPrompt) Render(uncached []UncachedPart) string {
	parts := make([]string, 0, 3)

	if text := renderStaticParts(p.StaticParts); text != "" {
		parts = append(parts, text)
	}
	if text := renderDynamicParts(p.DynamicParts); text != "" {
		parts = append(parts, dynamicContextHeader+"\n\n"+text)
	}
	if text := renderUncachedParts(uncached); text != "" {
		parts = append(parts, uncachedContextHeader+"\n\n"+text)
	}

	return strings.Join(parts, "\n\n")
}

// BuildMessages creates the system message and prepends it to the history array.
func (p SystemPrompt) BuildMessages(uncached []UncachedPart, history []conversation.Message) []conversation.Message {
	systemText := p.Render(uncached)

	messages := make([]conversation.Message, 0, len(history)+1)
	if systemText != "" {
		messages = append(messages, conversation.NewSystemMessage(systemText))
	}
	messages = append(messages, history...)
	return messages
}

// EstimateTokens calculates the token size of the stable system prompt parts.
func (p SystemPrompt) EstimateTokens() int {
	text := p.Render(nil)
	if text == "" {
		return 0
	}
	msg := conversation.NewSystemMessage(text)
	return conversation.EstimateTokens([]conversation.Message{msg})
}

// renderStaticParts joins static prompt parts while skipping empty parts.
func renderStaticParts(parts []StaticPart) string {
	texts := make([]string, 0, len(parts))
	for _, part := range parts {
		if text := strings.TrimSpace(string(part)); text != "" {
			texts = append(texts, text)
		}
	}
	return strings.Join(texts, "\n\n")
}

// renderDynamicParts joins dynamic prompt parts while skipping empty parts.
func renderDynamicParts(parts []DynamicPart) string {
	texts := make([]string, 0, len(parts))
	for _, part := range parts {
		if text := strings.TrimSpace(string(part)); text != "" {
			texts = append(texts, text)
		}
	}
	return strings.Join(texts, "\n\n")
}

// renderUncachedParts joins volatile prompt parts while skipping empty parts.
func renderUncachedParts(parts []UncachedPart) string {
	texts := make([]string, 0, len(parts))
	for _, part := range parts {
		if text := strings.TrimSpace(string(part)); text != "" {
			texts = append(texts, text)
		}
	}
	return strings.Join(texts, "\n\n")
}
