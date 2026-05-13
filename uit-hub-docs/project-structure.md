# Project Structure: Fake UIT Server

Tài liệu này quy định cấu trúc thư mục và trách nhiệm của từng tầng trong dịch vụ `fake-uit-server`. Dự án áp dụng **Layered Architecture** (theo hướng Clean/Hexagonal Architecture) kết hợp với **Dependency Injection (DI)** thông qua thư viện `samber/do`.

## Tổng quan sơ đồ thư mục (Monorepo)

```text
uit-hub/
├── docs/                        <-- Tài liệu dự án
│   ├── api-docs.md
│   ├── project-structure.md     <-- (Bạn đang ở đây)
│   └── convention.md
├── apps/                        <-- Chứa các dịch vụ
│   ├── fake-uit-server/         <-- Backend Go (Fake Server)
│   │   ├── cmd/                 <-- Entry points (chương trình thực thi)
│   │   │   └── server/
│   │   │       └── main.go      <-- File chạy app chính
│   │   ├── internal/            <-- Source code private (không thể import từ ngoài)
│   │   │   ├── bootstrap/       <-- Khởi tạo DI container và wiring app
│   │   │   ├── bus/             <-- In-memory EventBus cho nội bộ dự án
│   │   │   ├── config/          <-- Quản lý cấu hình (.env, env vars)
│   │   │   ├── controller/      <-- Tầng giao tiếp HTTP (Nhận request, gọi Service)
│   │   │   ├── infrastructure/  <-- Code thực thi kết nối bên ngoài (DB, Redis, Email, Mock Data)
│   │   │   │   └── mock/        <-- Chứa logic truy xuất data tĩnh giả lập (Mock Repositories)
│   │   │   ├── middleware/      <-- Các middleware cho HTTP (vd: Chaos Middleware tạo delay/lỗi)
│   │   │   ├── model/           <-- Data structs, Entities, DTOs
│   │   │   ├── route/           <-- Đăng ký đường dẫn HTTP cho các controller
│   │   │   ├── service/         <-- Tầng chứa Business Logic cốt lõi & Hợp đồng (Interfaces)
│   │   │   ├── ws/              <-- Tầng giao tiếp WebSocket
│   │   │   └── worker/          <-- Xử lý tác vụ ngầm (Lắng nghe EventBus, Cronjob)
│   │   ├── go.mod
│   │   └── go.sum
│   ├── mcp-server/              <-- MCP Server
│   └── agent/                   <-- AI Agent
└── scripts/                     <-- Scripts tiện ích (Build, Docker, v.v.)
```

---

## Trách nhiệm của từng thư mục (internal/)

### 1. `bootstrap/`
Nơi duy nhất thực hiện việc "kết nối" (wiring) các thành phần của ứng dụng. Sử dụng `samber/do` để quản lý Dependency Injection.
- **Quy tắc:** Không chứa logic nghiệp vụ. Chỉ dùng để đăng ký các hàm `New...` vào Injector.

### 2. `controller/` & `worker/`
Nơi tiếp nhận các yêu cầu từ phía ngoài (HTTP REST, EventBus, Cronjob).
- **Quy tắc:** Chỉ thực hiện việc xác thực đầu vào (Validation) và trả về kết quả (HTTP Status Code). Tuyệt đối không viết logic nghiệp vụ hay gọi trực tiếp xuống Database. Phải gọi thông qua `service`.

### 3. `service/`
Chứa toàn bộ Business Logic (Nghiệp vụ). 
- **Quy tắc Vàng (Interfaces belong to Consumer):** Tầng `service` cần dữ liệu gì thì sẽ tự định nghĩa ra `Interface` (hợp đồng) ngay bên trong thư mục của nó (vd: `StudentRepository interface`). Tầng này không cần quan tâm thành phần nào thực thi hợp đồng đó, dữ liệu đến từ Mongo hay File JSON tĩnh.

### 4. `infrastructure/`
Chứa các thành phần thực thi (Implementation) phụ thuộc vào công nghệ bên ngoài.
- **Bao gồm:** Kết nối MongoDB, Redis, AI Client, SMTP Email, và các **Repositories**.
- Đối với Fake Server, thư mục này sẽ chứa `infrastructure/mock/` (Đọc dữ liệu từ file tĩnh JSON). Các struct tại đây cần được cài đặt (implement) sao cho tương thích với `Interface` mà tầng `service` đã định nghĩa.
- **Quy tắc:** Nếu chuyển từ Mock JSON sang dùng DB thật, chỉ cần tạo thêm thư mục `infrastructure/postgres/`, viết đoạn mã thực thi mới và cập nhật cấu hình ở tầng bootstrap. Tầng `service` hoàn toàn không bị ảnh hưởng.

### 5. `middleware/`
Chứa các hàm đứng giữa request và controller.
- **Đặc biệt cho Fake Server:** Cần cài đặt `ChaosMiddleware` tại đây để giả lập độ trễ mạng (delay), tỷ lệ rớt server (503), quá tải (429) giúp Fake Server mô phỏng sát với thực tế nhất.

---

## Nguyên tắc cốt lõi (Core Principles)

1.  **Dependency Inversion:** Luôn nhận vào `Interface` ở Constructor, không nhận vào Struct cụ thể của layer khác.
2.  **No `repo/` Folder:** Không có thư mục `repo/` ngang hàng với `service`. Interface định nghĩa ở `service/`, Code thực thi truy xuất data đặt ở `infrastructure/`.
3.  **No Global State:** Không sử dụng biến toàn cục cho DB, Config, Logger. Mọi thứ phải được truyền qua DI.
4.  **Implicit Interface:** Trong Go, implementation không cần phải khai báo `implements InterfaceX`. Chỉ cần cung cấp đủ các hàm theo đúng chữ ký là hợp lệ. Mọi Interface phải được định nghĩa tại nơi tiêu thụ (Consumer).
