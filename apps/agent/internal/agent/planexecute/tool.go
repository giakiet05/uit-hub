package planexecute

import (
	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

func (a *PlanAndExecuteAgent) stepToolDefinitions(input agent.RunInput, step planStep) []tool.Definition {
	if input.Tools == nil {
		return nil
	}

	names := []string{
		"load_mcp_tool",
		"calculator",
		"read_file",
		"write_file",
		"memory_read",
		"memory_write",
		"memory_list",
	}
	names = append(names, step.SuggestedTools...)

	seen := map[string]struct{}{}
	definitions := []tool.Definition{}
	for _, name := range names {
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}

		definition, exists := input.ToolDefinition(name)
		if !exists {
			continue
		}
		definitions = append(definitions, definition)
	}

	if len(definitions) == 0 {
		return input.ToolDefinitions()
	}
	return definitions
}
