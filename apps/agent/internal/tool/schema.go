package tool

import "encoding/json"

// JSONSchema is the provider-neutral JSON Schema shape used for tool inputs.
type JSONSchema map[string]any

// EmptyInputSchema returns a schema for tools that accept no arguments.
func EmptyInputSchema() JSONSchema {
	return ObjectSchema(map[string]any{})
}

// ObjectSchema creates a closed object JSON Schema for tool arguments.
func ObjectSchema(properties map[string]any, required ...string) JSONSchema {
	if properties == nil {
		properties = map[string]any{}
	}
	if required == nil {
		required = []string{}
	}
	return JSONSchema{
		"type":                 "object",
		"additionalProperties": false,
		"properties":           properties,
		"required":             required,
	}
}

// StringProperty creates a string property schema.
func StringProperty(description string) map[string]any {
	return map[string]any{
		"type":        "string",
		"description": description,
	}
}

// StringEnumProperty creates a string enum property schema.
func StringEnumProperty(description string, values ...string) map[string]any {
	property := StringProperty(description)
	property["enum"] = values
	return property
}

// NumberProperty creates a number property schema.
func NumberProperty(description string) map[string]any {
	return map[string]any{
		"type":        "number",
		"description": description,
	}
}

// Clone returns a deep copy that can be safely handed to provider SDKs.
func (s JSONSchema) Clone() map[string]any {
	if s == nil {
		return EmptyInputSchema().Clone()
	}

	data, err := json.Marshal(map[string]any(s))
	if err != nil {
		return emptySchemaMap()
	}

	out := map[string]any{}
	if err := json.Unmarshal(data, &out); err != nil {
		return emptySchemaMap()
	}
	return out
}

func emptySchemaMap() map[string]any {
	return map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"properties":           map[string]any{},
		"required":             []string{},
	}
}
