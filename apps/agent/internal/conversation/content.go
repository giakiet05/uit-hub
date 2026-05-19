// Package conversation defines the message model shared by agents and LLM
// providers.
package conversation

// ContentType identifies the kind of content stored in a message part.
type ContentType string

const (
	// ContentTypeText is plain text content.
	ContentTypeText ContentType = "text"
)

// ContentPart is one typed piece of message content.
type ContentPart struct {
	Type ContentType
	Text string
}

// NewTextContent wraps plain text in the canonical content-part slice.
func NewTextContent(text string) []ContentPart {
	return []ContentPart{
		{Type: ContentTypeText, Text: text},
	}
}
