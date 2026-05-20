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

## Roadmap lớn

1. Core CLI và bootstrap
   - Tạo CLI mỏng để nhận prompt, đọc config cần thiết, chọn provider/model và khởi chạy agent.
   - Chưa cần TUI phức tạp, chưa cần backend web.

2. Agent loop
   - Xây vòng lặp chính gồm message history, gọi model, nhận response, phát hiện tool call, append tool result và lặp đến khi hoàn thành.
   - Đây là lõi sống còn của toàn bộ agent.

3. Tool system MVP
   - Định nghĩa tool interface tự mô tả: name, description, schema/input và execute.
   - Ban đầu chỉ cần vài tool gọi fake UIT server hoặc mock data.
   - Đến mốc này phải có single-agent MVP dùng được.

4. State và session
   - Tách runtime/session state khỏi conversation state.
   - Lưu message history tối thiểu, mode chạy, context sinh viên và thông tin token/cost nếu có.

5. Memory
   - Thêm memory sau khi agent loop và tool system đã ổn.
   - Memory có thể gồm user profile, thông tin học tập đã biết, sở thích, và các nội dung quan trọng từ lịch sử hội thoại.

6. Nâng cấp tool
   - Thêm timeout, retry, structured errors, policy/permission nhẹ, result rendering và concurrency khi cần.
   - Với domain UIT, ưu tiên độ tin cậy API và tính đúng của dữ liệu hơn là shell safety.

7. Multi-agent và task orchestration
   - Khi single agent bắt đầu phình, tách thành các sub-agent/task như academic advisor, schedule agent, course lookup agent, policy/FAQ agent.
   - Sub-agent nên dùng cùng core loop, không tạo code path riêng.

8. Hooks và extensibility
   - Thêm hook trước/sau tool call, trước/sau model call, logging, guardrail và context injection.
   - Đây là nền tảng hữu ích khi chuyển sang backend service.

9. Backend web mode
   - Khi CLI core đã ổn, bọc agent thành service API.
   - Thêm session endpoint, streaming response, user auth, persistence và observability.

10. Production hardening
    - Tối ưu performance, token budget, caching, eval/test scenarios, monitoring, rate limit và failure recovery.

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
