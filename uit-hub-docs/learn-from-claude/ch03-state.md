# Chapter 3: State

## Claude Code làm gì

Claude Code chia state thành hai tầng:

- Bootstrap state: singleton mutable sống theo process, dùng cho infrastructure như session id, cwd, model override, cost tracking, prompt cache context, telemetry.
- AppState: reactive store cho UI, dùng để render và trigger side effect khi state đổi.

Điểm quan trọng không phải là copy đúng hai tầng này, mà là chia state theo access pattern:

- Dữ liệu sống bao lâu?
- Ai được đọc và ghi?
- Khi nó đổi thì có cần side effect không?

## Bài học chính

State boundary là bản đồ quyền sở hữu dữ liệu theo vòng đời. Nếu không có bản đồ này, code agent sẽ dần rối vì conversation, prompt, tool registry, memory, MCP connection, usage stats, TUI state bị truyền lẫn vào nhau.

Claude Code cũng giữ nhiều context theo session snapshot để bảo vệ prompt cache. Ví dụ memory/project context được load một lần đầu session, không rebuild tùy tiện mỗi turn.

## Áp dụng cho UIT Hub Agent

Project mình nên chia state theo bốn nhóm:

### App State

Sống suốt process.

Ví dụ:

- config
- LLM provider
- base tool registry
- MCP server manager
- memory store
- logger
- prompt builder factory hoặc session prompt builder config

App state không nên chứa conversation của user.

### Session State

Sống trong một phiên chat.

Ví dụ:

- session id
- conversation
- system prompt snapshot
- memory snapshot đã load lúc start session
- MCP catalog snapshot
- runtime tool registry
- usage tracker theo session

Session state là nơi phù hợp cho runtime registry LRU của MCP tools. Tool đã load không nên sống mãi trong app state.

### Run State

Sống trong một lần user gửi prompt.

Ví dụ:

- round count
- max rounds
- số LLM calls trong run
- số tool calls trong run
- temporary events
- tool failures
- duration

Run state nên bị bỏ sau khi agent trả lời xong, nhưng có thể cộng dồn vào usage tracker của session.

### UI State

Chỉ phục vụ TUI.

Ví dụ:

- input text
- cursor position
- scroll position
- pane width/height
- rendered conversation lines
- rendered log lines
- status text

TUI không nên sở hữu MCP manager, memory store, provider, hay tool registry.

## Áp dụng ngay

Chưa refactor lớn ngay. Hiện tại chỉ cần dùng note này làm nguyên tắc cho các feature tiếp theo.

Những điểm nên giữ:

- Prompt snapshot build theo session, không rebuild tùy tiện mỗi round.
- Memory write trong session hiện tại không rebuild system prompt ngay. Session hiện tại biết nhờ observation/conversation; session sau mới load memory vào prompt snapshot.
- MCP catalog load lúc start session, full tool schema deferred qua `load_mcp_tool`.
- Runtime registry có LRU và thuộc về session/run boundary, không thuộc base registry.

## Để sau

Khi code lớn hơn, cân nhắc refactor:

- Tạo struct session rõ ràng để gom conversation, prompt snapshot, runtime registry, usage tracker.
- Tách usage/cost tracker khỏi ReAct agent loop.
- Tập trung side effect vào runtime/session manager thay vì rải trong từng callsite.
- Chuẩn bị cho backend mode, nơi nhiều user/session chạy song song.

## Không copy từ Claude Code

Không cần copy reactive AppState riêng của Claude Code ở giai đoạn này. Claude Code cần nó vì dùng React/Ink. UIT Hub Agent hiện dùng Bubble Tea, bản thân TUI model/update đã là reactive state.

Không cần sticky latch prompt cache ngay. Ý tưởng này hay nhưng chỉ nên làm khi mình thật sự có feature toggles ảnh hưởng cache key.

## Kết luận

Chưa implement state architecture ngay. Đây là thiết kế dẫn đường để tránh coupling khi thêm usage tracker, session resume, compaction, multi-agent, hoặc backend service sau này.
