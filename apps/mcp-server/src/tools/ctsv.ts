import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { callFakeServer, formatMcpResponse } from '../utils/api.js';

export function registerCtsvTools(server: McpServer) {
  server.tool(
    "ctsv_create_confirm_letter",
    "Đăng ký giấy xác nhận sinh viên (nghĩa vụ quân sự, ký túc xá, v.v.).",
    {
      token: z.string().describe("Bearer token"),
      language: z.enum(["VI", "EN"]),
      reason: z.enum(["MILITARY_DEFERMENT", "DORM_EXTEND", "TAX_DEDUCTION_DOCS", "DEFENSE_EDU_REGISTRATION", "OTHER"]),
      other_reason: z.string().optional().describe("Chi tiết lý do khác (Bắt buộc khi reason là OTHER)"),
      request_type: z.enum(["NEW", "REISSUE"]),
      note: z.string().optional()
    },
    async (args) => {
      const { token, ...body } = args;
      if (body.reason === "OTHER" && !body.other_reason) {
        return formatMcpResponse({
          ok: false, tool: "ctsv_create_confirm_letter",
          error: { code: "BAD_REQUEST", message: "Cần cung cấp other_reason khi chọn reason là OTHER." }
        });
      }
      const result = await callFakeServer("ctsv_create_confirm_letter", "POST", "/student/confirm-letter", token, body);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "ctsv_create_bank_loans",
    "Đăng ký xác nhận vay vốn ngân hàng.",
    {
      token: z.string().describe("Bearer token"),
      benefit: z.enum(["NO_DISCOUNT", "TUITION_REDUCTION", "TUITION_EXEMPTION"]),
      orphan_status: z.enum(["NOT_ORPHAN", "ORPHAN"]),
      template: z.enum(["LEGACY", "STEM"]),
      note: z.string().optional()
    },
    async (args) => {
      const { token, ...body } = args;
      const result = await callFakeServer("ctsv_create_bank_loans", "POST", "/student/bank-loans", token, body);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "ctsv_create_training_point_confirm",
    "Đăng ký giấy xác nhận điểm rèn luyện.",
    {
      token: z.string().describe("Bearer token"),
      language: z.enum(["VI", "EN"]),
      note: z.string().optional()
    },
    async (args) => {
      const { token, ...body } = args;
      const result = await callFakeServer("ctsv_create_training_point_confirm", "POST", "/student/training-point-confirm", token, body);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "ctsv_create_language_certificate",
    "Nộp chứng chỉ ngoại ngữ.",
    {
      token: z.string().describe("Bearer token"),
      document_type: z.enum(["CONFIRMATION", "DIPLOMA", "CERTIFICATE"]),
      birth_date: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, "Định dạng YYYY-MM-DD"),
      id_number: z.string().min(6).describe("CCCD / CMND"),
      listening_score: z.number().min(0),
      reading_score: z.number().min(0),
      total_score: z.number().min(0),
      exam_date: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, "Định dạng YYYY-MM-DD"),
      image_file: z.string().describe("Base64 hoặc URL của ảnh chứng chỉ")
    },
    async (args) => {
      const { token, ...body } = args;
      const result = await callFakeServer("ctsv_create_language_certificate", "POST", "/student/language-certificate", token, body);
      return formatMcpResponse(result);
    }
  );
}
