# Student Support AI Agent

Student Support AI Agent is an experimental Go-based AI agent prototype for student-support workflows. It focuses on agent architecture: LLM provider abstraction, tool calling, MCP integration, memory, session persistence, context management, and evaluation.

This project is not a production system. It is a learning and research project for understanding how AI agents can be built from lower-level components instead of relying on an agent framework.

## Current Status

- Prototype stage.
- Main interface is a terminal UI for fast local testing.
- The agent currently uses OpenAI-compatible models.
- MCP support is available through a local client adapter.
- Some features are intentionally simple and may change as the architecture evolves.

## Features

- ReAct agent loop for tool-using conversations.
- Plan-and-Execute baseline for comparing agent behavior on longer tasks.
- LLM provider abstraction with an OpenAI implementation.
- MCP client adapter for connecting external tool servers.
- Lazy MCP tool loading with a runtime tool registry.
- Local tools for development, memory, file sandboxing, and calculation.
- SQLite-backed session persistence and resume support.
- Persistent conversation history, token usage, and restored runtime tools.
- Prompt construction with stable and session-specific sections to improve prompt caching.
- Memory tools for storing and reading simple user preferences or notes.
- Context management with result budgeting, micro-compaction, snip compaction, and context summaries.
- Interrupt handling for stopping agent runs while allowing approved write tools to finish safely.
- Evaluation runner for testing agent workflows and comparing agent loop designs.
- Terminal UI for inspecting conversations, logs, tool calls, compaction events, and usage metrics.

## Architecture

The code is organized around a few main layers:

```text
app
  Loads configuration, creates shared services, connects MCP servers, and starts the selected interface.

session
  Owns one conversation session: conversation state, runtime tools, usage, memory context, and persistence.

agent
  Runs one agent loop for one user request. Implementations include ReAct and Plan-and-Execute.

tool
  Defines tool interfaces, registries, execution, result budgeting, and metadata.

mcpadapter
  Converts MCP server tools into native agent tools.

llm
  Defines the model/provider interface and provider implementations.
```

The intended direction is:

```text
app -> session -> agent
              -> tool services
              -> storage / memory / prompt context
```

Inner layers should avoid depending on outer layers. This keeps the agent loop easier to test and makes it possible to reuse the core runtime in a future backend service.

## Project Layout

```text
apps/agent
├── cmd
│   ├── agent/        # Terminal UI entrypoint
│   └── eval/         # Evaluation runner entrypoint
├── evals/            # Evaluation cases and reports
├── internal/
│   ├── agent/        # Agent interfaces and loop implementations
│   ├── app/          # Application wiring
│   ├── config/       # Env and YAML configuration
│   ├── conversation/ # Messages, tool calls, and compaction helpers
│   ├── eval/         # Evaluation runner logic
│   ├── event/        # Internal event bus and log handlers
│   ├── llm/          # LLM provider interfaces and implementations
│   ├── localtool/    # Local development tools
│   ├── mcpadapter/   # MCP client adapter
│   ├── memory/       # File-backed memory store
│   ├── prompt/       # System prompt construction
│   ├── session/      # Session orchestration and state
│   ├── storage/      # SQLite persistence
│   ├── tool/         # Tool registry and executor
│   ├── tui/          # Terminal UI
│   └── usage/        # Token usage data types
└── testdata/         # Test fixtures

apps/mock-mcp-server
└── cmd/server/       # Mock MCP server used for local testing
```

## Getting Started

Prerequisites:

- Go 1.25+
- An OpenAI API key or an OpenAI-compatible endpoint

Create local config files from the examples:

```bash
cd apps/agent
cp .env.example .env
cp mcp.example.yaml mcp.yaml
```

Edit the local `.env` file with your model settings and API key.

Run the terminal agent. If `mcp.yaml` points to the mock MCP server, the agent will connect to it during startup.

```bash
go run ./cmd/agent
```

Build a local binary:

```bash
go build -o agent ./cmd/agent
```

## Evaluation

The eval runner executes predefined cases and records agent behavior for comparison.

```bash
cd apps/agent
go run ./cmd/eval
```

Evaluation is currently used as a development tool, not as a full benchmark suite.

## Notes

- This repository is still evolving.
- The current terminal UI is mainly for development and debugging.
- The long-term goal is to reuse the agent core behind a backend service.
- The mock MCP server exists only to test tool-calling behavior before real student-support systems are connected.
