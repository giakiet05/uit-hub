package prompt

import "strings"

const (
	sessionContextHeader  = "=== SESSION CONTEXT ==="
	uncachedContextHeader = "=== UNCACHED CONTEXT ==="
)

// SystemPrompt groups prompt text by cache stability.
type SystemPrompt struct {
	StableText  string
	SessionText string
	Uncached    []UncachedSystemPrompt
}

// UncachedSystemPrompt is a prompt block that should not participate in global
// cache scope. Reason documents why this block is volatile.
type UncachedSystemPrompt struct {
	Text   string
	Reason string
}

// Render converts the structured prompt into one system prompt string.
func (p SystemPrompt) Render() string {
	parts := make([]string, 0, 3)

	if text := strings.TrimSpace(p.StableText); text != "" {
		parts = append(parts, text)
	}
	if text := strings.TrimSpace(p.SessionText); text != "" {
		parts = append(parts, sessionContextHeader+"\n\n"+text)
	}
	if text := renderUncached(p.Uncached); text != "" {
		parts = append(parts, uncachedContextHeader+"\n\n"+text)
	}

	return strings.Join(parts, "\n\n")
}

// Merge returns a prompt where non-empty fields from next are appended after
// the current prompt fields in their matching stability tier.
func (p SystemPrompt) Merge(next SystemPrompt) SystemPrompt {
	return SystemPrompt{
		StableText:  joinPromptText(p.StableText, next.StableText),
		SessionText: joinPromptText(p.SessionText, next.SessionText),
		Uncached:    append(append([]UncachedSystemPrompt{}, p.Uncached...), next.Uncached...),
	}
}

func renderUncached(blocks []UncachedSystemPrompt) string {
	parts := make([]string, 0, len(blocks))
	for _, block := range blocks {
		text := strings.TrimSpace(block.Text)
		if text == "" {
			continue
		}
		parts = append(parts, text)
	}
	return strings.Join(parts, "\n\n")
}

func joinPromptText(left string, right string) string {
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)

	switch {
	case left == "":
		return right
	case right == "":
		return left
	default:
		return left + "\n\n" + right
	}
}
