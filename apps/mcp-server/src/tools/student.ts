import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { callFakeServer, formatMcpResponse } from '../utils/api.js';
import { QueryYearSemesterSchema } from '../utils/types.js';

export function registerStudentTools(server: McpServer) {
  // --- Query Tools ---
  server.tool(
    "student_get_profile",
    "Lấy thông tin cá nhân của sinh viên.",
    { token: z.string().describe("Bearer token") },
    async ({ token }) => {
      const result = await callFakeServer("student_get_profile", "GET", "/student/profile", token);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "student_get_schedule",
    "Xem lịch học học kỳ.",
    { token: z.string().describe("Bearer token"), ...QueryYearSemesterSchema },
    async ({ token, year, semester }) => {
      const result = await callFakeServer("student_get_schedule", "GET", "/student/schedule", token, undefined, { year, semester });
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "student_get_exam_schedule",
    "Xem lịch thi học kỳ.",
    { token: z.string().describe("Bearer token"), ...QueryYearSemesterSchema },
    async ({ token, year, semester }) => {
      const result = await callFakeServer("student_get_exam_schedule", "GET", "/student/schedule/exam", token, undefined, { year, semester });
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "student_get_score",
    "Xem bảng điểm.",
    { token: z.string().describe("Bearer token") },
    async ({ token }) => {
      const result = await callFakeServer("student_get_score", "GET", "/student/score", token);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "student_get_tuition_fee",
    "Tra cứu học phí.",
    { token: z.string().describe("Bearer token") },
    async ({ token }) => {
      const result = await callFakeServer("student_get_tuition_fee", "GET", "/student/lookup/tuitionfee", token);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "student_get_insurance",
    "Tra cứu thông tin bảo hiểm y tế.",
    { token: z.string().describe("Bearer token") },
    async ({ token }) => {
      const result = await callFakeServer("student_get_insurance", "GET", "/student/insurance", token);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "student_get_courses",
    "Lấy danh sách các môn học.",
    { token: z.string().describe("Bearer token") },
    async ({ token }) => {
      const result = await callFakeServer("student_get_courses", "GET", "/student/courses", token);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "student_get_training_points",
    "Tra cứu điểm rèn luyện.",
    { token: z.string().describe("Bearer token") },
    async ({ token }) => {
      const result = await callFakeServer("student_get_training_points", "GET", "/student/training-points", token);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "student_get_survey_form",
    "Lấy danh sách form khảo sát sinh viên cần làm.",
    { token: z.string().describe("Bearer token") },
    async ({ token }) => {
      const result = await callFakeServer("student_get_survey_form", "GET", "/student/survey-form", token);
      return formatMcpResponse(result);
    }
  );

  // --- Composite Tools ---
  server.tool(
    "student_get_remaining_credits",
    "Tính số tín chỉ còn thiếu để tốt nghiệp.",
    {
      token: z.string().describe("Bearer token"),
      required_credits: z.number().min(0).describe("Tổng số tín chỉ yêu cầu của chương trình đào tạo")
    },
    async ({ token, required_credits }) => {
      const scoreResult = await callFakeServer("student_get_score", "GET", "/student/score", token);
      if (!scoreResult.ok) {
        return formatMcpResponse({
          ok: false,
          tool: "student_get_remaining_credits",
          error: {
            code: "NOT_FOUND",
            message: "Không đủ dữ liệu điểm để tính toán số tín chỉ tích lũy."
          }
        });
      }
      const earned_credits = scoreResult.data?.earned_credits || 0;
      const remaining = Math.max(required_credits - earned_credits, 0);

      return formatMcpResponse({
        ok: true,
        tool: "student_get_remaining_credits",
        data: {
          required_credits,
          earned_credits,
          remaining_credits: remaining
        }
      });
    }
  );

  // --- Mutation Tools ---
  server.tool(
    "student_create_contact",
    "Gửi form liên hệ.",
    {
      token: z.string().optional().describe("Bearer token (nếu có)"),
      name: z.string().min(1),
      email: z.string().email(),
      phone: z.string().optional(),
      subject: z.string().min(1),
      message: z.string().min(1)
    },
    async (args) => {
      const { token, ...body } = args;
      const result = await callFakeServer("student_create_contact", "POST", "/contact", token, body);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "student_create_transcript_registration",
    "Đăng ký in bảng điểm.",
    {
      token: z.string().describe("Bearer token"),
      copies: z.number().int().min(1).describe("Số lượng bản in"),
      language: z.enum(["VI", "EN"]).default("VI").describe("Ngôn ngữ"),
      delivery_method: z.enum(["PICKUP", "SHIP"]).describe("Phương thức nhận"),
      shipping_address: z.string().optional().describe("Địa chỉ nhận (bắt buộc nếu SHIP)"),
      phone: z.string().optional().describe("Số điện thoại (bắt buộc nếu SHIP)"),
      note: z.string().optional()
    },
    async (args) => {
      const { token, ...body } = args;
      if (body.delivery_method === "SHIP" && (!body.shipping_address || !body.phone)) {
        return formatMcpResponse({
          ok: false, tool: "student_create_transcript_registration",
          error: { code: "BAD_REQUEST", message: "Yêu cầu cung cấp địa chỉ và số điện thoại khi chọn SHIP." }
        });
      }
      const result = await callFakeServer("student_create_transcript_registration", "POST", "/student/transcript-regis", token, body);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "student_create_tuition_extend",
    "Xin gia hạn nộp học phí.",
    {
      token: z.string().describe("Bearer token"),
      year: z.number().int().min(1900),
      semester: z.number().int().min(1),
      requested_due_date: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, "Định dạng YYYY-MM-DD"),
      reason: z.string().min(1),
      phone: z.string().min(6).optional(),
      attachments: z.array(z.string()).default([])
    },
    async (args) => {
      const { token, ...body } = args;
      const result = await callFakeServer("student_create_tuition_extend", "POST", "/student/tuition-extend", token, body);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "student_create_monthly_parking",
    "Đăng ký vé xe tháng.",
    {
      token: z.string().describe("Bearer token"),
      vehicle_type: z.enum(["MOTORBIKE", "CAR", "BICYCLE"]),
      plate_number: z.string().min(3),
      months: z.number().int().min(1).max(12),
      start_month: z.string().regex(/^\d{4}-\d{2}$/, "Định dạng YYYY-MM"),
      owner_name: z.string().optional(),
      note: z.string().optional()
    },
    async (args) => {
      const { token, ...body } = args;
      const result = await callFakeServer("student_create_monthly_parking", "POST", "/student/monthly-parking", token, body);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "student_create_graduate_registration",
    "Đăng ký xét tốt nghiệp.",
    {
      token: z.string().describe("Bearer token"),
      year: z.number().int().min(1900),
      semester: z.number().int().min(1),
      email: z.string().email(),
      phone: z.string().min(6),
      address: z.string().optional(),
      note: z.string().optional()
    },
    async (args) => {
      const { token, ...body } = args;
      const result = await callFakeServer("student_create_graduate_registration", "POST", "/student/graduate", token, body);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "student_create_graduation_thesis_registration",
    "Đăng ký khóa luận tốt nghiệp.",
    {
      token: z.string().describe("Bearer token"),
      thesis_title: z.string().min(1),
      advisor_name: z.string().min(1),
      advisor_email: z.string().email().optional(),
      team_members: z.array(z.object({
        student_id: z.string().min(1),
        full_name: z.string().optional()
      })).min(1),
      note: z.string().optional()
    },
    async (args) => {
      const { token, ...body } = args;
      const result = await callFakeServer("student_create_graduation_thesis_registration", "POST", "/student/graduation-thesis", token, body);
      return formatMcpResponse(result);
    }
  );
}
