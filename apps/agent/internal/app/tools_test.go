package app

import (
	"context"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
)

func TestNewToolRegistryRegistersMockMCPTools(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping MCP stdio integration smoke test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	repoRoot := filepath.Clean("../../../..")
	registry, closers, err := newToolRegistry(ctx, config.Config{
		Provider: llm.ProviderTypeOpenAI,
		MCP: config.MCPConfig{
			ClientName:    "uit-hub-agent-test",
			ClientVersion: "0.1.0",
			Servers: []config.MCPServerConfig{
				{
					Name:      "mock_uit",
					Transport: "stdio",
					Command:   "go",
					Args:      []string{"run", "./apps/mock-mcp-server/cmd/server"},
					WorkDir:   repoRoot,
				},
			},
		},
	}, logging.NewNopLogger())
	if err != nil {
		t.Fatalf("newToolRegistry() error = %v", err)
	}
	defer closeAll(ctx, logging.NewNopLogger(), closers)

	names := []string{}
	for _, definition := range registry.Definitions() {
		names = append(names, definition.Name)
	}

	for _, want := range []string{
		"mock_uit__get_student_profile",
		"mock_uit__search_courses",
		"mock_uit__upsert_student_note",
		"mock_uit__list_student_notes",
		"mock_uit__delete_student_note",
		"mock_uit__submit_leave_request",
	} {
		if !slices.Contains(names, want) {
			t.Fatalf("registered tools missing %q: %v", want, names)
		}
	}
}
