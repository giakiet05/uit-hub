import { z } from 'zod';
import { callFakeServer, formatMcpResponse } from '../utils/api.js';
import { QueryYearSemesterSchema } from '../utils/types.js';
export function registerStudentTools(server) {
    // --- Query Tools ---
    server.tool("student_get_profile", "Lấy thông tin cá nhân của sinh viên.", {
        token: z.string().describe("Bearer token")
    }, async ({ token }) => {
        const result = await callFakeServer("student_get_profile", "GET", "/student/profile", token);
        return formatMcpResponse(result);
    });
    server.tool("student_get_schedule", "Xem lịch học học kỳ.", {
        token: z.string().describe("Bearer token"),
        ...QueryYearSemesterSchema
    }, async ({ token, year, semester }) => {
        const result = await callFakeServer("student_get_schedule", "GET", "/student/schedule", token, undefined, { year, semester });
        return formatMcpResponse(result);
    });
    server.tool("student_get_exam_schedule", "Xem lịch thi học kỳ.", {
        token: z.string().describe("Bearer token"),
        ...QueryYearSemesterSchema
    }, async ({ token, year, semester }) => {
        const result = await callFakeServer("student_get_exam_schedule", "GET", "/student/schedule/exam", token, undefined, { year, semester });
        return formatMcpResponse(result);
    });
    server.tool("student_get_score", "Xem bảng điểm.", {
        token: z.string().describe("Bearer token")
    }, async ({ token }) => {
        const result = await callFakeServer("student_get_score", "GET", "/student/score", token);
        return formatMcpResponse(result);
    });
    // --- Mutation Tools ---
    server.tool("student_create_contact", "Gửi form liên hệ.", {
        token: z.string().optional().describe("Bearer token (nếu có)"),
        name: z.string().min(1),
        email: z.string().email(),
        phone: z.string().optional(),
        subject: z.string().min(1),
        message: z.string().min(1)
    }, async (args) => {
        const { token, ...body } = args;
        const result = await callFakeServer("student_create_contact", "POST", "/contact", token, body);
        return formatMcpResponse(result);
    });
    // Todo: Thêm các tool đăng ký giấy tờ (transcript, graduate, parking, tuition_extend)
}
