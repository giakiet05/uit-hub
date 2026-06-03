package react

import (
	"context"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
)

const defaultMaxToolIterations = 8
const defaultToolTimeout = 15 * time.Second

// ReActAgent implements the baseline reason-act-observe loop for tool-calling
// agents.
type ReActAgent struct {
	model       llm.Model
	maxRounds   int
	toolTimeout time.Duration
	compaction  conversation.CompactionOptions
}

// ReActOption configures optional ReActAgent dependencies and limits.
type ReActOption func(*ReActAgent)

// NewReActAgent constructs a ReActAgent with defaults, then applies overrides.
func NewReActAgent(model llm.Model, opts ...ReActOption) *ReActAgent {
	agent := &ReActAgent{
		model:       model,
		maxRounds:   defaultMaxToolIterations,
		toolTimeout: defaultToolTimeout,
		compaction:  conversation.CompactionOptions{},
	}

	for _, opt := range opts {
		opt(agent)
	}

	return agent
}

// WithMaxRounds overrides the maximum ReAct loop rounds when maxRounds is
// positive.
func WithMaxRounds(maxRounds int) ReActOption {
	return func(agent *ReActAgent) {
		if maxRounds > 0 {
			agent.maxRounds = maxRounds
		}
	}
}

// WithToolTimeout overrides per-tool execution timeout when timeout is
// positive.
func WithToolTimeout(timeout time.Duration) ReActOption {
	return func(agent *ReActAgent) {
		if timeout > 0 {
			agent.toolTimeout = timeout
		}
	}
}

// WithCompaction configures lightweight conversation compaction.
func WithCompaction(options conversation.CompactionOptions) ReActOption {
	return func(agent *ReActAgent) {
		agent.compaction = options
	}
}

// Run appends the user prompt to the session and streams loop events until the
// model returns a final answer or the configured round limit is reached.
func (a *ReActAgent) Run(ctx context.Context, input agent.RunInput) <-chan agent.Event {
	events := make(chan agent.Event)
	go a.run(ctx, events, input)
	return events
}
