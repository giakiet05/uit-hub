package mcpadapter

import (
	"errors"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/agent/internal/tool"
)

// ToolMetadata is the compact MCP tool catalog entry shown to the model.
type ToolMetadata struct {
	Name              string
	ServerName        string
	ServerDescription string
	Description       string
}

type managedTool struct {
	tool     tool.Tool
	metadata ToolMetadata
}

// Manager owns connected MCP clients and the compact tool catalog built at
// session start.
type Manager struct {
	clients map[string]*Client
	tools   map[string]managedTool
	order   []string
}

// NewManager creates an empty MCP manager.
func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]*Client),
		tools:   make(map[string]managedTool),
		order:   []string{},
	}
}

// AddLoadedTools records one connected MCP client and its already-listed tools.
func (m *Manager) AddLoadedTools(serverName string, serverDescription string, client *Client, loadedTools []tool.Tool) error {
	if m == nil {
		return errors.New("add mcp tools: nil manager")
	}
	if serverName == "" {
		return errors.New("add mcp tools: empty server name")
	}
	if client == nil {
		return fmt.Errorf("add mcp tools for %q: nil client", serverName)
	}
	if _, exists := m.clients[serverName]; exists {
		return fmt.Errorf("add mcp tools for %q: duplicate server", serverName)
	}

	pending := make([]managedTool, 0, len(loadedTools))
	seen := make(map[string]struct{}, len(loadedTools))
	for _, loadedTool := range loadedTools {
		if loadedTool == nil {
			return fmt.Errorf("add mcp tools for %q: nil tool", serverName)
		}
		definition := loadedTool.Definition()
		if definition.Name == "" {
			return fmt.Errorf("add mcp tools for %q: empty tool name", serverName)
		}
		if _, exists := m.tools[definition.Name]; exists {
			return fmt.Errorf("add mcp tools for %q: duplicate tool %q", serverName, definition.Name)
		}
		if _, exists := seen[definition.Name]; exists {
			return fmt.Errorf("add mcp tools for %q: duplicate tool %q", serverName, definition.Name)
		}
		seen[definition.Name] = struct{}{}
		pending = append(pending, managedTool{
			tool: loadedTool,
			metadata: ToolMetadata{
				Name:              definition.Name,
				ServerName:        serverName,
				ServerDescription: serverDescription,
				Description:       definition.Description,
			},
		})
	}

	m.clients[serverName] = client
	for _, entry := range pending {
		m.tools[entry.metadata.Name] = entry
		m.order = append(m.order, entry.metadata.Name)
	}
	return nil
}

// Catalog returns compact MCP tool metadata in stable registration order.
func (m *Manager) Catalog() []ToolMetadata {
	if m == nil {
		return nil
	}

	catalog := make([]ToolMetadata, 0, len(m.order))
	for _, name := range m.order {
		catalog = append(catalog, m.tools[name].metadata)
	}
	return catalog
}

// LoadTool returns the full native tool wrapper for one namespaced MCP tool.
func (m *Manager) LoadTool(name string) (tool.Tool, bool) {
	if m == nil {
		return nil, false
	}

	managed, exists := m.tools[name]
	return managed.tool, exists
}

// Close closes every connected MCP client.
func (m *Manager) Close() error {
	if m == nil {
		return nil
	}

	var closeErr error
	for _, client := range m.clients {
		if err := client.Close(); err != nil {
			closeErr = errors.Join(closeErr, err)
		}
	}
	return closeErr
}
