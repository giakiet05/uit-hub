package tools

import (
	"context"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/mcp-servers/academic/internal/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAcademic registers academic-domain query and mutation tools.
func RegisterAcademic(s *server.MCPServer, c *client.Client) {
	tokenDesc := "Bearer token nhận từ auth_login"

	yearSemesterOpts := func() []mcp.ToolOption {
		return []mcp.ToolOption{
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
			mcp.WithNumber("year", mcp.Required(), mcp.Description("Năm học (VD: 2024)")),
			mcp.WithNumber("semester", mcp.Required(), mcp.Description("Học kỳ (1, 2, hoặc 3)")),
		}
	}

	// ── get_student_profile ─────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_student_profile",
			mcp.WithDescription("Lấy thông tin cá nhân của sinh viên đang đăng nhập (họ tên, MSSV, email, SĐT, ngành). Đọc dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			resp, apiErr := c.Get("/student/profile", token, nil)
			return handleResp(resp, apiErr), nil
		},
	)

	// ── get_student_schedule ────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_student_schedule",
			append([]mcp.ToolOption{
				mcp.WithDescription("Lấy thời khoá biểu (lịch học) của sinh viên theo năm và học kỳ. Đọc dữ liệu."),
			}, yearSemesterOpts()...)...,
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			year, _ := optionalNumber(req, "year")
			semester, _ := optionalNumber(req, "semester")
			resp, apiErr := c.Get("/student/schedule", token, map[string]string{
				"year":     fmt.Sprintf("%.0f", year),
				"semester": fmt.Sprintf("%.0f", semester),
			})
			return handleResp(resp, apiErr), nil
		},
	)

	// ── get_exam_schedule ───────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_exam_schedule",
			append([]mcp.ToolOption{
				mcp.WithDescription("Lấy lịch thi của sinh viên theo năm và học kỳ. Đọc dữ liệu."),
			}, yearSemesterOpts()...)...,
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			year, _ := optionalNumber(req, "year")
			semester, _ := optionalNumber(req, "semester")
			resp, apiErr := c.Get("/student/schedule/exam", token, map[string]string{
				"year":     fmt.Sprintf("%.0f", year),
				"semester": fmt.Sprintf("%.0f", semester),
			})
			return handleResp(resp, apiErr), nil
		},
	)

	// ── get_student_scores ──────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_student_scores",
			mcp.WithDescription("Lấy toàn bộ bảng điểm (giữa kỳ, cuối kỳ, tổng) của sinh viên. Đọc dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			resp, apiErr := c.Get("/student/score", token, nil)
			return handleResp(resp, apiErr), nil
		},
	)

	// ── get_enrolled_courses ────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_enrolled_courses",
			mcp.WithDescription("Lấy danh sách các môn học mà sinh viên đã và đang đăng ký. Đọc dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			resp, apiErr := c.Get("/student/courses", token, nil)
			return handleResp(resp, apiErr), nil
		},
	)

	// ── get_training_points ─────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_training_points",
			mcp.WithDescription("Tra cứu điểm rèn luyện của sinh viên (điểm, xếp loại). Đọc dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			resp, apiErr := c.Get("/student/training-points", token, nil)
			return handleResp(resp, apiErr), nil
		},
	)

	// ── get_survey_forms ────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_survey_forms",
			mcp.WithDescription("Lấy danh sách phiếu khảo sát sinh viên cần hoàn thành. Đọc dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			resp, apiErr := c.Get("/student/survey-form", token, nil)
			return handleResp(resp, apiErr), nil
		},
	)

	// ── get_deadlines ───────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_deadlines",
			mcp.WithDescription("Lấy danh sách deadline bài tập sắp tới của sinh viên. Đọc dữ liệu."),
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
			mcp.WithDescription("Lấy danh sách tài liệu (slide, PDF, notebook) của một môn học. Đọc dữ liệu."),
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
			mcp.WithDescription("Lấy danh sách bài tập/đồ án của một môn học. Đọc dữ liệu."),
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
			mcp.WithDescription("Nộp bài tập. Hỗ trợ nộp text, link hoặc file. Trả về submission_id. Ghi dữ liệu."),
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
