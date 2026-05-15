package llm

import "context"

type Provider interface {
	Generate(ctx context.Context, request GenerateRequest) (GenerateResponse, error)
}
