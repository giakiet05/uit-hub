// Package mcpadapter adapts MCP server tools to the agent's native tool
// abstraction.
package mcpadapter

import (
	"context"
	"errors"
	"os"
	"os/exec"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Session is the subset of an MCP client session required by the adapter.
type Session interface {
	ListTools(ctx context.Context, params *mcp.ListToolsParams) (*mcp.ListToolsResult, error)
	CallTool(ctx context.Context, params *mcp.CallToolParams) (*mcp.CallToolResult, error)
	Close() error
}

// Client wraps one connected MCP server session.
type Client struct {
	serverName string
	session    Session
}

// ClientIdentity identifies this host application during MCP initialization.
type ClientIdentity struct {
	Name    string
	Version string
}

// StdioConfig contains settings for connecting to an MCP stdio server.
type StdioConfig struct {
	ServerName string
	Command    string
	Args       []string
	Env        []string
	WorkDir    string
	Client     ClientIdentity
}

// NewStdioClient starts an MCP stdio server process and connects to it.
func NewStdioClient(ctx context.Context, cfg StdioConfig) (*Client, error) {
	if cfg.ServerName == "" {
		return nil, errors.New("mcp server name is required")
	}
	if cfg.Command == "" {
		return nil, errors.New("mcp stdio command is required")
	}
	if cfg.Client.Name == "" {
		return nil, errors.New("mcp client name is required")
	}
	if cfg.Client.Version == "" {
		return nil, errors.New("mcp client version is required")
	}

	cmd := exec.Command(cfg.Command, cfg.Args...)
	cmd.Env = append(os.Environ(), cfg.Env...)
	cmd.Dir = cfg.WorkDir

	mcpClient := mcp.NewClient(&mcp.Implementation{
		Name:    cfg.Client.Name,
		Version: cfg.Client.Version,
	}, nil)

	session, err := mcpClient.Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		return nil, err
	}

	return NewClient(cfg.ServerName, session), nil
}

// NewClient creates an adapter for an existing MCP session.
func NewClient(serverName string, session Session) *Client {
	return &Client{
		serverName: serverName,
		session:    session,
	}
}

// Close closes the underlying MCP session.
func (c *Client) Close() error {
	return c.session.Close()
}

// Tools lists server tools and wraps them as native agent tools.
func (c *Client) Tools(ctx context.Context) ([]tool.Tool, error) {
	result, err := c.session.ListTools(ctx, nil)
	if err != nil {
		return nil, err
	}

	tools := make([]tool.Tool, 0, len(result.Tools))
	for _, serverTool := range result.Tools {
		tools = append(tools, NewTool(c.serverName, c.session, serverTool))
	}
	return tools, nil
}
