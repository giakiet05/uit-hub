# Agent Evals

This folder contains scenario specs for comparing agent loop strategies.

The first target is to compare the current ReAct baseline with future Plan-and-Execute implementations.

These evals are not unit tests yet. They are structured scenario definitions that describe:

- user input
- expected final output
- expected agent workflow
- required tools
- success criteria
- failure signals

## Layout

```text
apps/agent/evals/
  cases/
    react-chain-pass.yaml
    react-branching-fail.yaml
  rubrics/
    workflow-rubric.md
  results/
    .gitkeep
```

## How to use automatically

Run from `apps/agent` so relative MCP paths in `mcp.yaml` resolve correctly:

```bash
go run ./cmd/eval -case evals/cases/react-chain-pass.yaml
```

Run every case:

```bash
go run ./cmd/eval -cases evals/cases
```

Useful flags:

```text
-out          result report directory
-mcp-config   MCP YAML config path
-memory       eval memory directory
-max-rounds   maximum agent rounds per case
-timeout      timeout per case
```

Each run writes:

```text
evals/results/<case-id>/<timestamp>/result.json
evals/results/<case-id>/<timestamp>/trace.md
```

## How to use manually

1. Start the agent TUI.
2. Copy the `prompt` from one case file.
3. Run the prompt against the target agent type.
4. Compare the trace, logs, and final answer against `expected_workflow`, `success_criteria`, and `failure_signals`.
5. Save notable manual run notes under `results/` if needed.

## Design notes

Do not judge only the final answer. For agent evaluation, workflow matters.

The current ReAct baseline should handle long linear chains reasonably well, but it is expected to struggle with branching tasks that require planning, independent subtask tracking, and merge steps.
