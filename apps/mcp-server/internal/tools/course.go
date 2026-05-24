package tools

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/mcp-server/internal/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterCourse registers course/assignment-related tools.
func RegisterCourse(s *server.MCPServer, c *client.Client) {
	tokenDesc := "Bearer token nhận từ auth_login"

	// ── get_deadlines ───────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_deadlines",
			mcp.WithDescription("Lấy danh sách deadline bài tập sắp tới của sinh viên."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			resp, apiErr := c.Get("/student/deadlines", token, nil)
			return handleResp(resp, apiErr), nil
		},
	)

	// ── get_course_materials ────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_course_materials",
			mcp.WithDescription("Lấy danh sách tài liệu (slide, PDF, notebook) của một môn học."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
			mcp.WithString("course_id", mcp.Required(), mcp.Description("Mã môn học (VD: IT010, MA006)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			courseID, err := requireString(req, "course_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			resp, apiErr := c.Get("/student/courses/"+courseID+"/materials", token, nil)
			return handleResp(resp, apiErr), nil
		},
	)

	// ── get_course_assignments ──────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_course_assignments",
			mcp.WithDescription("Lấy danh sách bài tập/đồ án của một môn học."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
			mcp.WithString("course_id", mcp.Required(), mcp.Description("Mã môn học (VD: IT010, MA006)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			courseID, err := requireString(req, "course_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			resp, apiErr := c.Get("/student/courses/"+courseID+"/assignments", token, nil)
			return handleResp(resp, apiErr), nil
		},
	)

	// ── submit_assignment ───────────────────────────────────────
	s.AddTool(
		mcp.NewTool("submit_assignment",
			mcp.WithDescription("Nộp bài tập. Hỗ trợ nộp text, link hoặc file. Trả về submission_id."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
			mcp.WithString("assignment_id", mcp.Required(), mcp.Description("Mã bài tập (VD: ASG001)")),
			mcp.WithString("submission_type", mcp.Required(), mcp.Description("Loại nộp bài"), mcp.Enum("FILE", "LINK", "TEXT")),
			mcp.WithString("text", mcp.Description("Nội dung text (khi type=TEXT)")),
			mcp.WithString("url", mcp.Description("URL bài nộp (khi type=LINK)")),
			mcp.WithString("comment", mcp.Description("Ghi chú cho giảng viên")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			assignmentID, err := requireString(req, "assignment_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			body := buildBody(req, "submission_type", "text", "url", "file_urls", "comment")
			resp, apiErr := c.Post("/student/assignments/"+assignmentID+"/submissions", token, body)
			return handleResp(resp, apiErr), nil
		},
	)
}
