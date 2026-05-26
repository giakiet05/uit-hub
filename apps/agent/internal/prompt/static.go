package prompt

import _ "embed"

//go:embed templates/static/identity.md
var identityStaticPrompt string

//go:embed templates/static/behavior.md
var behaviorStaticPrompt string

//go:embed templates/static/tool_usage.md
var toolUsageStaticPrompt string

// IdentityStaticPrompt returns the stable identity and role prompt.
func IdentityStaticPrompt() StaticPart {
	return StaticPart(identityStaticPrompt)
}

// BehaviorStaticPrompt returns stable system behavior rules.
func BehaviorStaticPrompt() StaticPart {
	return StaticPart(behaviorStaticPrompt)
}

// ToolUsageStaticPrompt returns stable detailed tool usage rules.
func ToolUsageStaticPrompt() StaticPart {
	return StaticPart(toolUsageStaticPrompt)
}
