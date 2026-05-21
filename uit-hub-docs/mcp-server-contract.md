# MCP Server Contract

File này ghi các quy ước tối thiểu cho team implement MCP server để agent có thể dùng tool ổn định.

## Mục tiêu

MCP server là nơi expose các API/hành động của hệ thống trường thành tool cho agent.

Agent sẽ:

- kết nối MCP server
- gọi `tools/list` để lấy tool definitions
- passthrough `inputSchema` gần như nguyên bản cho LLM
- nhận tool call từ LLM
- gửi arguments sang MCP server bằng `tools/call`

Vì vậy schema của MCP tool là contract quan trọng giữa MCP server, agent và LLM.

## Tool Name

Tool name phải:

- dùng `snake_case`
- ngắn nhưng rõ nghĩa
- ổn định, không đổi tùy tiện
- không tự nhét prefix trường/server vào name, agent sẽ tự namespace

Ví dụ tốt:

```text
get_student_profile
search_courses
submit_leave_request
upsert_student_note
delete_student_note
```

## Description

Description phải nói rõ:

- tool làm gì
- khi nào nên dùng
- tool là read-only hay write operation nếu có side effect

Ví dụ:

```text
Get a student's academic profile by student ID. Read-only.
Create or update an advisor note for a student. Write operation.
```

Tránh description quá dài hoặc nhồi quá nhiều chi tiết backend không cần thiết.

## Input Schema

`inputSchema` phải là JSON Schema object.

Quy ước tối thiểu:

- top-level luôn là `type: object`
- khai báo `properties` rõ ràng
- mỗi property phải có `description`
- field bắt buộc phải nằm trong `required`
- dùng `additionalProperties: false` nếu có thể
- dùng `enum` nếu field chỉ nhận một tập giá trị cố định
- array phải khai báo `items`
- object nested phải khai báo `properties`

Ví dụ:

```json
{
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "student_id": {
      "type": "string",
      "description": "Student ID, for example 22520001."
    },
    "visibility": {
      "type": "string",
      "description": "Who can see this note.",
      "enum": ["private", "advisor", "public"]
    }
  },
  "required": ["student_id"]
}
```

Schema càng rõ thì LLM càng dễ gọi đúng tool.

## Validation

MCP server phải tự validate input thật. Không được tin LLM hoặc agent.

Server cần xử lý:

- thiếu required field
- type sai
- enum sai
- ID không tồn tại
- input rỗng hoặc mơ hồ
- write operation không đủ điều kiện thực hiện

Với lỗi nghiệp vụ, trả tool error ngắn gọn và an toàn. Không dump stack trace hoặc chi tiết nội bộ.

## Output

Output nên là structured JSON ổn định.

Quy ước:

- field name dùng `snake_case`
- read không có dữ liệu thì trả `found: false` nếu hợp lý
- write thành công thì trả `id`, `status`, hoặc cờ `created`, `updated`, `deleted`
- lỗi nghiệp vụ trả tool error rõ ràng

Ví dụ output tốt:

```json
{
  "found": true,
  "profile": {
    "student_id": "22520001",
    "full_name": "Nguyen Van An",
    "gpa": 3.42,
    "academic_warning": false
  }
}
```

## Side Effects

Read tool và write tool phải phân biệt rõ trong description.

Write tool nên:

- idempotent nếu có thể
- validate input chặt
- không tự ý tạo side effect nếu input mơ hồ
- trả trạng thái rõ ràng sau khi chạy

Ví dụ:

```json
{
  "created": true,
  "note_id": "note-001"
}
```

## Security

MCP server không được:

- trả secret, token, password trong output
- log raw sensitive input
- trả stack trace cho agent
- expose tool nguy hiểm khi chưa có permission model

MCP server nên:

- timeout hợp lý
- giới hạn kích thước output
- sanitize lỗi trước khi trả về agent
- phân loại rõ tool nào read-only, tool nào có side effect

## Nguyên tắc chốt

MCP tool schema là contract cho LLM gọi tool. Agent sẽ passthrough `inputSchema` gần như nguyên bản cho model, nên schema càng rõ thì model càng gọi đúng.

Nhưng MCP server vẫn phải tự validate mọi input, vì LLM có thể gọi sai.
