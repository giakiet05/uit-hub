package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const serverVersion = "0.1.0"

type pingInput struct{}

type pingOutput struct {
	OK      bool   `json:"ok" jsonschema:"whether the server is alive"`
	Message string `json:"message" jsonschema:"status message"`
}

func main() {
	attempt, err := nextAttempt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "mock retry mcp server attempt tracking failed: %v\n", err)
		os.Exit(1)
	}
	if attempt < 3 {
		fmt.Fprintf(os.Stderr, "mock retry mcp server intentionally failed on attempt %d\n", attempt)
		os.Exit(1)
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "uit-hub-retry-mcp-server",
		Version: serverVersion,
	}, nil)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "retry_ping",
		Description: "Return a successful response after retry startup failures.",
	}, ping)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("mock retry mcp server failed: %v", err)
	}
}

func nextAttempt() (int, error) {
	path := attemptFilePath()
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return 0, err
	}

	attempt := 0
	text := strings.TrimSpace(string(data))
	if text != "" {
		parsed, err := strconv.Atoi(text)
		if err != nil {
			return 0, err
		}
		attempt = parsed
	}
	attempt++

	if err := os.WriteFile(path, []byte(strconv.Itoa(attempt)), 0o600); err != nil {
		return 0, err
	}
	return attempt, nil
}

func attemptFilePath() string {
	if path := strings.TrimSpace(os.Getenv("UIT_HUB_RETRY_MCP_ATTEMPT_FILE")); path != "" {
		return path
	}
	if loadID := strings.TrimSpace(os.Getenv("UIT_HUB_MCP_LOAD_ID")); loadID != "" {
		return filepath.Join(os.TempDir(), "uit-hub-retry-mcp-attempt-"+safeFilePart(loadID))
	}
	return filepath.Join(os.TempDir(), "uit-hub-retry-mcp-attempt")
}

func safeFilePart(value string) string {
	var builder strings.Builder
	for _, char := range value {
		if char >= 'a' && char <= 'z' ||
			char >= 'A' && char <= 'Z' ||
			char >= '0' && char <= '9' ||
			char == '-' ||
			char == '_' {
			builder.WriteRune(char)
			continue
		}
		builder.WriteByte('_')
	}
	return builder.String()
}

func ping(ctx context.Context, req *mcp.CallToolRequest, input pingInput) (*mcp.CallToolResult, pingOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, pingOutput{}, err
	}
	return nil, pingOutput{
		OK:      true,
		Message: "retry server connected",
	}, nil
}
