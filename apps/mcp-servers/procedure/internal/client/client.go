// Package client wraps all HTTP calls to the fake-uit-server.
// Every method returns (map[string]any, error).
//   - Business errors (4xx): returned as map with "error_code"+"message" keys, error is nil.
//   - System errors (timeout, network): returned as (nil, err).
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Client is a thin HTTP adapter for the fake-uit-server.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// New creates a Client. It reads API_BASE_URL from env, defaulting to localhost:3000.
func New() *Client {
	base := os.Getenv("API_BASE_URL")
	if base == "" {
		base = "http://localhost:3000/api/v1"
	}
	return &Client{
		BaseURL: base,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ----- low-level helpers -----

// do executes a request and returns the parsed JSON body.
// On 2xx it returns the envelope data field.
// On 4xx it returns the full body as-is (business error).
// On 5xx or network error it returns a Go error.
func (c *Client) do(req *http.Request) (map[string]any, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var parsed map[string]any
	if len(body) > 0 {
		if err := json.Unmarshal(body, &parsed); err != nil {
			return nil, fmt.Errorf("decode json: %w", err)
		}
	}

	if resp.StatusCode >= 500 {
		msg := "internal server error"
		if m, ok := parsed["message"].(string); ok {
			msg = m
		}
		return nil, fmt.Errorf("server error %d: %s", resp.StatusCode, msg)
	}

	if resp.StatusCode >= 400 {
		// Business error – return as data, not Go error.
		if parsed == nil {
			parsed = map[string]any{}
		}
		parsed["_http_status"] = resp.StatusCode
		return parsed, nil
	}

	return parsed, nil
}

// Get issues GET with optional query params.
func (c *Client) Get(path string, token string, query map[string]string) (map[string]any, error) {
	u := c.BaseURL + path
	if len(query) > 0 {
		params := url.Values{}
		for k, v := range query {
			if v != "" {
				params.Set(k, v)
			}
		}
		u += "?" + params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.do(req)
}

// Post issues POST with a JSON body.
func (c *Client) Post(path string, token string, body any) (map[string]any, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("encode body: %w", err)
		}
	}

	req, err := http.NewRequest(http.MethodPost, c.BaseURL+path, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.do(req)
}

// ----- response helpers -----

// IsBusinessError checks if the response is a 4xx business error.
func IsBusinessError(resp map[string]any) bool {
	_, ok := resp["_http_status"]
	return ok
}

// ExtractData extracts the "data" field from a successful envelope.
func ExtractData(resp map[string]any) any {
	if d, ok := resp["data"]; ok {
		return d
	}
	return resp
}
