package tools

import (
	"context"

	"github.com/giakiet05/uit-hub/apps/mcp-server/internal/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterCTSV registers CTSV (student affairs) tools.
func RegisterCTSV(s *server.MCPServer, c *client.Client) {
	tokenDesc := "Bearer token nhận từ auth_login"

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
