package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReadsMCPConfigFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mcp.yaml")
	content := []byte(`client:
  name: uit-hub-agent-test
  version: 0.2.0
servers:
  - name: mock_uit
    description: Mock UIT academic MCP server.
    transport: stdio
    command: go
    args:
      - run
      - ../mock-mcp-server/cmd/server
    env:
      - UIT_API_BASE_URL=http://localhost:8080
    workdir: ../agent
`)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("write MCP config: %v", err)
	}

	t.Setenv("OPENAI_API_KEY", "test-key")
	t.Setenv("MCP_CONFIG_PATH", path)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.MCP.ClientName != "uit-hub-agent-test" {
		t.Fatalf("ClientName = %q, want %q", cfg.MCP.ClientName, "uit-hub-agent-test")
	}
	if cfg.MCP.ClientVersion != "0.2.0" {
		t.Fatalf("ClientVersion = %q, want %q", cfg.MCP.ClientVersion, "0.2.0")
	}
	if len(cfg.MCP.Servers) != 1 {
		t.Fatalf("len(Servers) = %d, want 1", len(cfg.MCP.Servers))
	}
	server := cfg.MCP.Servers[0]
	if server.Name != "mock_uit" || server.Transport != "stdio" || server.Command != "go" {
		t.Fatalf("server = %+v", server)
	}
	if server.Description != "Mock UIT academic MCP server." {
		t.Fatalf("server.Description = %q, want description", server.Description)
	}
	if len(server.Args) != 2 || server.Args[0] != "run" || server.Args[1] != "../mock-mcp-server/cmd/server" {
		t.Fatalf("server.Args = %#v", server.Args)
	}
	if len(server.Env) != 1 || server.Env[0] != "UIT_API_BASE_URL=http://localhost:8080" {
		t.Fatalf("server.Env = %#v", server.Env)
	}
	if server.WorkDir != "../agent" {
		t.Fatalf("server.WorkDir = %q, want %q", server.WorkDir, "../agent")
	}
}

func TestLoadUsesEmptyMCPConfigWhenFileIsMissing(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "test-key")
	t.Setenv("MCP_CONFIG_PATH", filepath.Join(t.TempDir(), "missing.yaml"))

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.MCP.ClientName != "uit-hub-agent" {
		t.Fatalf("ClientName = %q, want default", cfg.MCP.ClientName)
	}
	if len(cfg.MCP.Servers) != 0 {
		t.Fatalf("len(Servers) = %d, want 0", len(cfg.MCP.Servers))
	}
}

func TestLoadReadsMemoryConfig(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "test-key")
	t.Setenv("MEMORY_PATH", "custom-memory")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Memory.Path != "custom-memory" {
		t.Fatalf("Memory.Path = %q, want custom-memory", cfg.Memory.Path)
	}
}
