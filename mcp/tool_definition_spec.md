# MCP Tool Definition Spec (Completed)

Tài liệu này là **đặc tả hoàn chỉnh** cho toàn bộ MCP tools của UIT chatbot, bám theo:

- `uit-hub-docs/api-docs.md`
- `mcp/tool_catalog.md`
- trạng thái route thực tế trong `apps/fake-uit-server`

---

## 1. Global Runtime Contract

### 1.1 Base API

- Base URL: `http://localhost:<PORT>`
- Base path: `/api/v1`

### 1.2 MCP Success Envelope

```json
{
  "ok": true,
  "tool": "student_get_score",
  "endpoint": "GET /api/v1/student/score",
  "data": {},
  "meta": {
    "source": "fake-uit-server"
  }
}
```

### 1.3 MCP Error Envelope

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

### 1.4 Error Mapping

| HTTP | error_code |
|---|---|
| 400 | BAD_REQUEST |
| 401 | UNAUTHORIZED |
| 403 | FORBIDDEN |
| 404 | NOT_FOUND |
| 429 | TOO_MANY_REQUESTS |
| 503 | SERVICE_UNAVAILABLE |
| Other | INTERNAL_ERROR |

### 1.5 Timeout & Retry

- Query GET: timeout 5-10s, retry 1 lần khi lỗi mạng tạm thời.
- Mutation POST: timeout 10-15s, không auto-retry.

---

## 2. Shared Schema References

### 2.1 `QueryYearSemester`

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["year", "semester"],
  "properties": {
    "year": { "type": "integer", "minimum": 1900 },
    "semester": { "type": "integer", "minimum": 1 }
  }
}
```

### 2.2 `QueryRoomsAvailability`

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["date", "start", "end"],
  "properties": {
    "date": { "type": "string", "pattern": "^\\d{4}-\\d{2}-\\d{2}$" },
    "start": { "type": "string" },
    "end": { "type": "string" }
  }
}
```

### 2.3 `PathCourseId`

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["courseId"],
  "properties": { "courseId": { "type": "string" } }
}
```

### 2.4 `PathAssignmentId`

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["assignmentId"],
  "properties": { "assignmentId": { "type": "string" } }
}
```

---

## 3. Tool Definitions

## 3.1 Auth

### `auth_login`

```yaml
domain: auth
type: mutation
risk_level: low
auth_required: false
confirmation_policy: none
method: POST
path: /api/v1/login
status: implemented
```

**Purpose**: Đăng nhập lấy token.

**Input schema**:

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["student_id", "password"],
  "properties": {
    "student_id": { "type": "string" },
    "password": { "type": "string", "minLength": 1 },
    "remember_me": { "type": "boolean", "default": false }
  }
}
```

**Invocation rules**:

- Chỉ gọi khi user yêu cầu đăng nhập hoặc thiếu token cho tool private.
- Không log password.

---

## 3.2 Student Query Tools

### `student_get_profile`

```yaml
type: query
risk_level: low
auth_required: true
confirmation_policy: none
method: GET
path: /api/v1/student/profile
status: implemented
```

**Input schema**: `{}` (không tham số).

### `student_get_schedule`

```yaml
type: query
risk_level: low
auth_required: true
method: GET
path: /api/v1/student/schedule
status: planned
```

**Input schema**: `QueryYearSemester`.

**Invocation rule**: thiếu `year`/`semester` thì hỏi bù trước khi gọi.

### `student_get_exam_schedule`

```yaml
type: query
risk_level: low
auth_required: true
method: GET
path: /api/v1/student/schedule/exam
status: planned
```

**Input schema**: `QueryYearSemester`.

### `student_get_score`

```yaml
type: query
risk_level: low
auth_required: true
method: GET
path: /api/v1/student/score
status: planned
```

**Input schema**: `{}`.

### `student_get_tuition_fee`

```yaml
type: query
risk_level: low
auth_required: true
method: GET
path: /api/v1/student/lookup/tuitionfee
status: planned
```

**Input schema**: `{}`.

### `student_get_insurance`

```yaml
type: query
risk_level: low
auth_required: true
method: GET
path: /api/v1/student/insurance
status: planned
```

**Input schema**: `{}`.

### `student_get_courses`

```yaml
type: query
risk_level: low
auth_required: true
method: GET
path: /api/v1/student/courses
status: planned
```

**Input schema**: `{}`.

### `student_get_training_points`

```yaml
type: query
risk_level: low
auth_required: true
method: GET
path: /api/v1/student/training-points
status: planned
```

**Input schema**: `{}`.

### `student_get_survey_form`

```yaml
type: query
risk_level: low
auth_required: true
method: GET
path: /api/v1/student/survey-form
status: planned
```

**Input schema**: `{}`.

---

## 3.3 Student Mutation Tools

### `student_create_transcript_registration`

```yaml
type: mutation
risk_level: high
auth_required: true
confirmation_policy: strict
method: POST
path: /api/v1/student/transcript-regis
status: planned
```

**Input schema**:

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["copies", "delivery_method"],
  "properties": {
    "copies": { "type": "integer", "minimum": 1 },
    "language": { "type": "string", "enum": ["VI", "EN"], "default": "VI" },
    "delivery_method": { "type": "string", "enum": ["PICKUP", "SHIP"] },
    "shipping_address": { "type": "string", "minLength": 1 },
    "phone": { "type": "string", "minLength": 6 },
    "note": { "type": "string" }
  },
  "allOf": [
    {
      "if": { "properties": { "delivery_method": { "const": "SHIP" } }, "required": ["delivery_method"] },
      "then": { "required": ["shipping_address", "phone"] }
    }
  ]
}
```

### `student_create_tuition_extend`

```yaml
type: mutation
risk_level: high
auth_required: true
confirmation_policy: strict
method: POST
path: /api/v1/student/tuition-extend
status: planned
```

**Input schema**:

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["year", "semester", "requested_due_date", "reason"],
  "properties": {
    "year": { "type": "integer", "minimum": 1900 },
    "semester": { "type": "integer", "minimum": 1 },
    "requested_due_date": { "type": "string", "pattern": "^\\d{4}-\\d{2}-\\d{2}$" },
    "reason": { "type": "string", "minLength": 1 },
    "phone": { "type": "string", "minLength": 6 },
    "attachments": { "type": "array", "items": { "type": "string" }, "default": [] }
  }
}
```

### `student_create_monthly_parking`

```yaml
type: mutation
risk_level: medium
auth_required: true
confirmation_policy: soft
method: POST
path: /api/v1/student/monthly-parking
status: planned
```

**Input schema**:

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["vehicle_type", "plate_number", "months", "start_month"],
  "properties": {
    "vehicle_type": { "type": "string", "enum": ["MOTORBIKE", "CAR", "BICYCLE"] },
    "plate_number": { "type": "string", "minLength": 3 },
    "months": { "type": "integer", "minimum": 1, "maximum": 12 },
    "start_month": { "type": "string", "pattern": "^\\d{4}-\\d{2}$" },
    "owner_name": { "type": "string" },
    "note": { "type": "string" }
  }
}
```

### `student_create_graduate_registration`

```yaml
type: mutation
risk_level: high
auth_required: true
confirmation_policy: strict
method: POST
path: /api/v1/student/graduate
status: planned
```

**Input schema**:

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["year", "semester", "email", "phone"],
  "properties": {
    "year": { "type": "integer", "minimum": 1900 },
    "semester": { "type": "integer", "minimum": 1 },
    "email": { "type": "string", "format": "email" },
    "phone": { "type": "string", "minLength": 6 },
    "address": { "type": "string" },
    "note": { "type": "string" }
  }
}
```

### `student_create_graduation_thesis_registration`

```yaml
type: mutation
risk_level: high
auth_required: true
confirmation_policy: strict
method: POST
path: /api/v1/student/graduation-thesis
status: planned
```

**Input schema**:

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["thesis_title", "advisor_name", "team_members"],
  "properties": {
    "thesis_title": { "type": "string", "minLength": 1 },
    "advisor_name": { "type": "string", "minLength": 1 },
    "advisor_email": { "type": "string", "format": "email" },
    "team_members": {
      "type": "array",
      "minItems": 1,
      "items": {
        "type": "object",
        "additionalProperties": false,
        "required": ["student_id"],
        "properties": {
          "student_id": { "type": "string", "minLength": 1 },
          "full_name": { "type": "string" }
        }
      }
    },
    "note": { "type": "string" }
  }
}
```

### `student_create_contact`

```yaml
type: mutation
risk_level: medium
auth_required: optional
confirmation_policy: soft
method: POST
path: /api/v1/contact
status: planned
```

**Input schema**:

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["name", "email", "subject", "message"],
  "properties": {
    "name": { "type": "string", "minLength": 1 },
    "email": { "type": "string", "format": "email" },
    "phone": { "type": "string" },
    "subject": { "type": "string", "minLength": 1 },
    "message": { "type": "string", "minLength": 1 }
  }
}
```

---

## 3.4 Room Tools

### `room_get_availability`

```yaml
type: query
risk_level: low
auth_required: true
confirmation_policy: none
method: GET
path: /api/v1/rooms/availability
status: planned
```

**Input schema**: `QueryRoomsAvailability`.

---

## 3.5 Course Tools

### `course_get_deadlines`

```yaml
type: query
risk_level: low
auth_required: true
method: GET
path: /api/v1/student/deadlines
status: planned
```

**Input schema**: `{}`.

### `course_get_materials`

```yaml
type: query
risk_level: low
auth_required: true
method: GET
path: /api/v1/student/courses/{courseId}/materials
status: planned
```

**Input schema**: `PathCourseId`.

### `course_get_assignments`

```yaml
type: query
risk_level: low
auth_required: true
method: GET
path: /api/v1/student/courses/{courseId}/assignments
status: planned
```

**Input schema**: `PathCourseId`.

### `course_create_submission`

```yaml
type: mutation
risk_level: high
auth_required: true
confirmation_policy: strict
method: POST
path: /api/v1/student/assignments/{assignmentId}/submissions
status: planned
```

**Input schema**:

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["assignmentId", "submission_type"],
  "properties": {
    "assignmentId": { "type": "string" },
    "submission_type": { "type": "string", "enum": ["FILE", "LINK", "TEXT"] },
    "text": { "type": "string" },
    "url": { "type": "string", "format": "uri" },
    "file_urls": { "type": "array", "items": { "type": "string" }, "default": [] },
    "comment": { "type": "string" }
  },
  "allOf": [
    {
      "if": { "properties": { "submission_type": { "const": "TEXT" } }, "required": ["submission_type"] },
      "then": { "required": ["text"] }
    },
    {
      "if": { "properties": { "submission_type": { "const": "LINK" } }, "required": ["submission_type"] },
      "then": { "required": ["url"] }
    },
    {
      "if": { "properties": { "submission_type": { "const": "FILE" } }, "required": ["submission_type"] },
      "then": { "required": ["file_urls"] }
    }
  ]
}
```

---

## 3.6 CTSV Tools

### `ctsv_create_confirm_letter`

```yaml
type: mutation
risk_level: high
auth_required: true
confirmation_policy: strict
method: POST
path: /api/v1/student/confirm-letter
status: planned
```

**Input schema**:

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["language", "reason", "request_type"],
  "properties": {
    "language": { "type": "string", "enum": ["VI", "EN"] },
    "reason": {
      "type": "string",
      "enum": ["MILITARY_DEFERMENT", "DORM_EXTEND", "TAX_DEDUCTION_DOCS", "DEFENSE_EDU_REGISTRATION", "OTHER"]
    },
    "other_reason": { "type": "string", "minLength": 1, "pattern": "^Bổ sung hồ sơ\\b.*$" },
    "request_type": { "type": "string", "enum": ["NEW", "REISSUE"] },
    "note": { "type": "string" }
  },
  "allOf": [
    {
      "if": { "properties": { "reason": { "const": "OTHER" } }, "required": ["reason"] },
      "then": { "required": ["other_reason"] },
      "else": { "not": { "required": ["other_reason"] } }
    }
  ]
}
```

### `ctsv_create_bank_loans`

```yaml
type: mutation
risk_level: high
auth_required: true
confirmation_policy: strict
method: POST
path: /api/v1/student/bank-loans
status: planned
```

**Input schema**:

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["benefit", "orphan_status", "template"],
  "properties": {
    "benefit": { "type": "string", "enum": ["NO_DISCOUNT", "TUITION_REDUCTION", "TUITION_EXEMPTION"] },
    "orphan_status": { "type": "string", "enum": ["NOT_ORPHAN", "ORPHAN"] },
    "template": { "type": "string", "enum": ["LEGACY", "STEM"] },
    "note": { "type": "string" }
  }
}
```

### `ctsv_create_training_point_confirm`

```yaml
type: mutation
risk_level: high
auth_required: true
confirmation_policy: strict
method: POST
path: /api/v1/student/training-point-confirm
status: planned
```

**Input schema**:

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["language"],
  "properties": {
    "language": { "type": "string", "enum": ["VI", "EN"] },
    "note": { "type": "string" }
  }
}
```

### `ctsv_create_language_certificate`

```yaml
type: mutation
risk_level: high
auth_required: true
confirmation_policy: strict
method: POST
path: /api/v1/student/language-certificate
status: planned
```

**Input schema**:

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["document_type", "birth_date", "id_number", "listening_score", "reading_score", "total_score", "exam_date", "image_file"],
  "properties": {
    "document_type": { "type": "string", "enum": ["CONFIRMATION", "DIPLOMA", "CERTIFICATE"] },
    "birth_date": { "type": "string", "pattern": "^\\d{4}-\\d{2}-\\d{2}$" },
    "id_number": { "type": "string", "minLength": 6 },
    "listening_score": { "type": "number", "minimum": 0 },
    "reading_score": { "type": "number", "minimum": 0 },
    "total_score": { "type": "number", "minimum": 0 },
    "exam_date": { "type": "string", "pattern": "^\\d{4}-\\d{2}-\\d{2}$" },
    "image_file": { "type": "string" }
  }
}
```

---

## 3.7 Composite Tools

### `student_get_remaining_credits`

```yaml
type: composite
risk_level: low
auth_required: true
confirmation_policy: none
status: planned
depends_on:
  - student_get_score
```

**Purpose**: Trả số tín chỉ còn thiếu để tốt nghiệp.

**Input schema**:

```json
{
  "type": "object",
  "additionalProperties": false,
  "required": ["required_credits"],
  "properties": {
    "required_credits": { "type": "number", "minimum": 0 }
  }
}
```

**Computation rule**:

`remaining = max(required_credits - earned_credits, 0)`

Nếu không lấy được `earned_credits` từ nguồn dữ liệu thì trả `ok=false` với `code=NOT_FOUND` và message rõ thiếu dữ liệu.

### `room_plan_meeting`

```yaml
type: composite
risk_level: medium
auth_required: true
confirmation_policy: soft
status: planned
depends_on:
  - room_get_availability
```

**Purpose**: Đề xuất phòng phù hợp từ danh sách phòng trống.

**Input schema**: `QueryRoomsAvailability` + tùy chọn tiêu chí (`capacity`, `building`, `equipment[]`).

**Rule**:

1. Gọi `room_get_availability`.
2. Lọc theo tiêu chí.
3. Trả danh sách sắp xếp theo độ phù hợp.

---

## 4. Implementation Status Snapshot

| Tool | Status |
|---|---|
| `auth_login` | implemented |
| `student_get_profile` | implemented |
| All tools còn lại | planned |

---

## 5. Minimum Edge Cases (áp dụng cho mọi tool)

1. Thiếu field required.
2. Sai format ngày/email/URI.
3. Sai enum.
4. Thiếu hoặc token không hợp lệ.
5. Backend timeout hoặc 503.
6. Endpoint chưa triển khai (status planned).
