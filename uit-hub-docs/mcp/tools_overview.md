# UIT Hub MCP Tools Overview



- Fake server stack: Go + Gin.
- Base path API: `/api/v1`.
- Response chuẩn:
  - Success: `{ "message": "Successfully!", "data": ... }`
  - Error: `{ "message": "...", "error_code": "..." }`

## Trạng thái triển khai backend hiện tại

Theo code trong `apps/fake-uit-server/internal/route`:

- Đã có: `POST /api/v1/login`, `GET /api/v1/student/profile`
- Demo/public cũ: `GET /api/v1/students`, `GET /api/v1/students/:id`
- Các endpoint khác trong `api-docs.md`: chưa triển khai đầy đủ ở fake server hiện tại

Vì vậy catalog tool sẽ có cả **implemented** và **planned**.

## Nguyên tắc thiết kế MCP tool

1. **1 tool chính = 1 endpoint API** để dễ trace/log/debug.
2. **Validation tại MCP** theo JSON schema từ `api-docs.md`.
3. **Không đổi shape response backend**; MCP chỉ chuẩn hóa envelope cho agent.
4. **Auth rõ ràng**: route private phải có Bearer token.
5. **Không suy diễn dữ liệu ngoài response** (trừ composite tool có quy tắc rõ ràng).

## Bộ tài liệu trong thư mục `mcp/`

- `tools_overview.md`: Tổng quan kiến trúc và nguyên tắc.
- `tool_classification.md`: Phân loại tool theo domain/risk/confirm policy.
- `tool_catalog.md`: Danh mục tool đầy đủ và mapping endpoint.
- `tool_definition_spec.md`: Chuẩn viết spec cho từng tool.
- `AI_invocation_rules.md`: Luật chọn/gọi tool cho agent.

## Quy ước đặt tên

- Prefix theo domain:
  - `auth_*`
  - `student_*`
  - `course_*`
  - `ctsv_*`
  - `room_*`
- Verb:
  - `get_*` cho read-only
  - `create_*` cho thao tác ghi/đăng ký

## Envelope output đề xuất cho MCP

```json
{
  "ok": true,
  "tool": "student_get_score",
  "endpoint": "GET /api/v1/student/score",
  "data": {}
}
```

Khi lỗi:

```json
{
  "ok": false,
  "tool": "student_get_score",
  "endpoint": "GET /api/v1/student/score",
  "error": {
    "http_status": 401,
    "code": "UNAUTHORIZED",
    "message": "Unauthorized!"
  }
}
```

## Phạm vi use case chính của chatbot

- Tra cứu học vụ: profile, lịch học, lịch thi, bảng điểm, học phí.
- Hỗ trợ học phần/course: deadlines, materials, assignments, submissions.
- Hỗ trợ CTSV: đăng ký các loại giấy xác nhận.
- Hỗ trợ phòng học/họp: tra cứu phòng trống.

`Đặt phòng họp` hiện là use case composite; API docs mới có availability, chưa có endpoint booking trực tiếp.