// Package llm defines provider-neutral request and response types for language
// model calls.
package llm

import "context"

// ProviderType identifies a concrete LLM provider implementation.
type ProviderType string

const (
	// ProviderTypeOpenAI selects the OpenAI provider.
	ProviderTypeOpenAI ProviderType = "openai"
)

// Provider generates assistant messages from conversation history and available
// tool definitions.
type Provider interface {
	Generate(ctx context.Context, request GenerateRequest) (GenerateResponse, error)
}
