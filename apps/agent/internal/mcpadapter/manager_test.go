package mcpadapter

import (
	"testing"

	agenttool "github.com/giakiet05/uit-hub/apps/agent/internal/tool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestManagerAddLoadedToolsBuildsCatalogAndLoadsTool(t *testing.T) {
	client := NewClient("uit", &fakeSession{})
	loadedTool := NewTool("uit", &fakeSession{}, &mcp.Tool{
		Name:        "get_student",
		Description: "Get student profile.",
	})
	manager := NewManager()

	if err := manager.AddLoadedTools("uit", "UIT academic tools.", client, []agenttool.Tool{loadedTool}); err != nil {
		t.Fatalf("AddLoadedTools() error = %v", err)
	}

	catalog := manager.Catalog()
	if got, want := len(catalog), 1; got != want {
		t.Fatalf("catalog len = %d, want %d", got, want)
	}
	if got, want := catalog[0].Name, "uit__get_student"; got != want {
		t.Fatalf("catalog name = %q, want %q", got, want)
	}
	if got, want := catalog[0].ServerName, "uit"; got != want {
		t.Fatalf("catalog server = %q, want %q", got, want)
	}
	if got, want := catalog[0].ServerDescription, "UIT academic tools."; got != want {
		t.Fatalf("catalog server description = %q, want %q", got, want)
	}
	if got, want := catalog[0].Description, "Get student profile."; got != want {
		t.Fatalf("catalog description = %q, want %q", got, want)
	}

	loaded, ok := manager.LoadTool("uit__get_student")
	if !ok {
		t.Fatal("LoadTool() ok = false, want true")
	}
	if loaded == nil {
		t.Fatal("LoadTool() returned nil tool")
	}
}

func TestManagerLoadToolUnknown(t *testing.T) {
	manager := NewManager()

	loaded, ok := manager.LoadTool("missing__tool")
	if ok {
		t.Fatal("LoadTool() ok = true, want false")
	}
	if loaded != nil {
		t.Fatalf("LoadTool() = %T, want nil", loaded)
	}
}

func TestManagerCloseClosesClients(t *testing.T) {
	session := &fakeSession{}
	client := NewClient("uit", session)
	manager := NewManager()
	if err := manager.AddLoadedTools("uit", "", client, nil); err != nil {
		t.Fatalf("AddLoadedTools() error = %v", err)
	}

	if err := manager.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if !session.closed {
		t.Fatal("session closed = false, want true")
	}
}

func TestManagerRejectsDuplicateTools(t *testing.T) {
	manager := NewManager()
	client := NewClient("uit", &fakeSession{})
	tools := []agenttool.Tool{
		NewTool("uit", &fakeSession{}, &mcp.Tool{Name: "get_student"}),
		NewTool("uit", &fakeSession{}, &mcp.Tool{Name: "get_student"}),
	}

	if err := manager.AddLoadedTools("uit", "", client, tools); err == nil {
		t.Fatal("AddLoadedTools() expected duplicate error")
	}
}

func TestManagerAddLoadedToolsRejectsNilClient(t *testing.T) {
	manager := NewManager()

	if err := manager.AddLoadedTools("uit", "", nil, nil); err == nil {
		t.Fatal("AddLoadedTools() expected nil client error")
	}
}
