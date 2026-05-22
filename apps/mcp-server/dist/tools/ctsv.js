import { z } from 'zod';
import { callFakeServer, formatMcpResponse } from '../utils/api.js';
export function registerCtsvTools(server) {
    server.tool("ctsv_create_confirm_letter", "Đăng ký giấy xác nhận sinh viên (nghĩa vụ quân sự, ký túc xá, v.v.).", {
        token: z.string().describe("Bearer token"),
        language: z.enum(["VI", "EN"]),
        reason: z.enum(["MILITARY_DEFERMENT", "DORM_EXTEND", "TAX_DEDUCTION_DOCS", "DEFENSE_EDU_REGISTRATION", "OTHER"]),
        other_reason: z.string().optional(),
        request_type: z.enum(["NEW", "REISSUE"]),
        note: z.string().optional()
    }, async (args) => {
        const { token, ...body } = args;
        const result = await callFakeServer("ctsv_create_confirm_letter", "POST", "/student/confirm-letter", token, body);
        return formatMcpResponse(result);
    });
    // Todo: Thêm các tool bank_loans, training_point_confirm, language_certificate
}
