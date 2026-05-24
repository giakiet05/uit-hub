package tool

import (
	"encoding/json"
	"testing"
)

func TestEmptyInputSchemaRequiredIsArray(t *testing.T) {
	data, err := json.Marshal(EmptyInputSchema())
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if string(data) == `{"additionalProperties":false,"properties":{},"required":null,"type":"object"}` {
		t.Fatalf("schema has null required: %s", data)
	}

	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if _, ok := decoded["required"].([]any); !ok {
		t.Fatalf("required = %T, want JSON array", decoded["required"])
	}
}
