package mcpadapter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// FlattenToolResult converts MCP content blocks into text for the LLM.
func FlattenToolResult(result *mcp.CallToolResult) (string, error) {
	parts := make([]string, 0, len(result.Content)+1)
	for _, content := range result.Content {
		switch typed := content.(type) {
		case *mcp.TextContent:
			if typed.Text != "" {
				parts = append(parts, typed.Text)
			}
		case *mcp.EmbeddedResource:
			if typed.Resource != nil {
				parts = append(parts, flattenResource(*typed.Resource))
			}
		case *mcp.ResourceLink:
			parts = append(parts, fmt.Sprintf("resource: %s", typed.URI))
		default:
			data, err := content.MarshalJSON()
			if err != nil {
				return "", fmt.Errorf("marshal unsupported mcp content: %w", err)
			}
			parts = append(parts, string(data))
		}
	}

	if result.StructuredContent != nil {
		data, err := json.Marshal(result.StructuredContent)
		if err != nil {
			return "", fmt.Errorf("marshal structured mcp content: %w", err)
		}
		parts = append(parts, string(data))
	}

	return strings.Join(parts, "\n"), nil
}

func flattenResource(resource mcp.ResourceContents) string {
	if resource.Text != "" {
		return resource.Text
	}
	if len(resource.Blob) > 0 {
		return fmt.Sprintf("resource blob: %s bytes=%d", resource.URI, len(resource.Blob))
	}
	return fmt.Sprintf("resource: %s", resource.URI)
}
