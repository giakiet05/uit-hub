package conversation

// Conversation owns the ordered message history for one session.
type Conversation struct {
	messages []Message
}

// NewConversation creates an empty conversation history.
func NewConversation() Conversation {
	return Conversation{
		messages: []Message{},
	}
}

// Append adds a message to the end of the conversation.
func (c *Conversation) Append(message Message) {
	c.messages = append(c.messages, message)
}

// Messages returns a defensive copy of the message history.
func (c *Conversation) Messages() []Message {
	messages := make([]Message, len(c.messages))
	copy(messages, c.messages)
	return messages
}
