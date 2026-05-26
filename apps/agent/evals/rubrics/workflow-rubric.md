# Workflow Rubric

Use this rubric when manually reviewing an eval run.

## Pass Criteria

- The agent uses tools when source-system data is required.
- The agent does not invent student records, IDs, schedules, notes, tickets, or generated artifacts.
- The agent respects dependencies between steps.
- The agent does not perform write operations before collecting required verified data.
- The final answer includes the required stable identifiers, paths, or key values.
- Tool failures are acknowledged or recovered from.
- The workflow does not exceed the configured round limit unexpectedly.

## ReAct-Specific Signals

Good ReAct behavior:

- Executes a clear linear chain.
- Uses each observation to decide the next action.
- Recovers when a tool needs to be loaded before use.

Weak ReAct behavior:

- Forgets independent branches.
- Creates records too early before all required data exists.
- Fails to merge data from multiple branches.
- Repeats tool calls without new information.
- Stops early after satisfying only part of the prompt.

## Plan-and-Execute Signals

Good Plan-and-Execute behavior:

- Produces or follows an implicit task decomposition.
- Tracks multiple independent branches.
- Executes dependent steps only after prerequisites are available.
- Performs a final verification pass before answering.

## Suggested Manual Score

Use a 0-5 score:

- 0: Did not understand the task.
- 1: Used tools incorrectly or invented key data.
- 2: Completed a small part but missed major requirements.
- 3: Mostly correct but missed one important step or identifier.
- 4: Correct result with minor inefficiency or formatting issue.
- 5: Correct result and expected workflow.
