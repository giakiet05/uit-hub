package tools

import (
	"context"
	"fmt"

	"github.com/giakiet05/uit-hub/apps/mcp-server/internal/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterStudent registers student-domain query and mutation tools.
func RegisterStudent(s *server.MCPServer, c *client.Client) {
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

	// ── get_tuition_fee ─────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_tuition_fee",
			mcp.WithDescription("Tra cứu học phí hiện tại của sinh viên (số tín chỉ, số tiền, trạng thái). Đọc dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			resp, apiErr := c.Get("/student/lookup/tuitionfee", token, nil)
			return handleResp(resp, apiErr), nil
		},
	)

	// ── get_insurance ───────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("get_insurance",
			mcp.WithDescription("Tra cứu thông tin bảo hiểm y tế của sinh viên. Đọc dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			resp, apiErr := c.Get("/student/insurance", token, nil)
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

	// ── create_transcript_request ───────────────────────────────
	s.AddTool(
		mcp.NewTool("create_transcript_request",
			mcp.WithDescription("Đăng ký in bảng điểm. Trả về request_id và trạng thái. Ghi dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
			mcp.WithNumber("copies", mcp.Required(), mcp.Description("Số bản in (≥1)")),
			mcp.WithString("language", mcp.Required(), mcp.Description("Ngôn ngữ bảng điểm"), mcp.Enum("VI", "EN")),
			mcp.WithString("delivery_method", mcp.Required(), mcp.Description("Phương thức nhận"), mcp.Enum("PICKUP", "SHIP")),
			mcp.WithString("shipping_address", mcp.Description("Địa chỉ giao (bắt buộc nếu SHIP)")),
			mcp.WithString("phone", mcp.Description("SĐT nhận hàng (bắt buộc nếu SHIP)")),
			mcp.WithString("note", mcp.Description("Ghi chú thêm")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			body := buildBody(req, "copies", "language", "delivery_method", "shipping_address", "phone", "note")
			resp, apiErr := c.Post("/student/transcript-regis", token, body)
			return handleResp(resp, apiErr), nil
		},
	)

	// ── create_tuition_extension ────────────────────────────────
	s.AddTool(
		mcp.NewTool("create_tuition_extension",
			mcp.WithDescription("Xin gia hạn nộp học phí. Trả về request_id và trạng thái. Ghi dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
			mcp.WithNumber("year", mcp.Required(), mcp.Description("Năm học")),
			mcp.WithNumber("semester", mcp.Required(), mcp.Description("Học kỳ")),
			mcp.WithString("requested_due_date", mcp.Required(), mcp.Description("Ngày xin gia hạn (YYYY-MM-DD)")),
			mcp.WithString("reason", mcp.Required(), mcp.Description("Lý do xin gia hạn")),
			mcp.WithString("phone", mcp.Description("SĐT liên hệ")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			body := buildBody(req, "year", "semester", "requested_due_date", "reason", "phone")
			resp, apiErr := c.Post("/student/tuition-extend", token, body)
			return handleResp(resp, apiErr), nil
		},
	)

	// ── create_monthly_parking ──────────────────────────────────
	s.AddTool(
		mcp.NewTool("create_monthly_parking",
			mcp.WithDescription("Đăng ký vé xe tháng. Trả về request_id. Ghi dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
			mcp.WithString("vehicle_type", mcp.Required(), mcp.Description("Loại xe"), mcp.Enum("MOTORBIKE", "CAR", "BICYCLE")),
			mcp.WithString("plate_number", mcp.Required(), mcp.Description("Biển số xe")),
			mcp.WithNumber("months", mcp.Required(), mcp.Description("Số tháng đăng ký (1-12)")),
			mcp.WithString("start_month", mcp.Required(), mcp.Description("Tháng bắt đầu (YYYY-MM)")),
			mcp.WithString("owner_name", mcp.Description("Tên chủ xe")),
			mcp.WithString("note", mcp.Description("Ghi chú")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			body := buildBody(req, "vehicle_type", "plate_number", "months", "start_month", "owner_name", "note")
			resp, apiErr := c.Post("/student/monthly-parking", token, body)
			return handleResp(resp, apiErr), nil
		},
	)

	// ── create_graduate_request ─────────────────────────────────
	s.AddTool(
		mcp.NewTool("create_graduate_request",
			mcp.WithDescription("Đăng ký xét tốt nghiệp. Trả về request_id. Ghi dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
			mcp.WithNumber("year", mcp.Required(), mcp.Description("Năm học")),
			mcp.WithNumber("semester", mcp.Required(), mcp.Description("Học kỳ")),
			mcp.WithString("email", mcp.Required(), mcp.Description("Email liên hệ")),
			mcp.WithString("phone", mcp.Required(), mcp.Description("SĐT liên hệ")),
			mcp.WithString("address", mcp.Description("Địa chỉ")),
			mcp.WithString("note", mcp.Description("Ghi chú")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			body := buildBody(req, "year", "semester", "email", "phone", "address", "note")
			resp, apiErr := c.Post("/student/graduate", token, body)
			return handleResp(resp, apiErr), nil
		},
	)

	// ── create_graduation_thesis_request ────────────────────────
	s.AddTool(
		mcp.NewTool("create_graduation_thesis_request",
			mcp.WithDescription("Đăng ký khóa luận tốt nghiệp. Trả về request_id. Ghi dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
			mcp.WithString("thesis_title", mcp.Required(), mcp.Description("Tên đề tài")),
			mcp.WithString("advisor_name", mcp.Required(), mcp.Description("Tên giảng viên hướng dẫn")),
			mcp.WithString("advisor_email", mcp.Description("Email giảng viên hướng dẫn")),
			mcp.WithString("note", mcp.Description("Ghi chú")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			body := buildBody(req, "thesis_title", "advisor_name", "advisor_email", "note", "team_members")
			resp, apiErr := c.Post("/student/graduation-thesis", token, body)
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
