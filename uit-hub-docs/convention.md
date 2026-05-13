# Quy ước Phát triển (Development Conventions)

Tài liệu này quy định các tiêu chuẩn viết mã, quy ước đặt tên, xử lý lỗi và quy trình quản lý phiên bản mã nguồn (Git Workflow) cho dự án `fake-uit-server`. Mọi thành viên tham gia phát triển bắt buộc phải tuân thủ các quy tắc này để đảm bảo tính nhất quán và khả năng bảo trì của mã nguồn.

---

## 1. Quy ước Đặt tên (Naming Conventions)

### 1.1. Tên Package
- Phải là danh từ số ít, viết thường hoàn toàn (lowercase).
- Phản ánh đúng chức năng của package (ví dụ: `user`, `auth`, `payment`).
- **Không** sử dụng các tên vô nghĩa như `utils`, `helpers`, `common`.

### 1.2. Biến và Hàm (Variables & Functions)
- Sử dụng `camelCase` cho các biến/hàm private (không export).
- Sử dụng `PascalCase` cho các biến/hàm public (export).
- Tránh việc lặp lại tên package trong tên biến/hàm (vd: `user.UserLogin()` -> Sửa thành `user.Login()`).

### 1.3. Chữ viết tắt (Acronyms)
- Các từ viết tắt thông dụng (ID, HTTP, URL, JSON, API) phải viết hoa toàn bộ nếu là public, hoặc viết thường toàn bộ nếu đứng đầu chuỗi private.
- **Đúng:** `userID`, `studentID`, `HTTPClient`, `ParseJSON`.
- **Sai:** `userId`, `studentId`, `HttpClient`, `ParseJson`.

### 1.4. Interfaces
- Nếu Interface chỉ có 1 phương thức (method), tên Interface nên thêm hậu tố `-er` (ví dụ: `Reader`, `Writer`, `Authenticator`).
- Interface nên được định nghĩa tại **nơi tiêu thụ (Consumer)**, không phải tại nơi thực thi (Implementation).

---

## 2. Quy tắc Viết mã (Code Style & Architecture)

### 2.1. Dependency Injection & Trạng thái toàn cục (Global State)
- **Không** sử dụng biến toàn cục (`global variables`) để lưu trữ kết nối cơ sở dữ liệu, bộ đệm, hoặc cấu hình.
- Mọi phụ thuộc (dependencies) phải được tiêm (inject) qua hàm khởi tạo (Constructor).
- Nếu sử dụng `samber/do`, việc đăng ký (provide) phải được thực hiện ở tầng `bootstrap`.

### 2.2. Xử lý Lỗi (Error Handling)
- **Không** ẩn lỗi (swallow errors). Mọi lỗi trả về từ hàm phải được xử lý hoặc ném lên tầng trên.
- Sử dụng `fmt.Errorf("...: %w", err)` để bọc lỗi (wrap errors) nhằm giữ lại dấu vết (stack trace).
- **Early Return:** Kiểm tra và xử lý lỗi ngay lập tức để giữ cho luồng code chính yếu không bị thụt đầu dòng (indentation) quá sâu.
  ```go
  // ĐÚNG
  if err != nil {
      return err
  }
  // Logic chính tiếp tục ở thụt lề cấp 1
  ```

### 2.3. Mã nguồn & Định dạng
- Bắt buộc chạy `go fmt` (hoặc cấu hình IDE tự động format) trước khi commit mã nguồn.
- Khuyến khích cấu hình `golangci-lint` trong quá trình phát triển để bắt các lỗi phổ biến.

---

## 3. Quy trình Git (Git Workflow)

### 3.1. Đặt tên Nhánh (Branch Naming)
Sử dụng tiền tố để phân loại mục đích của nhánh:
- `feat/...`: Phát triển tính năng mới (vd: `feat/student-schedule`).
- `fix/...`: Sửa lỗi (vd: `fix/login-timeout`).
- `refactor/...`: Tái cấu trúc mã nguồn nhưng không đổi chức năng.
- `docs/...`: Cập nhật tài liệu.

### 3.2. Thông điệp Commit (Commit Messages)
Tuân thủ định dạng Conventional Commits:
`<type>[optional scope]: <description>`

- **Ví dụ:**
  - `feat(auth): implement login endpoint with JWT`
  - `fix(student): handle missing transcript data gracefully`
  - `docs(api): update JSON schema for exam schedule`

- Mô tả (description) phải rõ ràng, ngắn gọn và viết bằng tiếng Anh.
- Sử dụng động từ nguyên mẫu (imperative mood) ở đầu câu (vd: "add" thay vì "added" hay "adds").
