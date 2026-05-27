package tools

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/mcp-server/internal/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterRoom registers room-related tools.
func RegisterRoom(s *server.MCPServer, c *client.Client) {
	tokenDesc := "Bearer token nhận từ auth_login"

	// ── get_rooms_availability ──────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_rooms_availability",
			mcp.WithDescription("Tra cứu danh sách phòng trống trong khoảng thời gian chỉ định. Đọc dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
			mcp.WithString("date", mcp.Required(), mcp.Description("Ngày cần tra (YYYY-MM-DD)")),
			mcp.WithString("start", mcp.Required(), mcp.Description("Giờ bắt đầu (HH:MM, VD: 08:00)")),
			mcp.WithString("end", mcp.Required(), mcp.Description("Giờ kết thúc (HH:MM, VD: 10:00)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			date := optionalString(req, "date")
			start := optionalString(req, "start")
			end := optionalString(req, "end")

			resp, apiErr := c.Get("/rooms/availability", token, map[string]string{
				"date":  date,
				"start": start,
				"end":   end,
			})
			return handleResp(resp, apiErr), nil
		},
	)

	// ── list_students ───────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("list_students",
			mcp.WithDescription("Lấy danh sách tất cả sinh viên (public, không cần token). Đọc dữ liệu."),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			resp, apiErr := c.Get("/students", "", nil)
			return handleResp(resp, apiErr), nil
		},
	)
}
