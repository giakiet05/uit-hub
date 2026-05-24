package memory

import (
	"fmt"
	"strings"
)

// FormatIndex renders the lightweight memory table of contents for prompt
// injection.
func FormatIndex(entries []IndexEntry) string {
	if len(entries) == 0 {
		return "No saved memories."
	}

	var builder strings.Builder
	builder.WriteString("# Memory Index\n\n")
	for _, entry := range entries {
		fmt.Fprintf(
			&builder,
			"- [%s](%s) [%s] -- %s\n",
			entry.Name,
			entry.ID,
			entry.Type,
			entry.Description,
		)
	}
	return strings.TrimSpace(builder.String())
}
