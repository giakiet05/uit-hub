# Fake Server Design

Tài liệu này mô tả cách code tiếp `apps/fake-uit-server` theo đúng structure hiện tại. Fake server không dùng DB, không cần repository ở giai đoạn này. Data có thể static trong service, nhưng code vẫn phải giữ layer rõ ràng để dễ mở rộng theo `api-docs.md`.

## Base Path

Tất cả API chính đặt dưới:

```text
/api/v1
```

Ví dụ:

```text
GET /api/v1/students
GET /api/v1/students/:id
```

`/ping` chỉ dùng để health check nhanh, không nằm trong `/api/v1`.

## Current Structure

```text
apps/fake-uit-server/
├── cmd/server/main.go
├── internal/
│   ├── apperror/
│   ├── bootstrap/
│   ├── config/
│   ├── controller/
│   ├── dto/
│   ├── middleware/
│   ├── model/
│   ├── route/
│   └── service/
├── .env
├── .env.example
├── go.mod
└── go.sum
```

## Request Flow

Request đi qua các tầng theo thứ tự:

```text
cmd/server/main.go
  -> bootstrap.Init()
  -> Gin router
  -> middleware.Chaos()
  -> middleware.RequireAuth() nếu route protected
  -> route.Register...Routes()
  -> controller
  -> service
  -> static model data
```

Response đi ngược lại:

```text
service returns model
  -> controller maps model to dto
  -> dto.SendSuccess / dto.SendError
  -> JSON response
```

## Layer Rules

### route

`internal/route` chỉ khai báo URL và nối route tới controller method.

Không viết business logic trong route.

Ví dụ:

```go
func RegisterStudentRoutes(rg *gin.RouterGroup, c *controller.StudentController) {
	students := rg.Group("/students")

	students.GET("", c.GetStudents)
	students.GET("/:id", c.GetStudentByID)
}
```

### controller

`internal/controller` nhận request, đọc param/query/body, gọi service, map model sang DTO, rồi trả response.

Controller được phép quyết định HTTP status code.

Controller không chứa static data và không tự format response JSON thủ công. Dùng `dto.SendSuccess`, `dto.SendError`, hoặc `dto.AbortWithError`.

Ví dụ:

```go
func (c *StudentController) GetStudentByID(ctx *gin.Context) {
	student, found := c.service.GetStudentByID(ctx.Param("id"))
	if !found {
		err := apperror.ErrNotFound
		err.Message = "Student not found"
		dto.SendError(ctx, err)
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, dto.NewStudentResponse(*student))
}
```

### service

`internal/service` chứa logic xử lý use case và static data cho fake server.

Service trả về `model`, không trả DTO.

Lý do: `model` là shape nội bộ, còn DTO là API contract. Controller là nơi map model sang DTO trước khi trả HTTP response.

Ví dụ:

```go
type StudentService interface {
	GetStudents() []model.Student
	GetStudentByID(id string) (*model.Student, bool)
}
```

### model

`internal/model` chứa internal data shape/entity. Không expose trực tiếp model ra response trong controller.

Ví dụ:

```go
type Student struct {
	ID    string
	Name  string
	Email string
	Phone string
	Major string
}
```

### dto

`internal/dto` chứa request/response DTO và response helpers.

DTO là contract giữa fake server và client. Nếu `api-docs.md` đổi response shape, sửa DTO trước.

Ví dụ:

```go
type StudentResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Major string `json:"major"`
}
```

Mapper nên đặt trong DTO package:

```go
func NewStudentResponse(student model.Student) StudentResponse
func NewStudentResponses(students []model.Student) []StudentResponse
```

### middleware

`internal/middleware` chứa middleware của Gin.

Hiện có:

```text
Chaos()       giả lập lỗi/delay
RequireAuth() kiểm tra Authorization Bearer token
```

Middleware không chứa business logic sâu. Nếu cần kiểm tra token, gọi `service.AuthService`.

### apperror

`internal/apperror` gom các lỗi chuẩn:

```text
BAD_REQUEST
UNAUTHORIZED
FORBIDDEN
NOT_FOUND
TOO_MANY_REQUESTS
SERVICE_UNAVAILABLE
INTERNAL_ERROR
```

Controller và middleware dùng `apperror` để trả lỗi thống nhất:

```json
{
  "message": "Student not found",
  "error_code": "NOT_FOUND"
}
```

### bootstrap

`internal/bootstrap` là nơi wiring dependency bằng `samber/do`.

Không viết business logic trong bootstrap.

Khi thêm service/controller mới, đăng ký dependency tại đây:

```go
do.Provide(injector, func(i *do.Injector) (service.StudentService, error) {
	return service.NewStudentService(), nil
})

do.Provide(injector, func(i *do.Injector) (service.AuthService, error) {
	return service.NewAuthService(), nil
})

do.Provide(injector, func(i *do.Injector) (*controller.StudentController, error) {
	studentService := do.MustInvoke[service.StudentService](i)
	return controller.NewStudentController(studentService), nil
})

do.Provide(injector, func(i *do.Injector) (*controller.AuthController, error) {
	authService := do.MustInvoke[service.AuthService](i)
	return controller.NewAuthController(authService), nil
})
```

Sau đó register route:

```go
api := router.Group("/api/v1")
api.Use(middleware.Chaos(cfg.Chaos))

authService := do.MustInvoke[service.AuthService](injector)

route.RegisterAuthRoutes(api, do.MustInvoke[*controller.AuthController](injector))
route.RegisterStudentRoutes(api, do.MustInvoke[*controller.StudentController](injector), authService)
```

## Auth Flow

Fake server có auth mỏng để client test flow đăng nhập và lỗi `401`.

Không dùng JWT thật, refresh token, bcrypt, session store hoặc DB.

### Public Routes

```text
POST /api/v1/login
GET  /api/v1/students
GET  /api/v1/students/:id
```

`/students` hiện là demo endpoint cũ nên tạm public.

### Protected Routes

```text
GET /api/v1/student/profile
```

Các endpoint dữ liệu cá nhân theo `api-docs.md` nên đặt dưới `/student/...` và dùng `RequireAuth`.

### Login

Request:

```bash
curl -X POST http://localhost:3000/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"student_id":"1","password":"password"}'
```

Response:

```json
{
  "message": "Successfully!",
  "data": {
    "student_id": "1",
    "access_token": "fake-token-1",
    "token_type": "Bearer",
    "remember_me": false
  }
}
```

### Calling Protected Routes

```bash
curl http://localhost:3000/api/v1/student/profile \
  -H "Authorization: Bearer fake-token-1"
```

Nếu không có token hoặc token sai:

```json
{
  "message": "Unauthorized!",
  "error_code": "UNAUTHORIZED"
}
```

### Adding Protected Routes

Trong route:

```go
student := rg.Group("/student")
student.Use(middleware.RequireAuth(authService))
{
	student.GET("/profile", c.GetProfile)
	student.GET("/score", c.GetScore)
}
```

Trong controller lấy student ID từ Gin context:

```go
studentID, ok := ctx.Get(middleware.StudentIDKey)
if !ok {
	dto.SendError(ctx, apperror.ErrUnauthorized)
	return
}
```

## Response Format

Success response:

```json
{
  "message": "Successfully!",
  "data": {}
}
```

Error response:

```json
{
  "message": "Internal error happened!",
  "error_code": "INTERNAL_ERROR"
}
```

Không tự tạo response map trong controller như:

```go
ctx.JSON(200, gin.H{"data": data})
```

Luôn dùng helper trong `internal/dto`.

## Chaos Middleware

`internal/middleware/chaos.go` dùng để giả lập lỗi hoặc delay cho fake server.

Mục tiêu là giúp client test các case xấu như server lỗi, request chậm, service unavailable, rate limit.

### Header override

Ép server trả status lỗi:

```bash
curl -H "X-Fake-Status: 500" http://localhost:3000/api/v1/students
```

Ép delay:

```bash
curl -H "X-Fake-Delay: 2000" http://localhost:3000/api/v1/students
```

`X-Fake-Delay` tính bằng milliseconds.

### Query override

Tương tự header, nhưng dùng query param:

```bash
curl "http://localhost:3000/api/v1/students?__fake_status=503"
curl "http://localhost:3000/api/v1/students?__fake_delay=2000"
```

### Env config

Config nằm trong `.env` và `.env.example`:

```env
FAKE_CHAOS_ENABLED=false
FAKE_ERROR_RATE=0
FAKE_MIN_DELAY_MS=0
FAKE_MAX_DELAY_MS=0
```

Ý nghĩa:

```text
FAKE_CHAOS_ENABLED  bật/tắt random chaos
FAKE_ERROR_RATE     tỷ lệ random trả lỗi 500, từ 0 tới 1
FAKE_MIN_DELAY_MS   delay thấp nhất khi random chaos bật
FAKE_MAX_DELAY_MS   delay cao nhất khi random chaos bật
```

Ví dụ config random lỗi 10%, delay 200-800ms:

```env
FAKE_CHAOS_ENABLED=true
FAKE_ERROR_RATE=0.1
FAKE_MIN_DELAY_MS=200
FAKE_MAX_DELAY_MS=800
```

Header/query override vẫn tiện nhất khi muốn test chính xác một case cụ thể.

## How To Add A New Endpoint

Ví dụ thêm endpoint theo `api-docs.md`:

```text
GET /student/profile
```

Vì base path là `/api/v1`, URL thật sẽ là:

```text
GET /api/v1/student/profile
```

### 1. Add model

Thêm hoặc mở rộng model trong `internal/model`.

```go
type StudentProfile struct {
	ID    string
	Name  string
	Email string
	Major string
}
```

### 2. Add DTO

Thêm response DTO trong `internal/dto`.

```go
type StudentProfileResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Major string `json:"major"`
}
```

Thêm mapper:

```go
func NewStudentProfileResponse(profile model.StudentProfile) StudentProfileResponse {
	return StudentProfileResponse{
		ID: profile.ID,
		Name: profile.Name,
		Email: profile.Email,
		Major: profile.Major,
	}
}
```

### 3. Add service method

Service trả model:

```go
type StudentService interface {
	GetProfile() model.StudentProfile
}
```

Implementation dùng static data:

```go
func (s *studentService) GetProfile() model.StudentProfile {
	return s.profile
}
```

### 4. Add controller method

Controller gọi service và map sang DTO:

```go
func (c *StudentController) GetProfile(ctx *gin.Context) {
	profile := c.service.GetProfile()
	dto.SendSuccess(ctx, http.StatusOK, dto.NewStudentProfileResponse(profile))
}
```

### 5. Register route

Trong `internal/route/student_route.go`:

```go
student := rg.Group("/student")
student.GET("/profile", c.GetProfile)
```

Không quên URL thật có prefix `/api/v1`.

Nếu endpoint cần đăng nhập, route đó phải nằm trong group có `middleware.RequireAuth(authService)`.

## Adding A New Domain

Nếu thêm domain mới, ví dụ `course`:

```text
internal/model/course.go
internal/dto/course_dto.go
internal/service/course_service.go
internal/controller/course_controller.go
internal/route/course_route.go
```

Sau đó đăng ký service/controller/route trong `internal/bootstrap/init.go`.

## What Not To Do

Không gọi service trực tiếp trong route.

Không trả model trực tiếp ra response.

Không tự tạo JSON response riêng trong từng controller.

Không thêm repo/DB khi data vẫn static.

Không viết business logic trong bootstrap.

Không dùng Fiber trong `apps/fake-uit-server`; project này dùng Gin.

## Verification

Sau khi thêm hoặc sửa endpoint, chạy:

```bash
cd apps/fake-uit-server
go test ./...
go run ./cmd/server
```

Hoặc chạy từ root monorepo:

```bash
go test ./apps/fake-uit-server/...
go run ./apps/fake-uit-server/cmd/server
go build -o /tmp/fake-uit-server ./apps/fake-uit-server/cmd/server
```

Không chạy build từ root kiểu này:

```bash
go build ./apps/fake-uit-server/cmd/server
```

Vì Go sẽ cố tạo binary tên `server` ở root, nhưng root đã có folder `server/` cũ.

Test nhanh:

```bash
curl -i http://localhost:3000/api/v1/students
curl -i http://localhost:3000/api/v1/students/1
curl -i -H "X-Fake-Status: 500" http://localhost:3000/api/v1/students
curl -i http://localhost:3000/api/v1/student/profile
curl -i -X POST http://localhost:3000/api/v1/login -H "Content-Type: application/json" -d '{"student_id":"1","password":"password"}'
curl -i http://localhost:3000/api/v1/student/profile -H "Authorization: Bearer fake-token-1"
```
