// Package llm defines provider-neutral request and response types for language
// model calls.
package llm

// ProviderType identifies a concrete LLM provider implementation.
type ProviderType string

const (
	ProviderTypeOpenAI ProviderType = "openai"
)

// Provider represents a service that can instantiate language models.
type Provider interface {
	Name() string
	CreateModel(modelName string, maxContext int) (Model, error)
}

// BaseProvider holds common configuration for LLM providers.
// It is intended to be embedded in concrete provider implementations.
type BaseProvider struct {
	APIKey  string
	BaseURL string
}

// NewBaseProvider creates a new BaseProvider instance.
func NewBaseProvider(apiKey, baseURL string) BaseProvider {
	return BaseProvider{
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
}
