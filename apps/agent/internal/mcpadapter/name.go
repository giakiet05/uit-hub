package mcpadapter

import (
	"fmt"
	"strings"
)

const namespaceSeparator = "__"

// NamespacedName prefixes an MCP tool name with its server namespace.
func NamespacedName(serverName string, toolName string) string {
	return serverName + namespaceSeparator + toolName
}

// SplitNamespacedName splits a namespaced tool name into server and tool names.
func SplitNamespacedName(name string) (string, string, error) {
	serverName, toolName, ok := strings.Cut(name, namespaceSeparator)
	if !ok || serverName == "" || toolName == "" {
		return "", "", fmt.Errorf("invalid namespaced tool name %q", name)
	}
	return serverName, toolName, nil
}
