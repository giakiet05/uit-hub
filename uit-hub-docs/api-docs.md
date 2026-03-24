
# API Specification (UIT Hub)

## Base URL
- **Localhost:** `http://localhost:8080/api/v1`
- **Remote:** `https://simo.giakiet.io.vn/api/v1`

> Ghi chú: Tất cả endpoint bên dưới được hiểu là **relative path** so với Base URL.

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

## JSON Schema (Common)

### SuccessResponse (generic)
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

### ErrorResponse (generic)
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

### Common Params

#### Query: year + semester
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

#### Query: rooms availability (date + start + end)
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
		"start": { "type": "string", "description": "Start time/slot (format phụ thuộc backend)" },
		"end": { "type": "string", "description": "End time/slot (format phụ thuộc backend)" }
	}
}
```

#### Path: courseId
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

#### Path: assignmentId
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

#### Path: slugs
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

**Request:** (chưa có mô tả chi tiết trong tài liệu hiện tại)

**Request JSON Schema (generic):**
```json
{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "LoginRequest",
	"type": "object",
	"additionalProperties": true
}
```

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

---

## 2. Student

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

### 2.11 GET /student/deadlines
Danh sách deadlines/bài tập.

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.12 GET /student/courses/{courseId}/materials
Tài liệu môn học theo `courseId`.

**Path params:** `courseId`

**Path JSON Schema:** PathCourseId

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.13 GET /student/courses/{courseId}/assignments
Bài tập theo `courseId`.

**Path params:** `courseId`

**Path JSON Schema:** PathCourseId

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.14 POST /student/transcript-regis
Đăng ký transcript.

**Request:** (chưa có mô tả chi tiết trong tài liệu hiện tại)

**Request JSON Schema (generic):**
```json
{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "TranscriptRegisRequest",
	"type": "object",
	"additionalProperties": true
}
```

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.15 POST /student/referral
Tạo yêu cầu referral.

**Request:** (chưa có mô tả chi tiết trong tài liệu hiện tại)

**Request JSON Schema (generic):**
```json
{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "ReferralRequest",
	"type": "object",
	"additionalProperties": true
}
```

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.16 POST /student/tuition-extend
Gia hạn học phí.

**Request:** (chưa có mô tả chi tiết trong tài liệu hiện tại)

**Request JSON Schema (generic):**
```json
{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "TuitionExtendRequest",
	"type": "object",
	"additionalProperties": true
}
```

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.17 POST /student/outpatient-regis
Đăng ký ngoại trú.

**Request:** (chưa có mô tả chi tiết trong tài liệu hiện tại)

**Request JSON Schema (generic):**
```json
{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "OutpatientRegisRequest",
	"type": "object",
	"additionalProperties": true
}
```

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.18 POST /student/eor-regis
Đăng ký EOR.

**Request:** (chưa có mô tả chi tiết trong tài liệu hiện tại)

**Request JSON Schema (generic):**
```json
{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "EorRegisRequest",
	"type": "object",
	"additionalProperties": true
}
```

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.19 POST /student/monthly-parking
Đăng ký gửi xe tháng.

**Request:** (chưa có mô tả chi tiết trong tài liệu hiện tại)

**Request JSON Schema (generic):**
```json
{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "MonthlyParkingRequest",
	"type": "object",
	"additionalProperties": true
}
```

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.20 POST /student/graduate
Đăng ký tốt nghiệp.

**Request:** (chưa có mô tả chi tiết trong tài liệu hiện tại)

**Request JSON Schema (generic):**
```json
{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "GraduateRequest",
	"type": "object",
	"additionalProperties": true
}
```

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.21 POST /student/graduation-thesis
Đăng ký khóa luận tốt nghiệp.

**Request:** (chưa có mô tả chi tiết trong tài liệu hiện tại)

**Request JSON Schema (generic):**
```json
{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "GraduationThesisRequest",
	"type": "object",
	"additionalProperties": true
}
```

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 2.22 POST /student/assignments/{assignmentId}/submissions
Nộp bài cho assignment.

**Path params:** `assignmentId`

**Path JSON Schema:** PathAssignmentId

**Request:** (chưa có mô tả chi tiết trong tài liệu hiện tại)

**Request JSON Schema (generic):**
```json
{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "AssignmentSubmissionRequest",
	"type": "object",
	"additionalProperties": true
}
```

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

---

## 3. Notifications

### 3.1 GET /notifications
Danh sách notifications.

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

---

## 4. Rooms

### 4.1 GET /rooms/availability
Tra cứu phòng trống.

**Query:** `date`, `start`, `end`

**Query JSON Schema:** QueryRoomsAvailability

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

---

## 5. Other

### 5.1 GET /annual-plan
Lấy annual plan.

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 5.2 GET /program/{slugs}
Lấy thông tin chương trình theo slug.

**Path params:** `slugs`

**Path JSON Schema:** PathProgramSlugs

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 5.3 GET /tutorial
Lấy tutorial.

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 5.4 POST /contact
Gửi liên hệ.

**Request:** (chưa có mô tả chi tiết trong tài liệu hiện tại)

**Request JSON Schema (generic):**
```json
{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "ContactRequest",
	"type": "object",
	"additionalProperties": true
}
```

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse

### 5.5 POST /verification
Gửi/kiểm tra verification.

**Request:** (chưa có mô tả chi tiết trong tài liệu hiện tại)

**Request JSON Schema (generic):**
```json
{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"title": "VerificationRequest",
	"type": "object",
	"additionalProperties": true
}
```

**Response (200 OK):**
```json
{ "message": "Successfully!", "data": {} }
```

**Response JSON Schema:** SuccessResponse