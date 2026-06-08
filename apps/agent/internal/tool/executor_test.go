package tool

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestExecutorExecutesAndMarksRuntimeToolUsed(t *testing.T) {
	base, err := NewBaseRegistry()
	if err != nil {
		t.Fatalf("NewBaseRegistry() error = %v", err)
	}
	runtime := NewRuntimeRegistry(2)
	if err := runtime.Register(testTool{name: "first", content: "first-result"}); err != nil {
		t.Fatalf("Register(first) error = %v", err)
	}
	if err := runtime.Register(testTool{name: "second", content: "second-result"}); err != nil {
		t.Fatalf("Register(second) error = %v", err)
	}
	executor := NewExecutor(NewToolSet(base, runtime), 0, nil)

	result, err := executor.Execute(context.Background(), Call{
		ID:   "call-1",
		Name: "first",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got, want := result.Content, "first-result"; got != want {
		t.Fatalf("result content = %q, want %q", got, want)
	}

	if err := runtime.Register(testTool{name: "third", content: "third-result"}); err != nil {
		t.Fatalf("Register(third) error = %v", err)
	}
	if runtime.Has("second") {
		t.Fatal("second tool still registered, want evicted after first was marked used")
	}
	if !runtime.Has("first") {
		t.Fatal("first tool missing, want retained after executor marked it used")
	}
}

func TestExecutorReturnsSanitizedToolError(t *testing.T) {
	base, err := NewBaseRegistry(erroringTestTool{name: "failing"})
	if err != nil {
		t.Fatalf("NewBaseRegistry() error = %v", err)
	}
	executor := NewExecutor(NewToolSet(base, NewRuntimeRegistry(1)), 0, nil)

	_, err = executor.Execute(context.Background(), Call{Name: "failing"})
	if err == nil {
		t.Fatal("Execute() expected error")
	}
	var toolErr ToolError
	if !errors.As(err, &toolErr) {
		t.Fatalf("Execute() error = %T, want ToolError", err)
	}
	if got, want := toolErr.Type, ErrorTypeExecutionFailed; got != want {
		t.Fatalf("tool error type = %q, want %q", got, want)
	}
	if got, want := ErrorObservation(err), "tool error: execution failed"; got != want {
		t.Fatalf("observation = %q, want %q", got, want)
	}
}

func TestExecutorTimesOutTool(t *testing.T) {
	base, err := NewBaseRegistry(blockingTestTool{name: "blocking"})
	if err != nil {
		t.Fatalf("NewBaseRegistry() error = %v", err)
	}
	executor := NewExecutor(NewToolSet(base, NewRuntimeRegistry(1)), time.Millisecond, nil)

	_, err = executor.Execute(context.Background(), Call{Name: "blocking"})
	if err == nil {
		t.Fatal("Execute() expected timeout error")
	}
	if got, want := ErrorObservation(err), "tool error: execution timed out"; got != want {
		t.Fatalf("observation = %q, want %q", got, want)
	}
}

func TestExecutorBudgetsLongResult(t *testing.T) {
	base, err := NewBaseRegistry(budgetedTestTool{
		BaseTool: NewBaseTool(
			Definition{Name: "long", Description: "long result", InputSchema: EmptyInputSchema()},
			NewReadOnlyMetadata(true),
		),
		content: "0123456789",
	})
	if err != nil {
		t.Fatalf("NewBaseRegistry() error = %v", err)
	}
	executor := NewExecutor(NewToolSet(base, NewRuntimeRegistry(1)), 0, NewResultBudgeter(5))

	result, err := executor.Execute(context.Background(), Call{Name: "long"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(result.Content, "tool result truncated") {
		t.Fatalf("result content = %q, want truncation header", result.Content)
	}
	if !strings.HasSuffix(result.Content, "01234") {
		t.Fatalf("result content = %q, want truncated body", result.Content)
	}
}

func TestExecutorDropsReadToolResultAfterRunCancel(t *testing.T) {
	base, err := NewBaseRegistry(testTool{name: "read", content: "done"})
	if err != nil {
		t.Fatalf("NewBaseRegistry() error = %v", err)
	}
	executor := NewExecutor(NewToolSet(base, NewRuntimeRegistry(1)), 0, nil)
	runCtx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = executor.Execute(runCtx, Call{Name: "read"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Execute() error = %v, want context.Canceled", err)
	}
}

func TestExecutorFinishesApprovedToolAfterRunCancel(t *testing.T) {
	base, err := NewBaseRegistry(forceFinishTestTool{name: "send_email"})
	if err != nil {
		t.Fatalf("NewBaseRegistry() error = %v", err)
	}
	executor := NewExecutor(
		NewToolSet(base, NewRuntimeRegistry(1)),
		0,
		nil,
		WithShutdownContext(context.Background()),
	)
	runCtx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := executor.Execute(runCtx, Call{Name: "send_email"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got, want := result.Content, "sent"; got != want {
		t.Fatalf("result content = %q, want %q", got, want)
	}
}

type erroringTestTool struct {
	name string
}

func (t erroringTestTool) Definition() Definition {
	return Definition{Name: t.name, Description: "erroring tool", InputSchema: EmptyInputSchema()}
}

func (t erroringTestTool) Metadata() Metadata {
	return NewReadOnlyMetadata(true)
}

func (t erroringTestTool) Execute(ctx context.Context, call Call) (Result, error) {
	return Result{}, errors.New("secret internal error")
}

type blockingTestTool struct {
	name string
}

func (t blockingTestTool) Definition() Definition {
	return Definition{Name: t.name, Description: "blocking tool", InputSchema: EmptyInputSchema()}
}

func (t blockingTestTool) Metadata() Metadata {
	return NewReadOnlyMetadata(true)
}

func (t blockingTestTool) Execute(ctx context.Context, call Call) (Result, error) {
	<-ctx.Done()
	return Result{}, ctx.Err()
}

type budgetedTestTool struct {
	BaseTool
	content string
}

func (t budgetedTestTool) Execute(ctx context.Context, call Call) (Result, error) {
	return Result{
		CallID:  call.ID,
		Name:    call.Name,
		Content: t.content,
	}, nil
}

type forceFinishTestTool struct {
	name string
}

func (t forceFinishTestTool) Definition() Definition {
	return Definition{Name: t.name, Description: "write tool", InputSchema: EmptyInputSchema()}
}

func (t forceFinishTestTool) Metadata() Metadata {
	return NewWriteMetadata(false, true)
}

func (t forceFinishTestTool) Execute(ctx context.Context, call Call) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	return Result{CallID: call.ID, Name: call.Name, Content: "sent"}, nil
}
