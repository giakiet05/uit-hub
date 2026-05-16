package conversation

type Conversation struct {
	messages []Message
}

func NewConversation() Conversation {
	return Conversation{
		messages: []Message{},
	}
}

func (c *Conversation) Append(message Message) {
	c.messages = append(c.messages, message)
}

func (c *Conversation) Messages() []Message {
	messages := make([]Message, len(c.messages))
	copy(messages, c.messages)
	return messages
}
