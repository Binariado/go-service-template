# Coding Standards

## Naming

| Element | Convention | Example |
|---|---|---|
| Internal directories | singular | `handler/`, `dto/`, `repository/` |
| REST routes | plural nouns | `/users`, `/orders/{id}` |
| Migration files | `000NNN_snake_case.{up,down}.sql` | `000002_create_orders_table.up.sql` |
| SQL query files | `<module>.sql` | `user.sql` |
| Domain files | `entity.go`, `repository.go`, `service.go` | — |
| Use case files | `<verb>_<noun>.go` | `create_order.go` |
| Use case structs | `<Verb><Noun>UseCase` + `New...UseCase` | `CreateOrderUseCase` |
| Handler files | `<entity>.go`; package `<module>handler` | `order.go` |
| DTO structs | `<Entity>Response`, `Create<Entity>Request`, `Update<Entity>Request`; package `<module>dto` | `OrderResponse` |
| Permission constants | `Perm<Module><Entity><Action>` = `"module.entity.action"` | `PermOrderOrdersRead` |

## Import aliases

Use descriptive aliases when two packages share a name:

```go
domainuser  "github.com/binariado/go-service-template/internal/domain/user"
userhandler "github.com/binariado/go-service-template/internal/transport/rest/user/handler"
userdto     "github.com/binariado/go-service-template/internal/transport/rest/user/dto"
graph       "github.com/binariado/go-service-template/internal/transport/graphql"
```

## Rules

- Dependencies are injected through constructors; no globals, no `init()` wiring (except loading `.env`).
- Every I/O method takes `ctx context.Context` as its first parameter.
- Never read `os.Getenv` outside `config/`. Pass values through constructors.
- **Handlers** parse/validate input → call a use case → map errors → write a response. No business logic, no SQL.
- **Use cases**: one file, one `Execute(ctx, ...)`. Application-level validation (pagination defaults, sort allow-lists) lives here, via `application/common/listing`.
- **Domain**: no HTTP, SQL, or framework types. Errors that callers must distinguish are exported values (`var ErrX = errors.New(...)`) in the domain package.
- **Repositories** are the only place that executes SQL. Use sqlc; fall back to `QueryContext`/`ExecContext` only for dynamic queries.
- **DTOs** never cross into the domain; map entity → DTO explicitly (`toUserResponse`).
- Timestamps in REST responses: RFC3339 strings.
- Log with `log/slog` and always pass the request context (`slog.ErrorContext`). Log the real error, return a generic `InternalError` to the client.
- Best-effort side effects (e.g. audit logs) may discard their error with `_ =`, but must be commented as such.

## Errors → HTTP

| Situation | Helper | Status |
|---|---|---|
| Invalid input | `response.BadRequest` | 400 |
| Missing/invalid credentials | `response.Unauthorized` | 401 |
| Not allowed | `response.Forbidden` | 403 |
| Not found | `response.NotFound` | 404 |
| Rate limited | `response.TooManyRequests` | 429 |
| Unexpected | `response.InternalError` | 500 |

## Tests

- Unit tests live next to the code (`*_test.go`). Fake repositories implement the domain `Repository` interface.
- Integration tests against PostgreSQL are named `*_integration_test.go` and should skip when `POSTGRES_URL` is not set.
- Run `make test` and `make vet` before finishing any change.
