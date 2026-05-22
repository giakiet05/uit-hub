import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { callFakeServer, formatMcpResponse } from '../utils/api.js';
import { PathCourseIdSchema, PathAssignmentIdSchema } from '../utils/types.js';

export function registerCourseTools(server: McpServer) {
  server.tool(
    "course_get_deadlines",
    "Lấy danh sách các hạn chót bài tập.",
    {
      token: z.string().describe("Bearer token")
    },
    async ({ token }) => {
      const result = await callFakeServer("course_get_deadlines", "GET", "/student/deadlines", token);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "course_get_materials",
    "Lấy tài liệu của một môn học.",
    {
      token: z.string().describe("Bearer token"),
      ...PathCourseIdSchema
    },
    async ({ token, courseId }) => {
      const result = await callFakeServer("course_get_materials", "GET", `/student/courses/${courseId}/materials`, token);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "course_get_assignments",
    "Lấy danh sách bài tập của môn học.",
    {
      token: z.string().describe("Bearer token"),
      ...PathCourseIdSchema
    },
    async ({ token, courseId }) => {
      const result = await callFakeServer("course_get_assignments", "GET", `/student/courses/${courseId}/assignments`, token);
      return formatMcpResponse(result);
    }
  );

  server.tool(
    "course_create_submission",
    "Nộp bài tập (text, link hoặc file).",
    {
      token: z.string().describe("Bearer token"),
      ...PathAssignmentIdSchema,
      submission_type: z.enum(["FILE", "LINK", "TEXT"]),
      text: z.string().optional(),
      url: z.string().url().optional(),
      file_urls: z.array(z.string()).optional(),
      comment: z.string().optional()
    },
    async (args) => {
      const { token, assignmentId, ...body } = args;
      const result = await callFakeServer("course_create_submission", "POST", `/student/assignments/${assignmentId}/submissions`, token, body);
      return formatMcpResponse(result);
    }
  );
}
