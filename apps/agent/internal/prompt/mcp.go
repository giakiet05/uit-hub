package prompt

import (
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/mcpadapter"
)

type mcpServerCatalogEntry struct {
	Name        string
	Description string
	Tools       []mcpadapter.ToolMetadata
}

// MCPToolCatalogDynamicPrompt returns the session-start MCP tool catalog.
func MCPToolCatalogDynamicPrompt(catalog []mcpadapter.ToolMetadata) DynamicPart {
	if len(catalog) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("## MCP Tool Catalog\n\n")
	builder.WriteString("MCP servers expose external capabilities. Each catalog entry has format: tool_name: short description.\n")
	builder.WriteString("Use load_mcp_tool with the exact tool name before using any catalog tool or inspecting its input schema. ")
	builder.WriteString("After loading a tool, use it in the next round.\n\n")

	for _, server := range groupMCPToolsByServer(catalog) {
		builder.WriteString("### ")
		builder.WriteString(server.Name)
		builder.WriteString("\n")
		if server.Description != "" {
			builder.WriteString(server.Description)
			builder.WriteString("\n")
		}
		builder.WriteString("Tools:\n")
		for _, entry := range server.Tools {
			builder.WriteString("- ")
			builder.WriteString(entry.Name)
			if entry.Description != "" {
				builder.WriteString(": ")
				builder.WriteString(entry.Description)
			}
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}

	return DynamicPart(strings.TrimSpace(builder.String()))
}

func groupMCPToolsByServer(catalog []mcpadapter.ToolMetadata) []mcpServerCatalogEntry {
	servers := make([]mcpServerCatalogEntry, 0)
	indexByName := make(map[string]int)
	for _, entry := range catalog {
		index, exists := indexByName[entry.ServerName]
		if !exists {
			index = len(servers)
			indexByName[entry.ServerName] = index
			servers = append(servers, mcpServerCatalogEntry{
				Name:        entry.ServerName,
				Description: entry.ServerDescription,
				Tools:       []mcpadapter.ToolMetadata{},
			})
		}

		if servers[index].Description == "" && entry.ServerDescription != "" {
			servers[index].Description = entry.ServerDescription
		}
		servers[index].Tools = append(servers[index].Tools, entry)
	}

	return servers
}
