import { z } from 'zod';

export const QueryYearSemesterSchema = {
  year: z.number().int().min(1900).describe("Năm học (Ví dụ: 2026)"),
  semester: z.number().int().min(1).describe("Học kỳ (1, 2, 3)")
};

export const QueryRoomsAvailabilitySchema = {
  date: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, "Định dạng YYYY-MM-DD").describe("Ngày muốn tra cứu phòng (VD: 2026-10-15)"),
  start: z.string().describe("Giờ bắt đầu (VD: 07:30)"),
  end: z.string().describe("Giờ kết thúc (VD: 11:30)")
};

export const PathCourseIdSchema = {
  courseId: z.string().describe("ID của môn học")
};

export const PathAssignmentIdSchema = {
  assignmentId: z.string().describe("ID của bài tập (assignment)")
};
