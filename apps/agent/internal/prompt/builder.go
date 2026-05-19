// Package prompt builds the message stack sent to LLM providers.
package prompt

import "github.com/giakiet05/uit-hub/apps/agent/internal/conversation"

const reActAgentType = "react"

// Builder constructs provider-ready message stacks from structured prompts.
type Builder struct {
	systemPrompt SystemPrompt
}

// NewBuilder creates a prompt builder from a structured system prompt.
func NewBuilder(systemPrompt SystemPrompt) *Builder {
	return &Builder{systemPrompt: systemPrompt}
}

// BuildAgentMessages builds the system message for an agent type and appends
// conversation history after it.
func (b *Builder) BuildAgentMessages(agentType string, history []conversation.Message) []conversation.Message {
	systemPrompt := b.BuildSystemPrompt(agentType)
	systemText := systemPrompt.Render()

	messages := make([]conversation.Message, 0, len(history)+1)
	if systemText != "" {
		messages = append(messages, conversation.NewSystemMessage(systemText))
	}
	messages = append(messages, history...)
	return messages
}

// BuildSystemPrompt merges agent-specific prompt with the builder-level prompt.
func (b *Builder) BuildSystemPrompt(agentType string) SystemPrompt {
	return systemPromptForAgent(agentType).Merge(b.systemPrompt)
}

func systemPromptForAgent(agentType string) SystemPrompt {
	switch agentType {
	case reActAgentType:
		return ReActSystemPrompt()
	default:
		return SystemPrompt{}
	}
}
