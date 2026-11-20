# AGENT.md — go-service-template

Guidance for AI agents and contributors working in this template.

## Principles

- Think before coding; understand the request first.
- Simplicity first: minimum code that solves the problem, no speculative features.
- Surgical changes; do not refactor unrelated code.
- Follow existing patterns. This is a **template**: keep it generic — no business-specific, vendor-specific or auth-provider-specific code.

## Architecture (summary)

Layer-first, module-second: `transport → application → domain`, `infrastructure` implements domain contracts. Full details in [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) and [docs/CODING_STANDARDS.md](docs/CODING_STANDARDS.md).

Hard rules:

- Transport never imports infrastructure; it receives `*application.Application`.
- Domain imports no transport, drivers or frameworks.
- Only repositories execute SQL. DTOs never enter the domain.
- One use case = one file = one `Execute(ctx, ...)`.
- Dependencies are injected through constructors; `application.New(db)` is the composition root.
- Never call `os.Getenv` outside `config/`.

## Adding a feature (in order)

1. Migration: `internal/infrastructure/database/pg/migrations/000NNN_name.{up,down}.sql`
2. Query: `internal/infrastructure/database/pg/queries/<module>.sql`, then `make sqlc`
3. Domain: `internal/domain/<module>/{entity,repository,service}.go`
4. Repository: `internal/infrastructure/repository/<module>.go`
5. Use case: `internal/application/<module>/<verb>_<noun>.go`; register it in `application/<module>/app.go`
6. Wire the module in `internal/application/app.go`
7. REST: `internal/transport/rest/<module>/{router.go,handler/,dto/}`; add `<module>.RegisterRoutes(r, app)` in `transport/rest/router.go`
8. GraphQL (optional): update `schema.graphqls`, `make gqlgen`, implement the resolver

## Commands

| Command | Action |
|---|---|
| `make run` | sqlc + run the API |
| `make build` | sqlc + build the binary |
| `make test` / `make vet` | sqlc + tests / static analysis |
| `make lint` | golangci-lint (CI runs it too) |
| `make sqlc` | Regenerate sqlc code |
| `make gqlgen` | Regenerate GraphQL code |
| `make compose-up` | API + PostgreSQL with Docker Compose |

Run `make vet` and `make test` before considering work complete.

## Documentation

| File | Content |
|---|---|
| `docs/ARCHITECTURE.md` | Tree, layer rules, request flow |
| `docs/CODING_STANDARDS.md` | Naming, aliases, rules, error mapping |
| `docs/DATABASE.md` | Migrations, sqlc, transactions |
| `docs/API_GUIDELINES.md` | Routes, responses, pagination |
| `docs/AUTH_JWT.md` | Example JWT authentication |
| `docs/JOBS.md` | Example background jobs (Asynq + Valkey) |
| `docs/REVIEW_CHECKLIST.md` | PR checklist |
| `docs/adr/` | Architecture decision records |
