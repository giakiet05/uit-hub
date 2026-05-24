package app

import (
	"github.com/giakiet05/uit-hub/apps/agent/internal/memory"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
)

// newSessionPrompt builds one session-stable prompt snapshot for the agent.
func newSessionPrompt(memoryContext memory.Context) prompt.SessionPrompt {
	return prompt.SessionPrompt{
		StaticParts: []prompt.StaticPart{
			prompt.ReActStaticPrompt(),
		},
		DynamicParts: []prompt.DynamicPart{
			prompt.MemoryDynamicPrompt(memoryContext),
		},
	}
}
