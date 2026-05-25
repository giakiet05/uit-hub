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
	"github.com/giakiet05/uit-hub/apps/agent/internal/memory"
)

func TestNewToolSetLoadsMockMCPToolCatalog(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping MCP stdio integration smoke test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	repoRoot := filepath.Clean("../../../..")
	toolSet, mcpManager, closers, err := newToolSet(ctx, config.Config{
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
	}, logging.NewNopLogger(), nil)
	if err != nil {
		t.Fatalf("newToolSet() error = %v", err)
	}
	defer closeAll(ctx, logging.NewNopLogger(), closers)

	names := []string{}
	for _, definition := range toolSet.Definitions() {
		names = append(names, definition.Name)
	}
	if !slices.Contains(names, "load_mcp_tool") {
		t.Fatalf("base tools missing load_mcp_tool: %v", names)
	}
	if slices.Contains(names, "mock_uit__get_student_profile") {
		t.Fatalf("MCP tool was eagerly registered in base tools: %v", names)
	}

	catalogNames := []string{}
	for _, entry := range mcpManager.Catalog() {
		catalogNames = append(catalogNames, entry.Name)
	}
	for _, want := range []string{
		"mock_uit__get_student_profile",
		"mock_uit__search_courses",
		"mock_uit__upsert_student_note",
		"mock_uit__list_student_notes",
		"mock_uit__delete_student_note",
		"mock_uit__submit_leave_request",
	} {
		if !slices.Contains(catalogNames, want) {
			t.Fatalf("catalog missing %q: %v", want, catalogNames)
		}
	}
}

func TestNewToolSetKeepsWorkingWhenMCPServerFails(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping MCP stdio integration smoke test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	repoRoot := filepath.Clean("../../../..")
	attemptFile := filepath.Join(t.TempDir(), "retry-attempt")
	toolSet, mcpManager, closers, err := newToolSet(ctx, config.Config{
		Provider: llm.ProviderTypeOpenAI,
		MCP: config.MCPConfig{
			ClientName:    "uit-hub-agent-test",
			ClientVersion: "0.1.0",
			Servers: []config.MCPServerConfig{
				{
					Name:      "broken",
					Transport: "stdio",
					Command:   "go",
					Args:      []string{"run", "./apps/mock-mcp-server/cmd/failing-server"},
					Env:       []string{"UIT_HUB_RETRY_MCP_ATTEMPT_FILE=" + attemptFile},
					WorkDir:   repoRoot,
				},
				{
					Name:      "mock_uit",
					Transport: "stdio",
					Command:   "go",
					Args:      []string{"run", "./apps/mock-mcp-server/cmd/server"},
					WorkDir:   repoRoot,
				},
			},
		},
	}, logging.NewNopLogger(), nil)
	if err != nil {
		t.Fatalf("newToolSet() error = %v", err)
	}
	defer closeAll(ctx, logging.NewNopLogger(), closers)

	names := []string{}
	for _, definition := range toolSet.Definitions() {
		names = append(names, definition.Name)
	}
	if !slices.Contains(names, "load_mcp_tool") {
		t.Fatalf("base tools missing load_mcp_tool: %v", names)
	}

	catalogNames := []string{}
	for _, entry := range mcpManager.Catalog() {
		catalogNames = append(catalogNames, entry.Name)
	}
	if !slices.Contains(catalogNames, "broken__retry_ping") {
		t.Fatalf("retry server catalog missing retry_ping after third attempt: %v", catalogNames)
	}
	if !slices.Contains(catalogNames, "mock_uit__get_student_profile") {
		t.Fatalf("successful server catalog missing get_student_profile: %v", catalogNames)
	}
}

func TestNewToolSetRegistersMemoryTools(t *testing.T) {
	ctx := context.Background()

	toolSet, _, closers, err := newToolSet(ctx, config.Config{
		Provider: llm.ProviderTypeOpenAI,
		Memory: config.MemoryConfig{
			Path: t.TempDir(),
		},
	}, logging.NewNopLogger(), memory.NewFileStore(t.TempDir()))
	if err != nil {
		t.Fatalf("newToolSet() error = %v", err)
	}
	defer closeAll(ctx, logging.NewNopLogger(), closers)

	names := []string{}
	for _, definition := range toolSet.Definitions() {
		names = append(names, definition.Name)
	}

	for _, want := range []string{"memory_list", "memory_read", "memory_write"} {
		if !slices.Contains(names, want) {
			t.Fatalf("registered tools missing %q: %v", want, names)
		}
	}
}
