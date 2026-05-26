package memory

import "context"

// Context is the session-scoped memory snapshot injected into the system
// prompt.
type Context struct {
	IndexText string
}

// Loader builds the memory prompt context from a Store once per session.
type Loader struct {
	store Store
}

// NewLoader creates a memory context loader.
func NewLoader(store Store) *Loader {
	return &Loader{store: store}
}

// Load returns the prompt-ready memory index for the current session.
func (l *Loader) Load(ctx context.Context) (Context, error) {
	if l == nil || l.store == nil {
		return Context{}, nil
	}

	entries, err := l.store.List(ctx)
	if err != nil {
		return Context{}, err
	}

	return Context{
		IndexText: FormatIndex(entries),
	}, nil
}
