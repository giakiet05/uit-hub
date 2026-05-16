package prompt

import "github.com/giakiet05/uit-hub/apps/agent/internal/conversation"

type Builder struct {
	systemPrompt string
}

func NewBuilder(systemPrompt string) *Builder {
	return &Builder{systemPrompt: systemPrompt}
}

func (b *Builder) Messages(history []conversation.Message) []conversation.Message {
	messages := make([]conversation.Message, 0, len(history)+1)
	if b.systemPrompt != "" {
		messages = append(messages, conversation.NewSystemMessage(b.systemPrompt))
	}
	messages = append(messages, history...)
	return messages
}
