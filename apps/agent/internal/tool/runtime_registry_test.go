package tool

import (
	"context"
	"errors"
	"testing"
)

func TestRuntimeRegistryEvictsLeastRecentlyUsedTool(t *testing.T) {
	registry := NewRuntimeRegistry(2)

	if err := registry.Register(testTool{name: "first"}); err != nil {
		t.Fatalf("Register(first) error = %v", err)
	}
	if err := registry.Register(testTool{name: "second"}); err != nil {
		t.Fatalf("Register(second) error = %v", err)
	}
	registry.MarkAccessed("first")
	if err := registry.Register(testTool{name: "third"}); err != nil {
		t.Fatalf("Register(third) error = %v", err)
	}

	if registry.Has("second") {
		t.Fatal("second tool still registered, want evicted")
	}
	for _, name := range []string{"first", "third"} {
		if !registry.Has(name) {
			t.Fatalf("%s tool missing after eviction", name)
		}
	}
}

func TestRuntimeRegistryExecuteMarksToolUsed(t *testing.T) {
	registry := NewRuntimeRegistry(1)
	if err := registry.Register(testTool{name: "first", content: "ok"}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	result, err := registry.Execute(context.Background(), Call{
		ID:   "call-1",
		Name: "first",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got, want := result.Content, "ok"; got != want {
		t.Fatalf("result content = %q, want %q", got, want)
	}
}

func TestToolSetExecutesBaseThenRuntimeTools(t *testing.T) {
	base, err := NewBaseRegistry(testTool{name: "base", content: "base-result"})
	if err != nil {
		t.Fatalf("NewBaseRegistry() error = %v", err)
	}
	runtime := NewRuntimeRegistry(2)
	if err := runtime.Register(testTool{name: "runtime", content: "runtime-result"}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	toolSet := NewToolSet(base, runtime)

	tests := []struct {
		name        string
		call        Call
		wantContent string
		wantErr     error
	}{
		{
			name:        "base tool",
			call:        Call{Name: "base"},
			wantContent: "base-result",
		},
		{
			name:        "runtime tool",
			call:        Call{Name: "runtime"},
			wantContent: "runtime-result",
		},
		{
			name:    "missing tool",
			call:    Call{Name: "missing"},
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toolSet.Execute(context.Background(), tt.call)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Execute() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			if got := result.Content; got != tt.wantContent {
				t.Fatalf("result content = %q, want %q", got, tt.wantContent)
			}
		})
	}
}

func TestToolSetDefinitionsCombinesBaseAndRuntime(t *testing.T) {
	base, err := NewBaseRegistry(testTool{name: "base"})
	if err != nil {
		t.Fatalf("NewBaseRegistry() error = %v", err)
	}
	runtime := NewRuntimeRegistry(2)
	if err := runtime.Register(testTool{name: "runtime"}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	definitions := NewToolSet(base, runtime).Definitions()
	if got, want := len(definitions), 2; got != want {
		t.Fatalf("definitions len = %d, want %d", got, want)
	}
	if got, want := definitions[0].Name, "base"; got != want {
		t.Fatalf("definitions[0].Name = %q, want %q", got, want)
	}
	if got, want := definitions[1].Name, "runtime"; got != want {
		t.Fatalf("definitions[1].Name = %q, want %q", got, want)
	}
}

type testTool struct {
	name    string
	content string
}

func (t testTool) Definition() Definition {
	return Definition{
		Name:        t.name,
		Description: "test tool",
		InputSchema: EmptyInputSchema(),
	}
}

func (t testTool) Metadata() Metadata {
	return NewReadOnlyMetadata(true, DefaultMaxResultChars)
}

func (t testTool) Execute(ctx context.Context, call Call) (Result, error) {
	return Result{
		CallID:  call.ID,
		Name:    call.Name,
		Content: t.content,
	}, nil
}
