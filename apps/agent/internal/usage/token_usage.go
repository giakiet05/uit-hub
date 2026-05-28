// Package usage defines shared usage counters collected across agent layers.
package usage

// TokenUsage contains provider-neutral token counters.
type TokenUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// Add merges another token usage value into this one.
func (u *TokenUsage) Add(delta TokenUsage) {
	if u == nil {
		return
	}
	u.InputTokens += delta.InputTokens
	u.OutputTokens += delta.OutputTokens
}
