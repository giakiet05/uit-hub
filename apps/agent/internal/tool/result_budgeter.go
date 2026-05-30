package tool

import "fmt"

// DefaultMaxResultChars is the default observation size limit for tools.
const DefaultMaxResultChars = 12000

// ResultBudgeter bounds oversized observations before they are appended back to
// the conversation.
type ResultBudgeter struct{}

// NewResultBudgeter creates a result budgeter.
func NewResultBudgeter() ResultBudgeter {
	return ResultBudgeter{}
}

// Apply truncates a result when the tool metadata defines a positive character
// limit.
func (b ResultBudgeter) Apply(result Result, metadata Metadata) Result {
	limit := metadata.MaxResultChars
	contentRunes := []rune(result.Content)
	if limit <= 0 || len(contentRunes) <= limit {
		return result
	}

	result.Content = fmt.Sprintf(
		"[tool result truncated: original_chars=%d max_chars=%d]\n%s",
		len(contentRunes),
		limit,
		string(contentRunes[:limit]),
	)
	return result
}
