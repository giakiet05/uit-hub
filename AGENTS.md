# UIT Hub Agent Guide

This repo is a monorepo for UIT Hub. The active backend work is the fake UIT server in:

```text
apps/fake-uit-server
```

Project docs live in `uit-hub-docs/`, not `docs/`. Do not create or move docs into a `docs/` folder.

Read these first when needed:

- `uit-hub-docs/api-docs.md`
- `uit-hub-docs/project-structure.md`
- `uit-hub-docs/convention.md`
- `uit-hub-docs/fake-server-design.md`

The old `server/` folder is experimental code from an earlier attempt. Do not edit, migrate, or delete it unless explicitly asked.

Fake server stack:

- Go
- Gin
- `samber/do`
- `joho/godotenv`

Base API path:

```text
/api/v1
```

Current fake server layers:

```text
route -> middleware -> controller -> service -> model/static data
```

Rules:

- `route`: register URLs only.
- `controller`: bind request, call service, map model to DTO, choose HTTP status.
- `service`: fake business logic and static data; return `model`, not DTO.
- `model`: internal data shape.
- `dto`: API request/response structs, mappers, response helpers.
- `apperror`: shared app errors and error codes.
- `middleware`: Gin middleware such as auth and chaos.
- `bootstrap`: DI wiring and route mounting only.

Do not add DB/repository code while data is static.

Use `dto.SendSuccess`, `dto.SendError`, or `dto.AbortWithError` for API responses. Do not hand-roll `gin.H` response bodies in controllers.

Fake auth exists for protected student routes. Use `middleware.RequireAuth(authService)` for private endpoints.

Chaos middleware exists for testing failures and delay:

```text
X-Fake-Status
X-Fake-Delay
__fake_status
__fake_delay
```

`.env` is local and ignored. `.env.example` should be committed.

Use Gin for the fake server, not Fiber.

Keep code simple, but keep the layers clean so future endpoints follow the same pattern.
