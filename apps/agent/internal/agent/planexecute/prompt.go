package planexecute

import (
	"fmt"
	"strings"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

func plannerPrompt(userPrompt string, definitions []tool.Definition) string {
	return fmt.Sprintf(`Create a concise execution plan for this user task.

User task:
%s

Available native/runtime tools:
%s

Return only valid JSON with this shape:
{
  "goal": "short goal",
  "steps": [
    {
      "id": "step_1",
      "description": "one concrete task",
      "depends_on": [],
      "parallel_group": "group_1",
      "suggested_tools": ["tool_name"],
      "expected_output": "compact result this step must produce"
    }
  ]
}

Rules:
- Prefer 8-12 practical steps for multi-tool workflows with many ordered clauses.
- Keep every explicit user operation in the plan. Do not drop create, patch, write, read, or final answer requirements.
- Do not combine dependent operations that need generated IDs. Create, patch, downstream record creation, file write, and file read should be separate steps when one depends on the previous result.
- Use parallel_group to mark independent steps, even though this executor currently runs sequentially.
- If a required MCP tool is only listed in the MCP catalog, include load_mcp_tool before relying on that tool.
- Do not include explanation outside JSON.`, userPrompt, renderToolList(definitions))
}

func executorPrompt(userPrompt string, step planStep, state executionState) string {
	return fmt.Sprintf(`Execute exactly one planned step. Do not complete unrelated steps.

Original user task:
%s

Known execution state:
%s

Current step:
%s

Use tools if needed. When this step is complete, return only valid JSON:
{
  "step_id": %q,
  "status": "completed",
  "summary": "compact factual summary",
  "tool_calls": ["tool_name"],
  "state_updates": {
    "stable_key": "short value"
  }
}

Rules:
- Keep state_updates compact. Store IDs, names, statuses, selected course codes, numeric results, and file paths.
- Do not copy raw tool JSON into state_updates.
- Only include observations if a later step truly needs verbatim evidence.
- If a tool output is needed as an argument for another tool, call the first tool, wait for its observation, then call the dependent tool in the next model response.
- Never mark a step completed while an explicit action in the current step description is still unfinished.`, userPrompt, mustJSON(state), mustJSON(step), step.ID)
}

func replannerPrompt(userPrompt string, plan executionPlan, results []stepResult, state executionState) string {
	return fmt.Sprintf(`Repair the execution plan after a step failure.

Original user task:
%s

Current plan:
%s

Step results:
%s

Known execution state:
%s

Return only valid JSON using the same plan shape:
{
  "goal": "short goal",
  "steps": []
}`, userPrompt, mustJSON(plan), mustJSON(compactResults(results)), mustJSON(state))
}

func finalizerPrompt(userPrompt string, plan executionPlan, results []stepResult, state executionState) string {
	return fmt.Sprintf(`Write the final answer for the user.

Original user task:
%s

Plan:
%s

Known execution state:
%s

Step summaries:
%s

Rules:
- Answer in the user's language.
- Be concise.
- Do not invent IDs or facts that are not present in step results.
- Do not mention internal planner/executor details unless the user asked.`, userPrompt, mustJSON(plan), mustJSON(state), mustJSON(compactResults(results)))
}

func renderToolList(definitions []tool.Definition) string {
	if len(definitions) == 0 {
		return "(none)"
	}

	lines := make([]string, 0, len(definitions))
	for _, definition := range definitions {
		lines = append(lines, fmt.Sprintf("- %s: %s", definition.Name, definition.Description))
	}
	return strings.Join(lines, "\n")
}
