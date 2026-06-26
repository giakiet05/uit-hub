package tools

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/mcp-servers/campus/internal/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterCampus registers campus-domain query and mutation tools.
func RegisterCampus(s *server.MCPServer, c *client.Client) {
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

	// ── create_contact ──────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("create_contact",
			mcp.WithDescription("Gửi form liên hệ/phản hồi tới trường. Không cần token. Ghi dữ liệu."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Họ tên người gửi")),
			mcp.WithString("email", mcp.Required(), mcp.Description("Email liên hệ")),
			mcp.WithString("phone", mcp.Description("SĐT")),
			mcp.WithString("subject", mcp.Required(), mcp.Description("Tiêu đề")),
			mcp.WithString("message", mcp.Required(), mcp.Description("Nội dung")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			body := buildBody(req, "name", "email", "phone", "subject", "message")
			resp, apiErr := c.Post("/contact", "", body)
			return handleResp(resp, apiErr), nil
		},
	)
}
