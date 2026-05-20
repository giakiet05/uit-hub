package conversation

type Role string

const (
	// RoleSystem marks instructions supplied by the runtime.
	RoleSystem Role = "system"
	// RoleUser marks messages written by the end user.
	RoleUser Role = "user"
	// RoleAssistant marks model-authored messages.
	RoleAssistant Role = "assistant"
	// RoleTool marks observations returned from tool execution.
	RoleTool Role = "tool"
)

// Message is the closed set of conversation message variants.
type Message interface {
	isMessage()
}

// SystemMessage contains runtime instructions sent to the LLM.
type SystemMessage struct {
	Content []ContentPart
}

// UserMessage contains a user prompt.
type UserMessage struct {
	Content []ContentPart
}

// AssistantMessage contains model text and optional tool calls.
type AssistantMessage struct {
	Content   []ContentPart
	ToolCalls []ToolCall
}

// ToolResultMessage contains the observation for a previously requested tool
// call.
type ToolResultMessage struct {
	ToolCallID string
	Content    []ContentPart
}

func (SystemMessage) isMessage()     {}
func (UserMessage) isMessage()       {}
func (AssistantMessage) isMessage()  {}
func (ToolResultMessage) isMessage() {}

// NewSystemMessage creates a system message from plain text.
func NewSystemMessage(text string) SystemMessage {
	return SystemMessage{Content: NewTextContent(text)}
}

// NewUserMessage creates a user message from plain text.
func NewUserMessage(text string) UserMessage {
	return UserMessage{Content: NewTextContent(text)}
}

// NewAssistantMessage creates an assistant message with optional tool calls.
func NewAssistantMessage(text string, toolCalls []ToolCall) AssistantMessage {
	return AssistantMessage{
		Content:   NewTextContent(text),
		ToolCalls: toolCalls,
	}
}

// NewToolResultMessage creates a tool observation message linked to a tool call.
func NewToolResultMessage(callID string, content string) ToolResultMessage {
	return ToolResultMessage{
		ToolCallID: callID,
		Content:    NewTextContent(content),
	}
}

// RoleOf returns the protocol role for a message variant.
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

// Text extracts all text parts from a message.
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

// contentText joins text content parts and ignores unsupported part types.
func contentText(content []ContentPart) string {
	var text string
	for _, part := range content {
		if part.Type == ContentTypeText {
			text += part.Text
		}
	}
	return text
}
