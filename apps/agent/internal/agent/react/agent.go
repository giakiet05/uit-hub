package react

import (
	"context"
	"log/slog"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
	"github.com/giakiet05/uit-hub/apps/agent/internal/prompt"
	"github.com/giakiet05/uit-hub/apps/agent/internal/runtime"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

const defaultMaxToolIterations = 8
const defaultToolTimeout = 15 * time.Second

// ReActAgent implements the baseline reason-act-observe loop for tool-calling
// agents.
type ReActAgent struct {
	provider      llm.Provider
	promptBuilder *prompt.Builder
	tools         *tool.ToolSet
	logger        *slog.Logger
	maxRounds     int
	toolTimeout   time.Duration
}

// ReActOption configures optional ReActAgent dependencies and limits.
type ReActOption func(*ReActAgent)

// NewReActAgent constructs a ReActAgent with defaults, then applies overrides.
func NewReActAgent(provider llm.Provider, opts ...ReActOption) *ReActAgent {
	agent := &ReActAgent{
		provider:      provider,
		promptBuilder: prompt.NewBuilder(prompt.SessionPrompt{}),
		logger:        logging.NewNopLogger(),
		maxRounds:     defaultMaxToolIterations,
		toolTimeout:   defaultToolTimeout,
	}

	for _, opt := range opts {
		opt(agent)
	}

	return agent
}

// WithPromptBuilder overrides the prompt builder.
func WithPromptBuilder(promptBuilder *prompt.Builder) ReActOption {
	return func(agent *ReActAgent) {
		if promptBuilder != nil {
			agent.promptBuilder = promptBuilder
		}
	}
}

// WithTools sets the tool collection used by the agent.
func WithTools(tools *tool.ToolSet) ReActOption {
	return func(agent *ReActAgent) {
		agent.tools = tools
	}
}

// WithLogger sets the structured logger.
func WithLogger(logger *slog.Logger) ReActOption {
	return func(agent *ReActAgent) {
		if logger != nil {
			agent.logger = logger
		}
	}
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

// Run appends the user prompt to the session and streams loop events until the
// model returns a final answer or the configured round limit is reached.
func (a *ReActAgent) Run(ctx context.Context, session *runtime.Session, userPrompt string) <-chan agent.Event {
	events := make(chan agent.Event)
	go a.run(ctx, events, session, userPrompt)
	return events
}
