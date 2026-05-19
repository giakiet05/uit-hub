package prompt

import _ "embed"

//go:embed templates/react_system.md
var reActStableSystemPrompt string

// ReActSystemPrompt returns the agent-specific stable prompt for ReAct.
func ReActSystemPrompt() SystemPrompt {
	return SystemPrompt{
		StableText: reActStableSystemPrompt,
	}
}
