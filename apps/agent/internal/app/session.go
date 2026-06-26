package app

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/llm"
	"github.com/giakiet05/uit-hub/apps/agent/internal/memory"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// newSessionState creates one chat session with stable memory, prompt, and tool
// snapshots.
func newSessionState(ctx context.Context, state *State, model llm.Model, runtimeToolNames []string, sessionOpts ...session.Option) (*session.State, error) {
	memoryContext, err := memory.NewLoader(state.MemoryStore).Load(ctx)
	if err != nil {
		state.Logger.DebugContext(ctx, "Memory context load failed", "error", err)
		return nil, err
	}

	tools, err := newSessionToolSet(state.Config, state.MemoryStore, state.MCPManager, model, runtimeToolNames)
	if err != nil {
		return nil, err
	}

	promptSnapshot := newSessionPrompt(string(state.Config.Agent.Type), memoryContext, state.MCPManager.Catalog())
	opts := append([]session.Option{
		session.WithPromptSnapshot(promptSnapshot),
		session.WithTools(tools),
		session.WithResultBudgeter(newResultBudgeter(state.Config)),
		session.WithConcurrentTools(state.Config.Agent.ConcurrentTools),
	}, sessionOpts...)
	return session.NewState(opts...), nil
}

func newResultBudgeter(cfg config.Config) *tool.ResultBudgeter {
	if !cfg.Compaction.ToolResultBudgetEnabled {
		return nil
	}
	return tool.NewResultBudgeter(cfg.Compaction.ToolResultBudgetMaxChars)
}
