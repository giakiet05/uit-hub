// Package prompt builds the message stack sent to LLM providers.
package prompt

import "github.com/giakiet05/uit-hub/apps/agent/internal/conversation"

// Builder constructs provider-ready message stacks from structured prompts.
type Builder struct {
	sessionPrompt SessionPrompt
}

// NewBuilder creates a builder from one session-stable prompt snapshot.
func NewBuilder(sessionPrompt SessionPrompt) *Builder {
	return &Builder{sessionPrompt: sessionPrompt}
}

// BuildAgentMessages builds messages for one LLM call.
func (b *Builder) BuildAgentMessages(
	uncached []UncachedPart,
	history []conversation.Message,
) []conversation.Message {
	callPrompt := b.BuildCallPrompt(uncached)
	systemText := callPrompt.Render()

	messages := make([]conversation.Message, 0, len(history)+1)
	if systemText != "" {
		messages = append(messages, conversation.NewSystemMessage(systemText))
	}
	messages = append(messages, history...)
	return messages
}

// BuildCallPrompt combines the session prompt with one-call volatile prompt
// parts.
func (b *Builder) BuildCallPrompt(uncached []UncachedPart) CallPrompt {
	return CallPrompt{
		Session:       b.sessionPrompt,
		UncachedParts: uncached,
	}
}
