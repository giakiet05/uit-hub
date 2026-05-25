package prompt

import (
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/memory"
)

// MemoryDynamicPrompt returns session-scoped memory instructions and index.
func MemoryDynamicPrompt(memoryContext memory.Context) DynamicPart {
	indexText := strings.TrimSpace(memoryContext.IndexText)
	if indexText == "" {
		return ""
	}

	return DynamicPart(`## Long-Term Memory

The following memory index was loaded at session start. It may be stale during this session if memory tools write new memories later.

` + indexText + `

Use memory_read when an index entry looks relevant.
Use memory_write only for stable user preferences, durable feedback, project context, or references.
Do not save transient task data, student records, grades, or data that should be fetched from source systems.`)
}
