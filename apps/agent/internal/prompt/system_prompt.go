package prompt

import "strings"

const (
	dynamicContextHeader  = "=== DYNAMIC CONTEXT ==="
	uncachedContextHeader = "=== UNCACHED CONTEXT ==="
)

// StaticPart is prompt text that should stay stable for every runtime session.
type StaticPart string

// DynamicPart is session-scoped prompt text loaded at session start.
type DynamicPart string

// UncachedPart is volatile prompt text built for one LLM call.
type UncachedPart string

// SessionPrompt contains prompt parts stable for one runtime session.
type SessionPrompt struct {
	StaticParts  []StaticPart
	DynamicParts []DynamicPart
}

// CallPrompt contains one LLM-call prompt: session-stable prompt plus volatile
// per-call context.
type CallPrompt struct {
	Session       SessionPrompt
	UncachedParts []UncachedPart
}

// Render converts one call prompt into the final system prompt text.
func (p CallPrompt) Render() string {
	parts := make([]string, 0, 3)

	if text := renderStaticParts(p.Session.StaticParts); text != "" {
		parts = append(parts, text)
	}
	if text := renderDynamicParts(p.Session.DynamicParts); text != "" {
		parts = append(parts, dynamicContextHeader+"\n\n"+text)
	}
	if text := renderUncachedParts(p.UncachedParts); text != "" {
		parts = append(parts, uncachedContextHeader+"\n\n"+text)
	}

	return strings.Join(parts, "\n\n")
}

// renderStaticParts joins static prompt parts while skipping empty parts.
func renderStaticParts(parts []StaticPart) string {
	texts := make([]string, 0, len(parts))
	for _, part := range parts {
		if text := strings.TrimSpace(string(part)); text != "" {
			texts = append(texts, text)
		}
	}
	return strings.Join(texts, "\n\n")
}

// renderDynamicParts joins dynamic prompt parts while skipping empty parts.
func renderDynamicParts(parts []DynamicPart) string {
	texts := make([]string, 0, len(parts))
	for _, part := range parts {
		if text := strings.TrimSpace(string(part)); text != "" {
			texts = append(texts, text)
		}
	}
	return strings.Join(texts, "\n\n")
}

// renderUncachedParts joins volatile prompt parts while skipping empty parts.
func renderUncachedParts(parts []UncachedPart) string {
	texts := make([]string, 0, len(parts))
	for _, part := range parts {
		if text := strings.TrimSpace(string(part)); text != "" {
			texts = append(texts, text)
		}
	}
	return strings.Join(texts, "\n\n")
}
