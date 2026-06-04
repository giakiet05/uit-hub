package llm

import (
	"context"
	"log/slog"
)

// Model represents a specific language model configurations and capabilities.
type Model interface {
	Name() string
	MaxContextTokens() int
	Generate(ctx context.Context, request GenerateRequest) (GenerateResponse, error)
}

// StreamModel is implemented by models that can stream visible assistant
// text before returning the final response.
type StreamModel interface {
	Model
	Stream(ctx context.Context, request GenerateRequest) <-chan StreamEvent
}

// BaseModel provides common implementations for Model interface methods.
// It is intended to be embedded in concrete model implementations.
type BaseModel struct {
	name             string
	maxContextTokens int
	logger           *slog.Logger
}

// NewBaseModel creates a new BaseModel instance.
func NewBaseModel(name string, maxContextTokens int, logger *slog.Logger) BaseModel {
	return BaseModel{
		name:             name,
		maxContextTokens: maxContextTokens,
		logger:           logger,
	}
}

// Name returns the model's name.
func (m *BaseModel) Name() string {
	return m.name
}

// MaxContextTokens returns the maximum context tokens supported by the model.
func (m *BaseModel) MaxContextTokens() int {
	return m.maxContextTokens
}

// Logger returns the model's logger.
func (m *BaseModel) Logger() *slog.Logger {
	return m.logger
}
