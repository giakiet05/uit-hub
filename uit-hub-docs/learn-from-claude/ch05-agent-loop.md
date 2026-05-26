# Chapter 5: Agent Loop

## Claude Code làm gì

Claude Code coi agent loop là trung tâm của toàn bộ hệ thống. Một lần user gửi task sẽ chạy một loop gồm nhiều iteration:

1. Build context.
2. Call model.
3. Stream response.
4. Execute tools.
5. Append observations.
6. Continue hoặc stop với terminal reason rõ ràng.

Loop này không phải toàn bộ session. Session có thể chứa nhiều user turns; mỗi user turn gọi agent loop một lần.

## Bài học chính

Agent không phải là một LLM call. Agent là một state machine chạy nhiều round cho tới khi hoàn thành, lỗi, bị abort, hết budget, hoặc chạm giới hạn.

Điểm quan trọng nhất là state transition phải explicit. Mỗi lần continue nên biết lý do vì sao tiếp tục, ví dụ:

- `next_turn`
- `max_rounds`
- `tool_error`
- `model_error`
- `completed`
- `aborted`

## Áp dụng cho UIT Hub Agent

Agent hiện tại đã có ReAct baseline:

- model call
- tool call
- observation
- continue nhiều round
- event stream về TUI

Nhưng loop vẫn nên được nâng cấp dần để rõ state hơn.

### Nên làm sau

- Tạo `RunState` cho một lần `agent.Run`.
- Tạo `TransitionReason` cho lý do tiếp tục.
- Tạo `TerminalReason` cho lý do dừng.
- Mỗi round nên emit event rõ ràng hơn.
- Đảm bảo mọi tool call đều có observation, kể cả lỗi hoặc timeout.
- Thêm tool result size budget để tránh observation quá lớn.

### Chưa làm vội

- Stop hooks.
- Auto compact.
- Reactive compact.
- Model fallback.
- Token budget continuation.
- Streaming tool execution concurrent phức tạp.

## Context Management

Claude Code có nhiều tầng context management:

- tool result budget
- snip compact
- microcompact
- context collapse
- auto compact

Project mình chưa cần làm hết. Tầng đáng làm sớm nhất là tool result budget, vì tool/MCP có thể trả dữ liệu lớn và làm phình conversation.

## Error Recovery

Mọi retry phải có circuit breaker. Không có retry nào được chạy vô hạn.

Những giới hạn nên có về sau:

- LLM retry max attempts.
- Tool retry max attempts nếu có.
- Compact retry max attempts.
- Max rounds per run.

Hiện `max rounds` đã là circuit breaker đầu tiên.

## Stop Hooks

Stop hook là cơ chế kiểm tra khi model nghĩ nó đã xong.

Ý tưởng này rất hay nhưng để sau. Với UIT Hub Agent, stop hook có thể dùng để kiểm tra:

- task yêu cầu ghi file nhưng chưa ghi
- task yêu cầu tạo advisor note nhưng chưa có note id
- task yêu cầu dùng source-system data nhưng chưa gọi tool
- task dài còn thiếu bước

## Multi-Agent

Chap 5 chưa nói sâu về orchestration multi-agent. Cơ chế sub-agent wrapped as tool sẽ nằm ở các chương sau.

Ý tưởng quan trọng cần giữ lại: agent chính không cần biết agent con hoạt động thế nào. Agent con có thể được expose như một tool:

1. Main agent gọi sub-agent tool.
2. Sub-agent tự làm việc.
3. Main agent nhận kết quả như tool observation bình thường.

## Kết luận

Chap 5 xác nhận việc cần nâng cấp core loop, nhưng không nên đập toàn bộ ngay. Hướng tốt nhất là dùng ReAct agent hiện tại làm baseline, chuẩn hóa eval trước, rồi refactor loop từng bước để so sánh với Plan-and-Execute sau này.
