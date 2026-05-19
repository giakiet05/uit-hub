package conversation

// ToolCall is the provider-independent representation of a model-requested
// tool invocation.
type ToolCall struct {
	ID        string
	Name      string
	Arguments map[string]any
}
