// Package tools registers all MCP tools and their handlers.
package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/mcp-servers/academic/internal/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// resultJSON turns any Go value into a JSON text MCP result.
func resultJSON(v any) *mcp.CallToolResult {
	b, err := json.Marshal(v)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("marshal error: %v", err))
	}
	return mcp.NewToolResultText(string(b))
}

// handleResp is the central dispatcher:
//   - Go error (network/5xx) → MCP error (isError=true).
//   - Business error (4xx)   → structured JSON content (isError=false so agent sees the message).
//   - Success                → extract "data" and return as JSON.
func handleResp(resp map[string]any, err error) *mcp.CallToolResult {
	if err != nil {
		return mcp.NewToolResultError(err.Error())
	}
	if client.IsBusinessError(resp) {
		status := resp["_http_status"]
		delete(resp, "_http_status")
		out := map[string]any{
			"ok":          false,
			"http_status": status,
		}
		if code, ok := resp["error_code"].(string); ok {
			out["error_code"] = code
		}
		if msg, ok := resp["message"].(string); ok {
			out["message"] = msg
		}
		return resultJSON(out)
	}
	data := client.ExtractData(resp)
	out := map[string]any{
		"ok":   true,
		"data": data,
	}
	return resultJSON(out)
}

// getArgs returns the arguments map from a CallToolRequest.
func getArgs(req mcp.CallToolRequest) map[string]any {
	return req.GetArguments()
}

// requireString extracts a required string arg, returning a user-friendly error.
func requireString(req mcp.CallToolRequest, key string) (string, error) {
	args := getArgs(req)
	v, ok := args[key]
	if !ok || v == nil {
		return "", fmt.Errorf("%s is required", key)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string", key)
	}
	return s, nil
}

// optionalString extracts an optional string arg.
func optionalString(req mcp.CallToolRequest, key string) string {
	args := getArgs(req)
	v, ok := args[key]
	if !ok || v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

// optionalNumber extracts an optional numeric arg.
func optionalNumber(req mcp.CallToolRequest, key string) (float64, bool) {
	args := getArgs(req)
	v, ok := args[key]
	if !ok || v == nil {
		return 0, false
	}
	n, ok := v.(float64)
	return n, ok
}

// buildBody extracts specified keys from args into a map for POST body.
func buildBody(req mcp.CallToolRequest, keys ...string) map[string]any {
	args := getArgs(req)
	body := map[string]any{}
	for _, k := range keys {
		if v, ok := args[k]; ok && v != nil {
			body[k] = v
		}
	}
	return body
}

// Handler is shorthand for the MCP handler signature.
type Handler func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)
