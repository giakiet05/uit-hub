package tools

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/mcp-servers/procedure/internal/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterProcedure registers procedure-domain query and mutation tools.
func RegisterProcedure(s *server.MCPServer, c *client.Client) {
	tokenDesc := "Bearer token nhận từ auth_login"

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

	// ── create_confirm_letter ───────────────────────────────────
	s.AddTool(
		mcp.NewTool("create_confirm_letter",
			mcp.WithDescription("Đăng ký giấy xác nhận sinh viên (hoãn nghĩa vụ QS, ký túc xá, giảm thuế, v.v.). Ghi dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
			mcp.WithString("language", mcp.Required(), mcp.Description("Ngôn ngữ"), mcp.Enum("VI", "EN")),
			mcp.WithString("reason", mcp.Required(), mcp.Description("Lý do xin xác nhận"),
				mcp.Enum("MILITARY_DEFERMENT", "DORM_EXTEND", "TAX_DEDUCTION_DOCS", "DEFENSE_EDU_REGISTRATION", "OTHER"),
			),
			mcp.WithString("other_reason", mcp.Description("Chi tiết lý do (bắt buộc nếu reason=OTHER)")),
			mcp.WithString("request_type", mcp.Required(), mcp.Description("Loại yêu cầu"), mcp.Enum("NEW", "REISSUE")),
			mcp.WithString("note", mcp.Description("Ghi chú")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			body := buildBody(req, "language", "reason", "other_reason", "request_type", "note")
			resp, apiErr := c.Post("/student/confirm-letter", token, body)
			return handleResp(resp, apiErr), nil
		},
	)

	// ── create_bank_loan_request ────────────────────────────────
	s.AddTool(
		mcp.NewTool("create_bank_loan_request",
			mcp.WithDescription("Đăng ký xác nhận vay vốn ngân hàng cho sinh viên. Ghi dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
			mcp.WithString("benefit", mcp.Required(), mcp.Description("Chế độ ưu đãi"),
				mcp.Enum("NO_DISCOUNT", "TUITION_REDUCTION", "TUITION_EXEMPTION"),
			),
			mcp.WithString("orphan_status", mcp.Required(), mcp.Description("Tình trạng mồ côi"),
				mcp.Enum("NOT_ORPHAN", "ORPHAN"),
			),
			mcp.WithString("template", mcp.Required(), mcp.Description("Mẫu đơn"),
				mcp.Enum("LEGACY", "STEM"),
			),
			mcp.WithString("note", mcp.Description("Ghi chú")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			body := buildBody(req, "benefit", "orphan_status", "template", "note")
			resp, apiErr := c.Post("/student/bank-loans", token, body)
			return handleResp(resp, apiErr), nil
		},
	)

	// ── create_training_point_confirm ───────────────────────────
	s.AddTool(
		mcp.NewTool("create_training_point_confirm",
			mcp.WithDescription("Đăng ký giấy xác nhận điểm rèn luyện. Ghi dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
			mcp.WithString("language", mcp.Required(), mcp.Description("Ngôn ngữ"), mcp.Enum("VI", "EN")),
			mcp.WithString("note", mcp.Description("Ghi chú")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			body := buildBody(req, "language", "note")
			resp, apiErr := c.Post("/student/training-point-confirm", token, body)
			return handleResp(resp, apiErr), nil
		},
	)

	// ── upload_language_certificate ─────────────────────────────
	s.AddTool(
		mcp.NewTool("upload_language_certificate",
			mcp.WithDescription("Nộp chứng chỉ ngoại ngữ (TOEIC, IELTS, v.v.). Ghi dữ liệu."),
			mcp.WithString("token", mcp.Required(), mcp.Description(tokenDesc)),
			mcp.WithString("document_type", mcp.Required(), mcp.Description("Loại chứng chỉ"),
				mcp.Enum("CONFIRMATION", "DIPLOMA", "CERTIFICATE"),
			),
			mcp.WithString("birth_date", mcp.Required(), mcp.Description("Ngày sinh (YYYY-MM-DD)")),
			mcp.WithString("id_number", mcp.Required(), mcp.Description("Số CCCD/CMND")),
			mcp.WithNumber("listening_score", mcp.Required(), mcp.Description("Điểm nghe")),
			mcp.WithNumber("reading_score", mcp.Required(), mcp.Description("Điểm đọc")),
			mcp.WithNumber("total_score", mcp.Required(), mcp.Description("Tổng điểm")),
			mcp.WithString("exam_date", mcp.Required(), mcp.Description("Ngày thi (YYYY-MM-DD)")),
			mcp.WithString("image_file", mcp.Required(), mcp.Description("Ảnh chứng chỉ (base64)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			token, err := requireString(req, "token")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			body := buildBody(req, "document_type", "birth_date", "id_number", "listening_score", "reading_score", "total_score", "exam_date", "image_file")
			resp, apiErr := c.Post("/student/language-certificate", token, body)
			return handleResp(resp, apiErr), nil
		},
	)
}
