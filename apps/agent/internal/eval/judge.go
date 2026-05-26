package eval

import (
	"fmt"
	"slices"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/agent"
)

// Judge runs deterministic checks over an eval result.
func Judge(testCase Case, result Result) []CheckResult {
	checks := []CheckResult{
		{
			Name:    "terminal_completed",
			Passed:  result.TerminalReason == agent.TerminalCompleted,
			Details: fmt.Sprintf("terminal_reason=%s error=%q", result.TerminalReason, result.Error),
		},
		{
			Name:    "minimum_rounds",
			Passed:  result.RoundCount >= testCase.RequiredMinRounds,
			Details: fmt.Sprintf("rounds=%d required=%d", result.RoundCount, testCase.RequiredMinRounds),
		},
		{
			Name:    "final_answer_present",
			Passed:  strings.TrimSpace(result.FinalAnswer) != "",
			Details: fmt.Sprintf("chars=%d", len(result.FinalAnswer)),
		},
	}

	for _, name := range testCase.RequiredTools.Native {
		checks = append(checks, toolWasCalledCheck("native_tool_called:"+name, name, result.ToolCallNames))
	}
	for _, name := range testCase.RequiredTools.MCP {
		checks = append(checks, toolWasCalledCheck("mcp_tool_called:"+name, name, result.ToolCallNames))
	}

	return checks
}

func toolWasCalledCheck(checkName string, toolName string, calls []string) CheckResult {
	return CheckResult{
		Name:    checkName,
		Passed:  slices.Contains(calls, toolName),
		Details: fmt.Sprintf("tool=%s calls=%s", toolName, strings.Join(calls, ",")),
	}
}
