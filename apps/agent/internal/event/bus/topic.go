package bus

// Topic represents an event routing category.
type Topic string

const (
	// TopicStream is for high-frequency events like text stream deltas.
	TopicStream Topic = "stream"
	// TopicLifecycle is for low-frequency structural events like tool calls and completions.
	TopicLifecycle Topic = "lifecycle"
	// TopicAll is a special topic to receive all events.
	TopicAll Topic = "*"
)

// Event is the base interface that any published message must implement.
type Event interface {
	Topic() Topic
}
