import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import { registerAuthTools } from './tools/auth.js';
import { registerStudentTools } from './tools/student.js';
import { registerCourseTools } from './tools/course.js';
import { registerCtsvTools } from './tools/ctsv.js';
import { registerRoomTools } from './tools/room.js';

async function main() {
  const server = new McpServer({
    name: "uit-hub-mcp",
    version: "1.0.0"
  });

  // Đăng ký các modules tools
  registerAuthTools(server);
  registerStudentTools(server);
  registerCourseTools(server);
  registerCtsvTools(server);
  registerRoomTools(server);

  // Chạy server MCP với standard input/output
  const transport = new StdioServerTransport();
  await server.connect(transport);

  console.error('UIT Hub MCP Server đang chạy và lắng nghe qua stdio...');
}

main().catch((err) => {
  console.error("Lỗi khi khởi chạy MCP Server:", err);
  process.exit(1);
});
