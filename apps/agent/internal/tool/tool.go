// Package tool defines provider-neutral tool contracts and execution metadata.
package tool

import "context"

// Tool is an executable capability that can be exposed to an LLM.
type Tool interface {
	Definition() Definition
	Execute(ctx context.Context, call Call) (Result, error)
}

// Definition describes a tool to the LLM provider.
type Definition struct {
	Name        string
	Description string
	InputSchema JSONSchema
}

// Call is one model-requested tool invocation.
type Call struct {
	ID        string
	Name      string
	Arguments map[string]any
}

// Result is the observation produced by a tool execution.
type Result struct {
	CallID  string
	Name    string
	Content string
}
