package agent

import (
	"encoding/json"
	"fmt"
	"strings"
)

// PreviewValue renders a structured value as a compact single-line log string.
func PreviewValue(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return PreviewText(fmt.Sprintf("%v", value))
	}
	return PreviewText(string(data))
}

// PreviewText collapses newlines so structured log fields stay on one line.
func PreviewText(text string) string {
	return strings.ReplaceAll(text, "\n", "\\n")
}
