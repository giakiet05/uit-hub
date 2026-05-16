package conversation

type ContentType string

const (
	ContentTypeText ContentType = "text"
)

type ContentPart struct {
	Type ContentType
	Text string
}

func NewTextContent(text string) []ContentPart {
	return []ContentPart{
		{Type: ContentTypeText, Text: text},
	}
}
