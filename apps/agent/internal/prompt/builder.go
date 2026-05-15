package prompt

import "github.com/giakiet05/uit-hub/apps/agent/internal/llm"

type Builder struct {
	systemPrompt string
}

func NewBuilder(systemPrompt string) *Builder {
	return &Builder{systemPrompt: systemPrompt}
}

func (b *Builder) Messages(conversation []llm.Message) []llm.Message {
	messages := make([]llm.Message, 0, len(conversation)+1)
	if b.systemPrompt != "" {
		messages = append(messages, llm.NewTextMessage(llm.RoleSystem, b.systemPrompt))
	}
	messages = append(messages, conversation...)
	return messages
}
