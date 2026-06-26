package tools

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/mcp-servers/procedure/internal/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAuth registers authentication tools.
func RegisterAuth(s *server.MCPServer, c *client.Client) {

	// ── auth_login ──────────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("auth_login",
			mcp.WithDescription("Đăng nhập bằng MSSV và mật khẩu, trả về token. Gọi khi cần xác thực. Đọc dữ liệu."),
			mcp.WithString("student_id",
				mcp.Required(),
				mcp.Description("Mã số sinh viên"),
			),
			mcp.WithString("password",
				mcp.Required(),
				mcp.Description("Mật khẩu"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sid, err := requireString(req, "student_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			pwd, err := requireString(req, "password")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			body := map[string]any{
				"student_id": sid,
				"password":   pwd,
			}
			resp, apiErr := c.Post("/login", "", body)
			return handleResp(resp, apiErr), nil
		},
	)
}
