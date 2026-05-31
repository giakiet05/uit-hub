# AI Invocation Rules for UIT MCP Tools

Luật này quy định cách agent chọn và gọi MCP tools khi hội thoại với sinh viên UIT.

## 1. Mục tiêu

- Gọi đúng tool, đúng tham số, đúng thời điểm.
- Không gọi tool khi thiếu dữ liệu bắt buộc.
- Trả lời minh bạch khi endpoint chưa triển khai.

## 2. Quy trình gọi tool chuẩn

1. **Intent detection**: xác định user muốn tra cứu hay đăng ký.
2. **Tool selection**: map intent -> tool trong `mcp/tool_catalog.md`.
3. **Slot filling**: hỏi bù tham số bắt buộc còn thiếu.
4. **Auth check**: nếu tool cần auth mà thiếu token -> yêu cầu đăng nhập.
5. **Invoke**: gọi tool với payload đã validate.
6. **Response rendering**: diễn giải ngắn gọn cho user.

## 3. Mapping intent nhanh

| User message mẫu | Tool |
|---|---|
| “Đăng nhập cho mình” | `auth_login` |
| “Xem profile” | `student_get_profile` |
| “Xem lịch học học kỳ 1 năm 2026” | `student_get_schedule` |
| “Xem lịch thi” | `student_get_exam_schedule` |
| “Tra cứu bảng điểm” | `student_get_score` |
| “Còn bao nhiêu tín chỉ để tốt nghiệp” | `student_get_remaining_credits` (composite) |
| “Phòng nào trống 9h-11h ngày mai” | `room_get_availability` |
| “Nộp bài assignment A1 bằng link” | `course_create_submission` |
| “Đăng ký giấy xác nhận sinh viên” | `ctsv_create_confirm_letter` |

## 4. Slot filling rules

Agent phải hỏi tiếp khi thiếu:

- `student_get_schedule`, `student_get_exam_schedule`: thiếu `year` hoặc `semester`.
- `room_get_availability`: thiếu `date`, `start`, `end`.
- `course_create_submission`:
  - `TEXT` thiếu `text`
  - `LINK` thiếu `url`
  - `FILE` thiếu `file_urls`
- `student_create_transcript_registration`:
  - nếu `delivery_method=SHIP` thì bắt buộc `shipping_address`, `phone`.
- `ctsv_create_confirm_letter`:
  - nếu `reason=OTHER` bắt buộc `other_reason` hợp lệ.

## 5. Auth rules

- Tool có `Auth=yes` chỉ gọi khi có token.
- Nếu chưa có token:
  1. Gợi ý gọi `auth_login`.
  2. Không gọi backend private endpoint.
- Khi backend trả `UNAUTHORIZED`, agent yêu cầu người dùng đăng nhập lại.

## 6. Confirmation policy

- **none** (không cần xác nhận lại): tất cả query GET.
- **soft**: `student_create_contact`.
- **strict**: mọi mutation đăng ký hồ sơ/đơn từ hoặc nộp bài.

Với `strict`, agent phải tóm tắt payload và hỏi xác nhận trước khi gọi.

## 7. Error handling policy

| error_code | Hành vi agent |
|---|---|
| `BAD_REQUEST` | Báo field sai/thiếu và yêu cầu nhập lại |
| `UNAUTHORIZED` | Yêu cầu đăng nhập lại |
| `FORBIDDEN` | Báo user không đủ quyền |
| `NOT_FOUND` | Báo không tìm thấy dữ liệu tương ứng |
| `TOO_MANY_REQUESTS` | Đề nghị thử lại sau |
| `SERVICE_UNAVAILABLE` | Báo hệ thống bận, cho retry |
| `INTERNAL_ERROR` | Báo lỗi hệ thống và đề xuất thử lại |

## 8. Retry and timeout

- Query GET: timeout 5-10s, retry 1 lần khi lỗi mạng tạm thời.
- Mutation POST: timeout 10-15s, không auto-retry (tránh gửi trùng).
- Nếu gặp timeout: trả lời rõ “hệ thống phản hồi chậm”, không bịa kết quả.

## 9. Rules khi endpoint chưa triển khai

Nếu tool có status `planned`:

1. Không giả vờ có dữ liệu thật.
2. Thông báo endpoint đang trong giai đoạn phát triển fake server.
3. Nếu có tool gần đúng (ví dụ profile/score) thì đề xuất thay thế.

## 10. Composite rule: tín chỉ còn lại để tốt nghiệp

Tool: `student_get_remaining_credits`

Quy tắc:

1. Lấy dữ liệu tích lũy từ `student_get_score` (nếu available).
2. Lấy `required_credits` từ cấu hình chương trình (không hardcode mơ hồ trong lời đáp).
3. Tính: `remaining = max(required_credits - earned_credits, 0)`.
4. Nếu thiếu một trong hai nguồn, trả về trạng thái “không đủ dữ liệu”.

## 11. Safety and privacy

- Không in token trong trả lời hoặc log.
- Không in `image_file` base64 đầy đủ.
- Không tự ý lưu thông tin nhạy cảm của sinh viên ngoài ngữ cảnh phiên.

## 12. Output style cho câu trả lời cuối

- Trả lời ngắn gọn, trực tiếp.
- Nêu kết quả trước, chi tiết sau.
- Nếu lỗi, nêu nguyên nhân và hành động cần làm tiếp theo.