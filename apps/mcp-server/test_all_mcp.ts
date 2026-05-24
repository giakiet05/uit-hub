import { registerAuthTools } from './src/tools/auth.js';
import { registerStudentTools } from './src/tools/student.js';
import { registerRoomTools } from './src/tools/room.js';
import { registerCourseTools } from './src/tools/course.js';
import { registerCtsvTools } from './src/tools/ctsv.js';
import { config } from './src/config.js';

async function testAll() {
  console.log('Testing ALL MCP Tools to Fake-UIT-Server on', config.apiBaseUrl);

  const tools: Record<string, Function> = {};
  const mockServer = {
    tool: (name: string, desc: string, schema: any, handler: Function) => {
      tools[name] = handler;
    }
  } as any;

  registerAuthTools(mockServer);
  registerStudentTools(mockServer);
  registerRoomTools(mockServer);
  registerCourseTools(mockServer);
  registerCtsvTools(mockServer);

  let token = "";
  let successCount = 0;
  let totalCount = 0;

  const runTest = async (toolName: string, args: any) => {
    totalCount++;
    console.log(`\n--- Testing ${toolName} ---`);
    try {
      const res = await tools[toolName](args);
      const parsed = JSON.parse(res.content[0].text);
      if (parsed.ok) {
        console.log(`✅ OK: ${toolName}`);
        successCount++;
        return parsed.data;
      } else {
        console.log(`❌ ERROR: ${toolName}`, parsed.error);
        return null;
      }
    } catch (e: any) {
      console.log(`💥 EXCEPTION: ${toolName}`, e.message);
      return null;
    }
  };

  // 1. AUTH
  const loginData = await runTest('auth_login', { student_id: '22520001', password: 'pass123' });
  if (!loginData) return console.log('Login failed, aborting suite.');
  token = loginData.token;

  // 2. STUDENT QUERIES
  await runTest('student_get_profile', { token });
  await runTest('student_get_schedule', { token, year: 2024, semester: 1 });
  await runTest('student_get_exam_schedule', { token, year: 2024, semester: 1 });
  await runTest('student_get_score', { token });
  await runTest('student_get_tuition_fee', { token });
  await runTest('student_get_insurance', { token });
  await runTest('student_get_courses', { token });
  await runTest('student_get_training_points', { token });
  await runTest('student_get_survey_form', { token });
  await runTest('student_get_remaining_credits', { token, required_credits: 150 });

  // 3. STUDENT MUTATIONS
  await runTest('student_create_contact', { token, name: "Nguyễn Văn An", email: "22520001@gm.uit.edu.vn", subject: "Hỏi về học phí", message: "Tôi xin hỏi." });
  await runTest('student_create_transcript_registration', { token, copies: 2, language: "VI", delivery_method: "PICKUP" });
  await runTest('student_create_tuition_extend', { token, year: 2024, semester: 1, requested_due_date: "2024-12-31", reason: "Khó khăn" });
  await runTest('student_create_monthly_parking', { token, vehicle_type: "MOTORBIKE", plate_number: "59X1-12345", months: 3, start_month: "2024-05" });
  await runTest('student_create_graduate_registration', { token, year: 2024, semester: 1, email: "22520001@gm.uit.edu.vn", phone: "0901234001" });
  await runTest('student_create_graduation_thesis_registration', { token, thesis_title: "AI in Healthcare", advisor_name: "Dr. A", team_members: [{ student_id: "22520001" }] });

  // 4. ROOMS
  await runTest('room_get_availability', { token, date: "2024-05-20", start: "08:00", end: "10:00" });
  await runTest('room_plan_meeting', { token, date: "2024-05-20", start: "08:00", end: "10:00", capacity: 50 });

  // 5. COURSE
  await runTest('course_get_deadlines', { token });
  await runTest('course_get_materials', { token, courseId: "IT010" });
  await runTest('course_get_assignments', { token, courseId: "IT010" });
  await runTest('course_create_submission', { token, assignmentId: "ASG001", submission_type: "LINK", url: "https://github.com/abc" });

  // 6. CTSV
  await runTest('ctsv_create_confirm_letter', { token, language: "VI", reason: "DORM_EXTEND", request_type: "NEW" });
  await runTest('ctsv_create_bank_loans', { token, benefit: "NO_DISCOUNT", orphan_status: "NOT_ORPHAN", template: "LEGACY" });
  await runTest('ctsv_create_training_point_confirm', { token, language: "VI" });
  await runTest('ctsv_create_language_certificate', { token, document_type: "CERTIFICATE", birth_date: "2004-01-01", id_number: "079204000001", listening_score: 450, reading_score: 400, total_score: 850, exam_date: "2024-01-01", image_file: "aGVsbG8=" });

  console.log(`\n============== REPORT ==============`);
  console.log(`Total tests: ${totalCount}`);
  console.log(`Passed: ${successCount}`);
  console.log(`Failed: ${totalCount - successCount}`);
}

testAll();
