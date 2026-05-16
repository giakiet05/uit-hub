package conversation

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

type Message interface {
	isMessage()
}

type SystemMessage struct {
	Content []ContentPart
}

type UserMessage struct {
	Content []ContentPart
}

type AssistantMessage struct {
	Content   []ContentPart
	ToolCalls []ToolCall
}

type ToolResultMessage struct {
	ToolCallID string
	Content    []ContentPart
}

func (SystemMessage) isMessage()     {}
func (UserMessage) isMessage()       {}
func (AssistantMessage) isMessage()  {}
func (ToolResultMessage) isMessage() {}

func NewSystemMessage(text string) SystemMessage {
	return SystemMessage{Content: NewTextContent(text)}
}

func NewUserMessage(text string) UserMessage {
	return UserMessage{Content: NewTextContent(text)}
}

func NewAssistantMessage(text string, toolCalls []ToolCall) AssistantMessage {
	return AssistantMessage{
		Content:   NewTextContent(text),
		ToolCalls: toolCalls,
	}
}

func NewToolResultMessage(callID string, content string) ToolResultMessage {
	return ToolResultMessage{
		ToolCallID: callID,
		Content:    NewTextContent(content),
	}
}

func RoleOf(message Message) Role {
	switch message.(type) {
	case SystemMessage:
		return RoleSystem
	case UserMessage:
		return RoleUser
	case AssistantMessage:
		return RoleAssistant
	case ToolResultMessage:
		return RoleTool
	default:
		return ""
	}
}

func Text(message Message) string {
	switch typed := message.(type) {
	case SystemMessage:
		return contentText(typed.Content)
	case UserMessage:
		return contentText(typed.Content)
	case AssistantMessage:
		return contentText(typed.Content)
	case ToolResultMessage:
		return contentText(typed.Content)
	default:
		return ""
	}
}

func contentText(content []ContentPart) string {
	var text string
	for _, part := range content {
		if part.Type == ContentTypeText {
			text += part.Text
		}
	}
	return text
}
