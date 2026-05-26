package prompt

import _ "embed"

//go:embed templates/react_system.md
var reActStableSystemPrompt string

// ReActStaticPrompt returns the agent-specific static prompt for ReAct.
func ReActStaticPrompt() StaticPart {
	return StaticPart(reActStableSystemPrompt)
}
