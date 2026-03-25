# API Specification (UIT Hub)

## Base URL

## Success Response Format

```json
{
  "message": "Successfully!",
  "data": {}
}
```

## Error Response Format

```json
{
  "message": "Internal error happened!",
  "error_code": "INTERNAL_ERROR"
}
```

---

## JSON Schema

### Common

#### SuccessResponse (generic)

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "SuccessResponse",
  "type": "object",
  "additionalProperties": false,
  "required": ["message", "data"],
  "properties": {
    "message": { "type": "string" },
    "data": {}
  }
}
```

#### ErrorResponse (generic)

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "ErrorResponse",
  "type": "object",
  "additionalProperties": false,
  "required": ["message", "error_code"],
  "properties": {
    "message": { "type": "string" },
    "error_code": { "type": "string" }
  }
}
```

#### Common Params

##### Query: year + semester

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "QueryYearSemester",
  "type": "object",
  "additionalProperties": false,
  "required": ["year", "semester"],
  "properties": {
    "year": { "type": "integer", "minimum": 1900 },
    "semester": { "type": "integer", "minimum": 1 }
  }
}
```

##### Query: rooms availability (date + start + end)

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "QueryRoomsAvailability",
  "type": "object",
  "additionalProperties": false,
  "required": ["date", "start", "end"],
  "properties": {
    "date": {
      "type": "string",
      "description": "ISO date (YYYY-MM-DD)",
      "pattern": "^\\d{4}-\\d{2}-\\d{2}$"
    },
    "start": {
      "type": "string",
      "description": "Start time/slot (format phụ thuộc backend)"
    },
    "end": {
      "type": "string",
      "description": "End time/slot (format phụ thuộc backend)"
    }
  }
}
```

##### Path: courseId

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "PathCourseId",
  "type": "object",
  "additionalProperties": false,
  "required": ["courseId"],
  "properties": {
    "courseId": { "type": "string" }
  }
}
```

##### Path: assignmentId

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "PathAssignmentId",
  "type": "object",
  "additionalProperties": false,
  "required": ["assignmentId"],
  "properties": {
    "assignmentId": { "type": "string" }
  }
}
```

##### Path: slugs

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "PathProgramSlugs",
  "type": "object",
  "additionalProperties": false,
  "required": ["slugs"],
  "properties": {
    "slugs": { "type": "string", "description": "Slug(s) path segment" }
  }
}
```

---

## API List

## 1. Auth

### 1.1 POST /login

Đăng nhập.

**Request JSON Schema:**

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "LoginRequest",
  "type": "object",
  "additionalProperties": false,
  "required": ["student_id", "password"],
  "properties": {
    "student_id": {
      "type": "string",
      "description": "MSSV (ví dụ 2252xxxx)"
    },
    "password": { "type": "string", "minLength": 1 },
    "remember_me": { "type": "boolean", "default": false }
  }
}
```

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

---

## 2. Student (daa.uit.edu.vn)

### 2.1 GET /student/profile

Lấy thông tin profile sinh viên.

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.2 GET /student/schedule

Lấy lịch học theo năm & học kỳ.

**Query:** `year`, `semester`

**Query JSON Schema:** QueryYearSemester

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.3 GET /student/schedule/exam

Lấy lịch thi theo năm & học kỳ.

**Query:** `year`, `semester`

**Query JSON Schema:** QueryYearSemester

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.4 GET /student/score

Lấy điểm / bảng điểm.

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.5 GET /student/lookup/tuitionfee

Tra cứu học phí.

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.6 GET /student/insurance

Tra cứu thông tin bảo hiểm.

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.7 GET /student/courses

Danh sách môn học / lớp học phần.

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.8 GET /student/training-points

Tra cứu điểm rèn luyện.

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.9 GET /student/lookup/office365

Tra cứu thông tin Office365.

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.10 GET /student/survey-form

Lấy survey form (nếu có).

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.11 POST /student/transcript-regis

Đăng ký transcript.

**Request JSON Schema:**

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "TranscriptRegisRequest",
  "type": "object",
  "additionalProperties": false,
  "required": ["copies", "delivery_method"],
  "properties": {
    "copies": { "type": "integer", "minimum": 1, "default": 1 },
    "language": { "type": "string", "enum": ["VI", "EN"], "default": "VI" },
    "delivery_method": { "type": "string", "enum": ["PICKUP", "SHIP"] },
    "shipping_address": {
      "type": "string",
      "minLength": 1,
      "description": "Bắt buộc nếu delivery_method=SHIP"
    },
    "phone": { "type": "string", "minLength": 6 },
    "note": { "type": "string" }
  },
  "allOf": [
    {
      "if": {
        "properties": { "delivery_method": { "const": "SHIP" } },
        "required": ["delivery_method"]
      },
      "then": { "required": ["shipping_address", "phone"] }
    }
  ]
}
```

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.12 POST /student/referral

Tạo yêu cầu giấy giới thiệu / referral.

**Request JSON Schema:**

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "ReferralRequest",
  "type": "object",
  "additionalProperties": false,
  "required": ["recipient", "purpose"],
  "properties": {
    "recipient": {
      "type": "string",
      "minLength": 1,
      "description": "Đơn vị/Người nhận (vd: Công ty/Phòng ban)"
    },
    "purpose": { "type": "string", "minLength": 1 },
    "from_date": {
      "type": "string",
      "description": "YYYY-MM-DD",
      "pattern": "^\\d{4}-\\d{2}-\\d{2}$"
    },
    "to_date": {
      "type": "string",
      "description": "YYYY-MM-DD",
      "pattern": "^\\d{4}-\\d{2}-\\d{2}$"
    },
    "contact_email": { "type": "string", "format": "email" },
    "note": { "type": "string" }
  }
}
```

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.13 POST /student/tuition-extend

Gia hạn học phí.

**Request JSON Schema:**

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "TuitionExtendRequest",
  "type": "object",
  "additionalProperties": false,
  "required": ["year", "semester", "requested_due_date", "reason"],
  "properties": {
    "year": { "type": "integer", "minimum": 1900 },
    "semester": { "type": "integer", "minimum": 1 },
    "requested_due_date": {
      "type": "string",
      "description": "YYYY-MM-DD",
      "pattern": "^\\d{4}-\\d{2}-\\d{2}$"
    },
    "reason": { "type": "string", "minLength": 1 },
    "phone": { "type": "string", "minLength": 6 },
    "attachments": {
      "type": "array",
      "description": "Danh sách URL/ID tệp minh chứng (nếu có)",
      "items": { "type": "string" },
      "default": []
    }
  }
}
```

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.14 POST /student/outpatient-regis

Đăng ký ngoại trú.

**Request JSON Schema:**

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "OutpatientRegisRequest",
  "type": "object",
  "additionalProperties": false,
  "required": ["address", "start_date", "landlord_name", "landlord_phone"],
  "properties": {
    "address": { "type": "string", "minLength": 1 },
    "province": { "type": "string" },
    "district": { "type": "string" },
    "ward": { "type": "string" },
    "start_date": {
      "type": "string",
      "description": "YYYY-MM-DD",
      "pattern": "^\\d{4}-\\d{2}-\\d{2}$"
    },
    "end_date": {
      "type": "string",
      "description": "YYYY-MM-DD",
      "pattern": "^\\d{4}-\\d{2}-\\d{2}$"
    },
    "landlord_name": { "type": "string", "minLength": 1 },
    "landlord_phone": { "type": "string", "minLength": 6 },
    "note": { "type": "string" },
    "attachments": {
      "type": "array",
      "items": { "type": "string" },
      "default": []
    }
  }
}
```

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.15 POST /student/eor-regis

Đăng ký EOR.

**Request JSON Schema:**

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "EorRegisRequest",
  "type": "object",
  "additionalProperties": false,
  "required": ["year", "semester", "subject_code", "reason"],
  "properties": {
    "year": { "type": "integer", "minimum": 1900 },
    "semester": { "type": "integer", "minimum": 1 },
    "subject_code": { "type": "string", "minLength": 1 },
    "subject_name": { "type": "string" },
    "reason": { "type": "string", "minLength": 1 },
    "contact_phone": { "type": "string", "minLength": 6 },
    "attachments": {
      "type": "array",
      "items": { "type": "string" },
      "default": []
    }
  }
}
```

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.16 POST /student/monthly-parking

Đăng ký gửi xe tháng.

**Request JSON Schema:**

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "MonthlyParkingRequest",
  "type": "object",
  "additionalProperties": false,
  "required": ["vehicle_type", "plate_number", "months", "start_month"],
  "properties": {
    "vehicle_type": {
      "type": "string",
      "enum": ["MOTORBIKE", "CAR", "BICYCLE"]
    },
    "plate_number": { "type": "string", "minLength": 3 },
    "months": { "type": "integer", "minimum": 1, "maximum": 12 },
    "start_month": {
      "type": "string",
      "description": "YYYY-MM",
      "pattern": "^\\d{4}-\\d{2}$"
    },
    "owner_name": { "type": "string" },
    "note": { "type": "string" }
  }
}
```

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.17 POST /student/graduate

Đăng ký tốt nghiệp.

**Request JSON Schema:**

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "GraduateRequest",
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

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.18 POST /student/graduation-thesis

Đăng ký khóa luận tốt nghiệp.

**Request JSON Schema:**

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "GraduationThesisRequest",
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

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

---

## 3. Courses (course.uit.edu.vn)

### 3.1 GET /student/deadlines

Danh sách deadlines/bài tập.

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 3.2 GET /student/courses/{courseId}/materials

Tài liệu môn học theo `courseId`.

**Path params:** `courseId`

**Path JSON Schema:** PathCourseId

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 3.3 GET /student/courses/{courseId}/assignments

Bài tập theo `courseId`.

**Path params:** `courseId`

**Path JSON Schema:** PathCourseId

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 3.4 POST /student/assignments/{assignmentId}/submissions

Nộp bài cho assignment.

**Path params:** `assignmentId`

**Path JSON Schema:** PathAssignmentId

**Request JSON Schema:**

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "AssignmentSubmissionRequest",
  "type": "object",
  "additionalProperties": false,
  "required": ["submission_type"],
  "properties": {
    "submission_type": { "type": "string", "enum": ["FILE", "LINK", "TEXT"] },
    "text": { "type": "string" },
    "url": { "type": "string", "format": "uri" },
    "file_urls": {
      "type": "array",
      "items": { "type": "string" },
      "default": []
    },
    "comment": { "type": "string" }
  },
  "allOf": [
    {
      "if": {
        "properties": { "submission_type": { "const": "TEXT" } },
        "required": ["submission_type"]
      },
      "then": { "required": ["text"] }
    },
    {
      "if": {
        "properties": { "submission_type": { "const": "LINK" } },
        "required": ["submission_type"]
      },
      "then": { "required": ["url"] }
    },
    {
      "if": {
        "properties": { "submission_type": { "const": "FILE" } },
        "required": ["submission_type"]
      },
      "then": { "required": ["file_urls"] }
    }
  ]
}
```

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

---

## 4. Announcement

### 4.1 GET /announcement

Danh sách announcement.

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

---

## 5. Rooms

### 5.1 GET /rooms/availability

Tra cứu phòng trống.

**Query:** `date`, `start`, `end`

**Query JSON Schema:** QueryRoomsAvailability

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

---

### 6.1 POST /contact

Gửi liên hệ.

**Request JSON Schema:**

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "ContactRequest",
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

**Response (200 OK):**

```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse
