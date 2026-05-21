package mcpadapter

import (
	"encoding/json"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// InputSchemaFromMCP converts an MCP input schema value to the native tool
// schema shape.
func InputSchemaFromMCP(schema any) tool.JSONSchema {
	if schema == nil {
		return emptyObjectSchema()
	}

	data, err := json.Marshal(schema)
	if err != nil {
		return emptyObjectSchema()
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		return emptyObjectSchema()
	}
	return tool.JSONSchema(out)
}

func emptyObjectSchema() tool.JSONSchema {
	return tool.EmptyInputSchema()
}
