package llm

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

type Message struct {
	Role       Role
	Content    []ContentPart
	ToolCalls  []ToolCall
	ToolCallID string
	ToolName   string
}

func NewTextMessage(role Role, text string) Message {
	return Message{
		Role: role,
		Content: []ContentPart{
			{Type: ContentTypeText, Text: text},
		},
	}
}

func NewToolResultMessage(callID string, name string, content string) Message {
	return Message{
		Role:       RoleTool,
		ToolCallID: callID,
		ToolName:   name,
		Content: []ContentPart{
			{Type: ContentTypeText, Text: content},
		},
	}
}

func (m Message) ContentText() string {
	var text string
	for _, part := range m.Content {
		if part.Type == ContentTypeText {
			text += part.Text
		}
	}
	return text
}

type ContentType string

const (
	ContentTypeText ContentType = "text"
)

type ContentPart struct {
	Type ContentType
	Text string
}

type ToolDefinition struct {
	Name        string
	Description string
	InputSchema map[string]any
}

type ToolCall struct {
	ID        string
	Name      string
	Arguments map[string]any
}

type GenerateRequest struct {
	SessionID string
	Model     string
	Messages  []Message
	Tools     []ToolDefinition
}

type FinishReason string

const (
	FinishReasonStop     FinishReason = "stop"
	FinishReasonToolCall FinishReason = "tool_call"
)

type GenerateResponse struct {
	Message      Message
	FinishReason FinishReason
	Usage        Usage
}

type Usage struct {
	InputTokens  int
	OutputTokens int
}
