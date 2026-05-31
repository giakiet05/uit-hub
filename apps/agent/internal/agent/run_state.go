package agent

import "time"

// RunState owns mutable state for one Agent.Run invocation.
type RunState struct {
	SessionID string
	StartedAt time.Time
	Round     int
	MaxRounds int
	Stats     *RunStats
}

// NewRunState creates run-scoped state with fresh counters.
func NewRunState(sessionID string, maxRounds int) *RunState {
	return &RunState{
		SessionID: sessionID,
		StartedAt: time.Now().UTC(),
		MaxRounds: maxRounds,
		Stats:     NewRunStats(),
	}
}

// Finish records the final duration for the run stats.
func (s *RunState) Finish() {
	if s == nil || s.Stats == nil {
		return
	}
	s.Stats.Finish()
}
