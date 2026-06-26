// Package eval runs scenario-based agent evaluations outside the TUI.
package eval

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Case defines one scenario-based agent evaluation.
type Case struct {
	ID                string        `yaml:"id" json:"id"`
	Title             string        `yaml:"title" json:"title"`
	Purpose           string        `yaml:"purpose" json:"purpose"`
	AgentType         string        `yaml:"agent_type" json:"agent_type"`
	ExpectedResult    string        `yaml:"expected_result" json:"expected_result"`
	RequiredMinRounds int           `yaml:"required_min_rounds" json:"required_min_rounds"`
	RequiredTools     RequiredTools `yaml:"required_tools" json:"required_tools"`
	Prompt            string        `yaml:"prompt" json:"prompt"`
	Prompts           []string      `yaml:"prompts" json:"prompts"`
	ExpectedWorkflow  []string      `yaml:"expected_workflow" json:"expected_workflow"`
	SuccessCriteria   []string      `yaml:"success_criteria" json:"success_criteria"`
	FailureSignals    []string      `yaml:"failure_signals" json:"failure_signals"`
	Notes             string        `yaml:"notes" json:"notes"`
}

// RequiredTools lists tool names that should appear in an eval trace.
type RequiredTools struct {
	Native []string `yaml:"native" json:"native"`
	MCP    []string `yaml:"mcp" json:"mcp"`
}

// LoadCase reads one YAML eval case from disk.
func LoadCase(path string) (Case, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Case{}, err
	}

	var testCase Case
	if err := yaml.Unmarshal(data, &testCase); err != nil {
		return Case{}, err
	}
	return testCase, nil
}
