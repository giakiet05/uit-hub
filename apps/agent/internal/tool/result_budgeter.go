package tool

import "fmt"

// DefaultMaxResultChars is the default observation size limit for tools.
const DefaultMaxResultChars = 12000

// ResultBudgeter bounds oversized observations before they are appended back to
// the conversation.
type ResultBudgeter struct {
	MaxChars int
}

// NewResultBudgeter creates a result budgeter.
func NewResultBudgeter(maxChars int) *ResultBudgeter {
	if maxChars <= 0 {
		maxChars = DefaultMaxResultChars
	}
	return &ResultBudgeter{
		MaxChars: maxChars,
	}
}

// Apply truncates a result when the budget defines a positive character limit.
func (b *ResultBudgeter) Apply(result Result) Result {
	if b == nil {
		return result
	}

	limit := b.MaxChars
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
