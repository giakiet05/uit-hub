package tool

import "context"

type Tool interface {
	Definition() Definition
	Execute(ctx context.Context, call Call) (Result, error)
}

type Definition struct {
	Name        string
	Description string
	InputSchema map[string]any
}

type Call struct {
	ID        string
	Name      string
	Arguments map[string]any
}

type Result struct {
	CallID  string
	Name    string
	Content string
}
