package localtool

import (
	"context"
	"strings"
	"testing"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

func TestKnowledgeLookupFindsPrivateKnowledge(t *testing.T) {
	lookup := NewKnowledgeLookup()

	result, err := lookup.Execute(context.Background(), tool.Call{
		ID:   "call-knowledge",
		Name: "knowledge_lookup",
		Arguments: map[string]any{
			"query": "What is the UIT Hub secret code?",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(result.Content, "UIT-HUB-2026-ALPHA") {
		t.Fatalf("result content = %q, want secret code", result.Content)
	}
}

func TestKnowledgeLookupFindsVietnameseSecretCodeQuery(t *testing.T) {
	lookup := NewKnowledgeLookup()

	result, err := lookup.Execute(context.Background(), tool.Call{
		ID:   "call-knowledge",
		Name: "knowledge_lookup",
		Arguments: map[string]any{
			"query": "mã bí mật nội bộ UIT Hub là gì?",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(result.Content, "UIT-HUB-2026-ALPHA") {
		t.Fatalf("result content = %q, want secret code", result.Content)
	}
}
