package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// MCPConfig contains external MCP server settings.
type MCPConfig struct {
	ConfigPath    string
	ClientName    string
	ClientVersion string
	Servers       []MCPServerConfig
}

// MCPServerConfig contains settings for one MCP stdio server.
type MCPServerConfig struct {
	Name      string
	Transport string
	Command   string
	Args      []string
	WorkDir   string
}

type mcpConfigFile struct {
	Client  mcpClientConfigFile   `yaml:"client"`
	Servers []mcpServerConfigFile `yaml:"servers"`
}

type mcpClientConfigFile struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type mcpServerConfigFile struct {
	Name      string   `yaml:"name"`
	Transport string   `yaml:"transport"`
	Command   string   `yaml:"command"`
	Args      []string `yaml:"args"`
	WorkDir   string   `yaml:"workdir"`
}

func loadMCPConfig(path string) (MCPConfig, error) {
	cfg := MCPConfig{
		ConfigPath:    path,
		ClientName:    "uit-hub-agent",
		ClientVersion: "0.1.0",
	}
	if strings.TrimSpace(path) == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return MCPConfig{}, fmt.Errorf("read MCP config %q: %w", path, err)
	}

	var file mcpConfigFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return MCPConfig{}, fmt.Errorf("parse MCP config %q: %w", path, err)
	}

	if name := strings.TrimSpace(file.Client.Name); name != "" {
		cfg.ClientName = name
	}
	if version := strings.TrimSpace(file.Client.Version); version != "" {
		cfg.ClientVersion = version
	}

	for _, server := range file.Servers {
		name := strings.TrimSpace(server.Name)
		command := strings.TrimSpace(server.Command)
		if name == "" || command == "" {
			continue
		}
		transport := strings.TrimSpace(server.Transport)
		if transport == "" {
			transport = "stdio"
		}
		cfg.Servers = append(cfg.Servers, MCPServerConfig{
			Name:      name,
			Transport: transport,
			Command:   command,
			Args:      server.Args,
			WorkDir:   strings.TrimSpace(server.WorkDir),
		})
	}

	return cfg, nil
}
