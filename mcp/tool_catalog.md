# MCP Tool Catalog (UIT Hub)

Catalog này là nguồn tham chiếu chính cho agent khi gọi MCP tools.

## 1. Chuẩn cột

- **Auth**: `yes` / `no` / `optional`
- **Status**:
  - `implemented`: đã có route trong `apps/fake-uit-server`
  - `planned`: có trong `api-docs.md` nhưng fake server chưa implement đủ

## 2. Core Auth Tools

| Tool | Method | Endpoint | Auth | Input chính | Status |
|---|---|---|---|---|---|
| `auth_login` | POST | `/api/v1/login` | no | `student_id`, `password`, `remember_me?` | implemented |

## 3. Student Tools

| Tool | Method | Endpoint | Auth | Input chính | Status |
|---|---|---|---|---|---|
| `student_get_profile` | GET | `/api/v1/student/profile` | yes | none | implemented |
| `student_get_schedule` | GET | `/api/v1/student/schedule` | yes | `year`, `semester` | planned |
| `student_get_exam_schedule` | GET | `/api/v1/student/schedule/exam` | yes | `year`, `semester` | planned |
| `student_get_score` | GET | `/api/v1/student/score` | yes | none | planned |
| `student_get_tuition_fee` | GET | `/api/v1/student/lookup/tuitionfee` | yes | none | planned |
| `student_get_insurance` | GET | `/api/v1/student/insurance` | yes | none | planned |
| `student_get_courses` | GET | `/api/v1/student/courses` | yes | none | planned |
| `student_get_training_points` | GET | `/api/v1/student/training-points` | yes | none | planned |
| `student_get_survey_form` | GET | `/api/v1/student/survey-form` | yes | none | planned |
| `student_create_transcript_registration` | POST | `/api/v1/student/transcript-regis` | yes | `copies`, `delivery_method`, ... | planned |
| `student_create_tuition_extend` | POST | `/api/v1/student/tuition-extend` | yes | `year`, `semester`, `requested_due_date`, `reason`, ... | planned |
| `student_create_monthly_parking` | POST | `/api/v1/student/monthly-parking` | yes | `vehicle_type`, `plate_number`, `months`, `start_month`, ... | planned |
| `student_create_graduate_registration` | POST | `/api/v1/student/graduate` | yes | `year`, `semester`, `email`, `phone`, ... | planned |
| `student_create_graduation_thesis_registration` | POST | `/api/v1/student/graduation-thesis` | yes | `thesis_title`, `advisor_name`, `team_members`, ... | planned |
| `student_create_contact` | POST | `/api/v1/contact` | optional | `name`, `email`, `subject`, `message`, `phone?` | planned |

## 4. Room Tools

| Tool | Method | Endpoint | Auth | Input chính | Status |
|---|---|---|---|---|---|
| `room_get_availability` | GET | `/api/v1/rooms/availability` | yes | `date`, `start`, `end` | planned |

## 5. Course Tools

| Tool | Method | Endpoint | Auth | Input chính | Status |
|---|---|---|---|---|---|
| `course_get_deadlines` | GET | `/api/v1/student/deadlines` | yes | none | planned |
| `course_get_materials` | GET | `/api/v1/student/courses/{courseId}/materials` | yes | `courseId` | planned |
| `course_get_assignments` | GET | `/api/v1/student/courses/{courseId}/assignments` | yes | `courseId` | planned |
| `course_create_submission` | POST | `/api/v1/student/assignments/{assignmentId}/submissions` | yes | `assignmentId`, `submission_type`, conditional fields | planned |

## 6. CTSV Tools

| Tool | Method | Endpoint | Auth | Input chính | Status |
|---|---|---|---|---|---|
| `ctsv_create_confirm_letter` | POST | `/api/v1/student/confirm-letter` | yes | `language`, `reason`, `request_type`, `other_reason?`, `note?` | planned |
| `ctsv_create_bank_loans` | POST | `/api/v1/student/bank-loans` | yes | `benefit`, `orphan_status`, `template`, `note?` | planned |
| `ctsv_create_training_point_confirm` | POST | `/api/v1/student/training-point-confirm` | yes | `language`, `note?` | planned |
| `ctsv_create_language_certificate` | POST | `/api/v1/student/language-certificate` | yes | `document_type`, `birth_date`, `id_number`, scores, `exam_date`, `image_file` | planned |

## 7. Composite Tools (đề xuất)

| Tool | Loại | Mục tiêu | Phụ thuộc |
|---|---|---|---|
| `student_get_remaining_credits` | composite query | Tính số tín chỉ còn thiếu | `student_get_score` + cấu hình số tín chỉ tốt nghiệp |
| `room_plan_meeting` | composite | Đề xuất phòng + khung giờ phù hợp | `room_get_availability` + rule chọn phòng |

## 8. Ghi chú triển khai

1. Tool catalog này bám API docs; khi fake server thêm endpoint, cập nhật `Status` từ `planned` -> `implemented`.
2. Mỗi tool cần có file spec chi tiết theo `mcp/tool_definition_spec.md`.
3. Các route demo cũ `/api/v1/students` và `/api/v1/students/:id` không nằm trong API docs chính nên không đưa vào catalog chính thức.