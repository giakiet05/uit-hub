# Định hướng UIT Agent

## Bối cảnh

UIT Agent là agent hỗ trợ sinh viên, không phải coding agent chạy local như Claude Code. Agent này sẽ tương tác với các hệ thống backend truyền thống của trường và về sau đóng vai trò backend cho frontend web hoặc mobile.

Giai đoạn đầu sẽ chạy bằng CLI đơn giản để test nhanh core agent. Cách này giúp mình tránh bị kéo vào những vấn đề riêng của backend/web trước khi agent loop, tool system và luồng hội thoại ổn định.

Project sẽ ưu tiên tự build bằng Go, không dựa vào framework agent như LangGraph. Lý do là framework agent có thể chậm, khó kiểm soát và tạo thêm độ phức tạp không cần thiết cho core runtime.

## Claude Code là reference, không phải bản thiết kế để copy

Bộ tài liệu `claude-code-from-source-main` được dùng làm bản đồ kiến trúc. Mình sẽ đọc theo từng phase, làm tới đâu đọc tới đó, tránh đọc hết một lượt và kéo quá nhiều context không liên quan.

Cần lưu ý khác biệt domain:

- Claude Code là local coding agent, tập trung mạnh vào filesystem, shell command, permission, bảo mật môi trường local, tài nguyên máy và terminal UI.
- UIT Agent là student assistant, tập trung vào workflow sinh viên, API/backend trường, dữ liệu học vụ, hội thoại và tool gọi service.

Vì vậy, các abstraction của Claude Code sẽ được điều chỉnh theo domain UIT thay vì copy nguyên xi.

## Roadmap lớn (Cập nhật)

1. Core CLI và bootstrap (Đã xong)
2. Agent loop (ReAct & Plan-and-Execute) (Đã xong)
3. Tool system MVP (Đã xong)
4. Nâng cấp tool: Concurrency, JSON Schema metadata (Đã xong)

**CÁC TÍNH NĂNG TIẾP THEO (Theo thứ tự ưu tiên mới):**

5. **Human in the loop (Đã xong):** Cho phép agent tạm dừng để hỏi xin ý kiến hoặc chờ user duyệt một tool nguy hiểm.
6. **Compact conversation:** Thu gọn/nén lịch sử chat (tóm tắt các vòng lặp cũ) để tiết kiệm token khi context phình to.
7. **Persistent conversation:** Lưu lịch sử chat vào Database thay vì chỉ giữ trong RAM, cho phép user tiếp tục session cũ.
8. **Skills and hooks:** Thêm hook can thiệp trước/sau khi gọi tool, gọi model (phục vụ logging, guardrails).
9. **Cancel context / Interrupt signal:** Xử lý bắt tín hiệu Ctrl+C từ người dùng để ngắt luồng an toàn (cancel context, abort goroutines).
10. **Multi-agent:** Kiến trúc đa agent chuyên biệt phối hợp làm việc (Để làm CUỐI CÙNG).
11. **Backend web mode:** Bọc agent thành service API.
12. **Production hardening:** Tối ưu performance, caching, eval, rate limit.

## Chiến lược đọc tài liệu Claude Code

Không đọc toàn bộ tài liệu một lúc. Đọc theo nhu cầu của từng phase:

- Bootstrap: `ch02-bootstrap.md`
- State: `ch03-state.md`
- API/model layer: `ch04-api-layer.md`
- Agent loop: `ch05-agent-loop.md`
- Tool system: `ch06-tools.md`
- Tool concurrency: `ch07-concurrency.md`
- Sub-agent và task: `ch08-sub-agents.md`, `ch09-fork-agents.md`, `ch10-coordination.md`
- Memory: `ch11-memory.md`
- Hooks/extensibility: `ch12-extensibility.md`
- CLI interaction: `ch14-input-interaction.md`
- MCP/remote/backend hóa: `ch15-mcp.md`, `ch16-remote.md`
- Performance: `ch17-performance.md`

Nguyên tắc đọc: lấy mục lục trước, chỉ đọc đúng section liên quan, sau đó rút ra bài học có thể áp dụng cho UIT Agent.
