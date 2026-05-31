# MCP Tool Classification

## 1. Phân loại theo domain

| Domain | Prefix | Mô tả |
|---|---|---|
| Authentication | `auth_` | Đăng nhập/lấy token |
| Student | `student_` | Tra cứu và đăng ký tác vụ học vụ |
| Room | `room_` | Tra cứu phòng trống |
| Course | `course_` | Deadline, tài liệu, bài tập, nộp bài |
| CTSV | `ctsv_` | Các dịch vụ giấy xác nhận |

## 2. Phân loại theo loại thao tác

| Type | Đặc điểm | Ví dụ |
|---|---|---|
| Query | Read-only, không thay đổi dữ liệu | `student_get_score` |
| Mutation | Có tạo/yêu cầu nghiệp vụ | `student_create_transcript_registration` |
| Composite | Gọi nhiều tool để tạo câu trả lời | `student_get_remaining_credits` (suy luận) |

## 3. Phân loại theo rủi ro và xác nhận

| Risk | Confirmation policy | Ví dụ |
|---|---|---|
| Low | none | Các GET tra cứu |
| Medium | soft | `student_create_contact` |
| High | strict | Nộp bài, đăng ký hồ sơ CTSV, đăng ký tốt nghiệp |

## 4. Bảng phân loại toàn bộ tool theo `api-docs.md`

| Tool | Domain | Type | Risk | Auth | Endpoint | Status |
|---|---|---|---|---|---|---|
| `auth_login` | auth | mutation | low | no | `POST /api/v1/login` | implemented |
| `student_get_profile` | student | query | low | yes | `GET /api/v1/student/profile` | implemented |
| `student_get_schedule` | student | query | low | yes | `GET /api/v1/student/schedule` | planned |
| `student_get_exam_schedule` | student | query | low | yes | `GET /api/v1/student/schedule/exam` | planned |
| `student_get_score` | student | query | low | yes | `GET /api/v1/student/score` | planned |
| `student_get_tuition_fee` | student | query | low | yes | `GET /api/v1/student/lookup/tuitionfee` | planned |
| `student_get_insurance` | student | query | low | yes | `GET /api/v1/student/insurance` | planned |
| `student_get_courses` | student | query | low | yes | `GET /api/v1/student/courses` | planned |
| `student_get_training_points` | student | query | low | yes | `GET /api/v1/student/training-points` | planned |
| `student_get_survey_form` | student | query | low | yes | `GET /api/v1/student/survey-form` | planned |
| `student_create_transcript_registration` | student | mutation | high | yes | `POST /api/v1/student/transcript-regis` | planned |
| `student_create_tuition_extend` | student | mutation | high | yes | `POST /api/v1/student/tuition-extend` | planned |
| `student_create_monthly_parking` | student | mutation | medium | yes | `POST /api/v1/student/monthly-parking` | planned |
| `student_create_graduate_registration` | student | mutation | high | yes | `POST /api/v1/student/graduate` | planned |
| `student_create_graduation_thesis_registration` | student | mutation | high | yes | `POST /api/v1/student/graduation-thesis` | planned |
| `room_get_availability` | room | query | low | yes | `GET /api/v1/rooms/availability` | planned |
| `student_create_contact` | student | mutation | medium | optional | `POST /api/v1/contact` | planned |
| `course_get_deadlines` | course | query | low | yes | `GET /api/v1/student/deadlines` | planned |
| `course_get_materials` | course | query | low | yes | `GET /api/v1/student/courses/{courseId}/materials` | planned |
| `course_get_assignments` | course | query | low | yes | `GET /api/v1/student/courses/{courseId}/assignments` | planned |
| `course_create_submission` | course | mutation | high | yes | `POST /api/v1/student/assignments/{assignmentId}/submissions` | planned |
| `ctsv_create_confirm_letter` | ctsv | mutation | high | yes | `POST /api/v1/student/confirm-letter` | planned |
| `ctsv_create_bank_loans` | ctsv | mutation | high | yes | `POST /api/v1/student/bank-loans` | planned |
| `ctsv_create_training_point_confirm` | ctsv | mutation | high | yes | `POST /api/v1/student/training-point-confirm` | planned |
| `ctsv_create_language_certificate` | ctsv | mutation | high | yes | `POST /api/v1/student/language-certificate` | planned |

## 5. Composite tools đề xuất (không map 1:1 endpoint)

| Tool | Dùng khi nào | Thành phần |
|---|---|---|
| `student_get_remaining_credits` | Hỏi còn bao nhiêu tín chỉ để tốt nghiệp | `student_get_score` + cấu hình chương trình |
| `room_plan_meeting` | Hỏi đặt phòng họp | `room_get_availability` + policy chọn phòng |

Lưu ý: composite tool phải ghi rõ công thức suy luận, không suy diễn mơ hồ.