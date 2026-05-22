import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { callFakeServer, formatMcpResponse } from '../utils/api.js';

export function registerAuthTools(server: McpServer) {
  server.tool(
    "auth_login",
    "Đăng nhập lấy token. Gọi khi user yêu cầu đăng nhập hoặc thiếu token.",
    {
      student_id: z.string().describe("Mã số sinh viên"),
      password: z.string().min(1).describe("Mật khẩu"),
      remember_me: z.boolean().default(false).optional()
    },
    async (args) => {
      const result = await callFakeServer(
        "auth_login",
        "POST",
        "/login",
        undefined,
        args
      );
      return formatMcpResponse(result);
    }
  );
}
