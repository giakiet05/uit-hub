package prompt

// AgentType identifies which agent-specific prompt should be used.
type AgentType string

const (
	// AgentTypeReAct builds prompts for the baseline ReAct tool-calling agent.
	AgentTypeReAct AgentType = "react"
)
