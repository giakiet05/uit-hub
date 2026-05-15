# Cấu trúc project Agent

## Tư duy thiết kế

Agent trong `apps/agent` được thiết kế như một core runtime tái sử dụng, không code cứng theo domain UIT. Agent không biết trực tiếp về sinh viên, môn học, lịch học hay API của trường. Những thông tin đó đi vào agent qua:

- Prompt và system instructions.
- MCP server chứa tool/API của từng trường.
- RAG MCP server hoặc các nguồn context bên ngoài.

Cách này giúp agent core có thể public và tái sử dụng cho trường khác. Muốn dùng cho domain khác thì đổi prompt, đổi MCP server, đổi nguồn RAG; code core không cần đổi.

## Cây thư mục

```text
apps/agent/
  cmd/
    agent/

  internal/
    app/
    config/
    cli/
    runtime/
    agent/
    llm/
    prompt/
    tool/
    mcp/
    localtool/
    memory/
    policy/
    hook/
    telemetry/
    testharness/

  testdata/
```

## Vai trò từng folder

### `cmd/agent`

Entry point của CLI agent. Package này chỉ nên parse tham số chạy, gọi bootstrap/wiring và chuyển quyền điều khiển cho app. Không đặt agent loop, tool logic hay provider logic ở đây.

### `internal/app`

Lớp khởi động và wiring cấp ứng dụng. Nhiệm vụ chính là ghép các subsystem lại với nhau: config, model client, MCP clients, tool registry, prompt builder, runtime session và CLI runner.

### `internal/config`

Định nghĩa cấu hình generic cho agent: model provider, MCP servers, prompt path, logging, timeout, memory mode và các tuỳ chọn runtime. Config không được chứa logic riêng của UIT.

### `internal/cli`

Presentation layer cho CLI. Ban đầu chỉ cần CLI đơn giản để nhập prompt, hiển thị output và test nhanh core agent. Folder này không nên chứa logic nghiệp vụ hay orchestration sâu.

### `internal/runtime`

Quản lý runtime/session state cấp thấp: session id, conversation state, lifecycle, cancellation, timeout, mode chạy và các dữ liệu vận hành chung. Runtime không gọi trực tiếp model hay tool nếu không thông qua agent loop.

### `internal/agent`

Lõi agent loop. Package này quản lý message history, gọi model, nhận tool call, chạy tool thông qua executor, append tool result và quyết định khi nào dừng. Agent loop chỉ phụ thuộc vào interface generic, không phụ thuộc domain.

### `internal/llm`

Abstraction cho LLM provider. Package này chịu trách nhiệm gửi request tới model, nhận response hoặc streaming response, chuẩn hoá message/tool-call format cho agent loop. Các provider cụ thể như OpenAI, Gemini, Claude hoặc echo/mock provider sẽ nằm dưới package này.

### `internal/prompt`

Xây prompt và message context đưa vào model. Đây là nơi nạp system prompt, instructions, memory context, MCP tool descriptions và các template cần thiết. Domain như UIT chỉ nên đi vào agent qua prompt/config, không đi vào code core.

### `internal/tool`

Native tool contract của agent. Package này định nghĩa interface tool, registry, executor, input schema, result và error format. Agent loop chỉ biết tool thông qua contract này.

### `internal/mcp`

MCP client và MCP adapter. Package này kết nối MCP server, discover MCP tools, chuyển MCP tool thành native tool trong `internal/tool`, rồi gọi MCP tool khi agent yêu cầu execute.

### `internal/localtool`

Chứa các local tool nếu sau này thật sự cần. Ví dụ tool tính toán đơn giản, thời gian hiện tại hoặc helper không phụ thuộc domain. Hiện tại đa số tool sẽ đến từ MCP, nên folder này có thể để trống trong MVP.

### `internal/memory`

Quản lý memory generic của agent. Memory không nên giả định domain UIT; nó chỉ lưu và truy xuất thông tin hội thoại, user preference hoặc context dài hạn theo cơ chế generic.

### `internal/policy`

Policy và guardrail generic cho tool/model/runtime. Ví dụ permission mode, allow/deny tool, timeout policy, mutation policy hoặc safety checks. Policy nên dựa trên metadata của tool và config, không hardcode domain.

### `internal/hook`

Lifecycle hooks cho agent runtime. Hook có thể chạy trước/sau model call, trước/sau tool call, khi session start/end hoặc khi có error. Đây là điểm mở rộng để thêm logging, guardrail, context injection hoặc custom behavior.

### `internal/telemetry`

Logging, tracing, metrics và usage tracking. Package này giúp quan sát agent runtime mà không làm bẩn agent loop.

### `internal/testharness`

Test harness và eval utilities cho agent. Về sau có thể dùng để chạy kịch bản hội thoại, mock model, mock MCP server và kiểm tra expected workflows.

### `testdata`

Fixture dùng cho test: sample prompt, sample MCP schema, sample model response, transcript hoặc eval scenario. Không chứa secret hoặc dữ liệu thật nhạy cảm.

## Nguyên tắc ranh giới

- Agent core không có package domain như `uit`, `student`, `academic` hay `schedule`.
- Tool nghiệp vụ nằm trong MCP server, không nằm trong agent core.
- Agent chỉ thấy tool qua native tool contract và MCP adapter.
- Prompt là nơi đưa domain knowledge vào agent.
- CLI chỉ là presentation đầu tiên; sau này có thể thêm backend web mode mà không đổi agent loop.
- Các package trong `internal/` phải ưu tiên interface generic để project có thể tái sử dụng như một agent framework.
