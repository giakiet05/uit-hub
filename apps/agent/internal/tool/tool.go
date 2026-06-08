// Package tool defines provider-neutral tool contracts and execution metadata.
package tool

import "context"

// Tool is an executable capability that can be exposed to an LLM.
type Tool interface {
	Definition() Definition
	Metadata() Metadata
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

// Metadata describes execution behavior that the agent runtime needs but the
// LLM provider does not.
type Metadata struct {
	ReadOnly          bool
	Destructive       bool
	ConcurrencySafe   bool
	RequireApproval   bool
	FinishOnInterrupt bool
}

// NewReadOnlyMetadata creates metadata for tools that do not mutate external
// state.
func NewReadOnlyMetadata(concurrencySafe bool) Metadata {
	return Metadata{
		ReadOnly:        true,
		ConcurrencySafe: concurrencySafe,
	}
}

// NewWriteMetadata creates metadata for tools that can mutate state.
func NewWriteMetadata(destructive bool, requireApproval bool) Metadata {
	return Metadata{
		ReadOnly:          false,
		Destructive:       destructive,
		RequireApproval:   requireApproval,
		FinishOnInterrupt: requireApproval,
	}
}

// BaseTool stores stable definition and execution metadata for concrete tools.
type BaseTool struct {
	definition Definition
	metadata   Metadata
}

// NewBaseTool creates reusable tool identity and metadata.
func NewBaseTool(definition Definition, metadata Metadata) BaseTool {
	return BaseTool{
		definition: definition,
		metadata:   metadata,
	}
}

// Definition returns the stable tool definition.
func (t BaseTool) Definition() Definition {
	return t.definition
}

// Metadata returns the stable tool execution metadata.
func (t BaseTool) Metadata() Metadata {
	return t.metadata
}
