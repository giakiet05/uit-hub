package planexecute

type executionPlan struct {
	Goal  string     `json:"goal"`
	Steps []planStep `json:"steps"`
}

type planStep struct {
	ID             string   `json:"id"`
	Description    string   `json:"description"`
	DependsOn      []string `json:"depends_on"`
	ParallelGroup  string   `json:"parallel_group"`
	SuggestedTools []string `json:"suggested_tools"`
	ExpectedOutput string   `json:"expected_output"`
}

type stepResult struct {
	StepID       string            `json:"step_id"`
	Status       string            `json:"status"`
	Summary      string            `json:"summary"`
	ToolCalls    []string          `json:"tool_calls"`
	Observations []string          `json:"observations,omitempty"`
	StateUpdates map[string]string `json:"state_updates,omitempty"`
	Error        string            `json:"error,omitempty"`
}

type compactStepResult struct {
	StepID       string            `json:"step_id"`
	Status       string            `json:"status"`
	Summary      string            `json:"summary"`
	ToolCalls    []string          `json:"tool_calls,omitempty"`
	StateUpdates map[string]string `json:"state_updates,omitempty"`
	Error        string            `json:"error,omitempty"`
}

type executionState struct {
	Facts map[string]string `json:"facts"`
}
