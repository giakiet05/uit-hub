# UIT Hub Agent

The UIT Hub Agent is an advanced, terminal-based AI assistant designed to interact with users, execute multi-step reasoning processes, and seamlessly interface with various external tools and systems. It serves as the core intelligent layer within the UIT Hub ecosystem, capable of performing complex code manipulation, reasoning tasks, and executing system-level operations.

The agent is built with a highly decoupled architecture, separating the core reasoning loop from the user interface and underlying infrastructure. This allows for flexible testing, robust event tracing, and future scalability.

## Core Features

- **Agent Frameworks:** Supports multiple agent reasoning loops, including `ReAct` (Reasoning and Acting) and `Plan & Execute`, allowing the agent to break down complex tasks into manageable steps.
- **Terminal User Interface (TUI):** A rich, interactive terminal interface built using `BubbleTea`, providing real-time streaming text, tool execution previews, and a session resumption picker.
- **Event-Driven Architecture:** A core Event Bus facilitates non-blocking communication between the agent logic and various background services (logging, metrics, and storage), ensuring the main reasoning loop is never blocked by I/O operations.
- **Model Context Protocol (MCP):** Native support for connecting to external MCP servers, allowing the agent to dynamically load and interact with third-party tools and data sources.
- **Contextual Memory:** Persistent memory capabilities that allow the agent to read and write context across different sessions.
- **Session Management:** Stores conversational history and token usage metrics in a local SQLite database, allowing users to seamlessly resume previous sessions.
- **Evaluation Harness:** A built-in evaluation framework to run automated test cases against the agent's logic to measure performance and correctness.

## Directory Structure

The project follows a standard Go monorepo layout, with the agent logic encapsulated within `apps/agent`.

```text
apps/agent
├── cmd
│   ├── agent/            # Main entrypoint for the Agent TUI application
│   └── eval/             # Entrypoint for running the automated evaluation harness
├── evals/                # Test cases, expected results, and rubrics for agent evaluation
└── internal/             # Private application and library code
    ├── agent/            # Core reasoning loops (ReAct, Plan & Execute) and agent interfaces
    ├── app/              # High-level application wiring, CLI parsing, and dependency injection
    ├── config/           # Environment configuration, YAML parsing, and secrets management
    ├── conversation/     # Domain models for chat messages, tool calls, and conversation history
    ├── eval/             # Execution logic and reporting for the evaluation harness
    ├── event/            # Pub/Sub Event Bus infrastructure and background Event Handlers (Logger, Storage, Usage)
    ├── hook/             # (WIP) Lifecycle hooks and guardrails for intervening during agent execution
    ├── llm/              # Interfaces and implementations for LLM providers (e.g., OpenAI)
    ├── localtool/        # Implementations of native, built-in tools (File sandbox, calculator, memory I/O)
    ├── logging/          # Centralized structured logging utilities
    ├── mcpadapter/       # Client implementation and schema mapping for the Model Context Protocol
    ├── memory/           # Persistent knowledge extraction and contextual memory storage
    ├── prompt/           # Dynamic system prompt construction and template management
    ├── session/          # Runtime session state, containing conversation history and active tools
    ├── storage/          # Database layer (SQLite + GORM) for persisting sessions and metrics
    ├── tool/             # Tool registry, execution engine, and payload budgeting constraints
    ├── tui/              # Terminal User Interface components and models (BubbleTea)
    └── usage/            # Token tracking and telemetry statistics
```
