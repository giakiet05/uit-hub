package app

import (
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/mcpadapter"
	"github.com/giakiet05/uit-hub/apps/agent/internal/memory"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
)

// newSessionPrompt builds one session-stable prompt snapshot for the agent.
func newSessionPrompt(agentType string, memoryContext memory.Context, mcpCatalog []mcpadapter.ToolMetadata) prompt.SystemPrompt {
	var agentStatic prompt.StaticPart
	switch agent.Type(agentType) {
	case agent.TypeOrchestrator:
		agentStatic = prompt.OrchestratorAgentStaticPrompt()
	default:
		agentStatic = prompt.ReActStaticPrompt()
	}

	return prompt.SystemPrompt{
		StaticParts: []prompt.StaticPart{
			prompt.IdentityStaticPrompt(),
			prompt.BehaviorStaticPrompt(),
			prompt.ToolUsageStaticPrompt(),
			agentStatic,
		},
		DynamicParts: []prompt.DynamicPart{
			prompt.MemoryDynamicPrompt(memoryContext),
			prompt.MCPToolCatalogDynamicPrompt(mcpCatalog),
		},
	}
}
