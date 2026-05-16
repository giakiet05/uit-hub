package conversation

type ToolCall struct {
	ID        string
	Name      string
	Arguments map[string]any
}
