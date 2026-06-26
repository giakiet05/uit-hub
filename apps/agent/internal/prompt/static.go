package prompt

import _ "embed"

//go:embed templates/static/identity.md
var identityStaticPrompt string

//go:embed templates/static/behavior.md
var behaviorStaticPrompt string

//go:embed templates/static/tool_usage.md
var toolUsageStaticPrompt string

//go:embed templates/static/academic_agent.md
var academicAgentStaticPrompt string

//go:embed templates/static/procedure_agent.md
var procedureAgentStaticPrompt string

//go:embed templates/static/campus_agent.md
var campusAgentStaticPrompt string

//go:embed templates/static/orchestrator_agent.md
var orchestratorAgentStaticPrompt string

//go:embed templates/static/subagent_rules.md
var subagentRulesStaticPrompt string

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

// AcademicAgentStaticPrompt returns stable identity for academic agent.
func AcademicAgentStaticPrompt() StaticPart {
	return StaticPart(academicAgentStaticPrompt)
}

// ProcedureAgentStaticPrompt returns stable identity for procedure agent.
func ProcedureAgentStaticPrompt() StaticPart {
	return StaticPart(procedureAgentStaticPrompt)
}

// CampusAgentStaticPrompt returns stable identity for campus agent.
func CampusAgentStaticPrompt() StaticPart {
	return StaticPart(campusAgentStaticPrompt)
}

// OrchestratorAgentStaticPrompt returns stable identity for orchestrator agent.
func OrchestratorAgentStaticPrompt() StaticPart {
	return StaticPart(orchestratorAgentStaticPrompt)
}

// SubAgentRulesStaticPrompt returns strict formatting rules for sub-agents.
func SubAgentRulesStaticPrompt() StaticPart {
	return StaticPart(subagentRulesStaticPrompt)
}
