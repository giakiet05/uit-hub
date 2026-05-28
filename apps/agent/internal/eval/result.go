package eval

import (
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
	"github.com/giakiet05/uit-hub/apps/agent/internal/usage"
)

// Result is the machine-readable summary written after one eval run.
type Result struct {
	usage.TokenUsage
	CaseID         string               `json:"case_id"`
	Title          string               `json:"title"`
	StartedAt      time.Time            `json:"started_at"`
	FinishedAt     time.Time            `json:"finished_at"`
	Duration       string               `json:"duration"`
	Passed         bool                 `json:"passed"`
	TerminalReason agent.TerminalReason `json:"terminal_reason"`
	Error          string               `json:"error,omitempty"`
	RoundCount     int                  `json:"round_count"`
	LLMCalls       int                  `json:"llm_calls"`
	ToolCalls      int                  `json:"tool_calls"`
	ToolFailures   int                  `json:"tool_failures"`
	FinalAnswer    string               `json:"final_answer"`
	ToolCallNames  []string             `json:"tool_call_names"`
	CheckResults   []CheckResult        `json:"check_results"`
}

// CheckResult describes one rule-based eval check.
type CheckResult struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Details string `json:"details,omitempty"`
}

// Trace is a compact event log for manual review.
type Trace struct {
	Events []TraceEvent
}

// TraceEvent is one human-readable event in an eval trace.
type TraceEvent struct {
	Time time.Time
	Text string
}
