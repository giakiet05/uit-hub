package agent

// Type identifies which agent implementation is running.
type Type string

const (
	// TypeReAct identifies the baseline ReAct tool-calling agent.
	TypeReAct Type = "react"
	// TypePlanAndExecute identifies the single-agent plan-and-execute agent.
	TypePlanAndExecute Type = "plan_and_execute"
)

// String returns the stable prompt/type identifier for this agent type.
func (t Type) String() string {
	return string(t)
}
