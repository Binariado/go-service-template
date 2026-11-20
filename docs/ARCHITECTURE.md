# Architecture

Clean Architecture, **layer-first, module-second**. Dependencies flow inward only.

```
Transport  (REST handlers/DTOs/routes, GraphQL resolvers)  → Application
Application (use cases, orchestration)                     → Domain
Domain     (entities, repository contracts, rules)         → nothing
Infrastructure (PostgreSQL, sqlc, migrations, ...)         → implements Domain contracts
```

## Tree

```
cmd/
└── api/main.go                       # Bootstrap: config → DB → migrations → app → router → server
config/                               # Environment-driven config (single source of truth)
docs/                                 # This documentation
internal/
├── domain/<module>/                  # entity.go, repository.go, service.go (+ errors.go when needed)
├── application/
│   ├── app.go                        # Top-level composition root: application.New(db)
│   ├── common/listing/               # Shared pagination/sort normalization
│   └── <module>/
│       ├── app.go                    # Module composition root (NewApp)
│       └── <verb>_<noun>.go          # One use case per file
├── infrastructure/
│   ├── database/
│   │   ├── database.go               # DB handle (Postgres, + future stores)
│   │   ├── transaction_context.go    # WithTransaction / TransactionFromContext
│   │   └── pg/
│   │       ├── connect.go            # Connection pool
│   │       ├── migrate.go            # RunMigrations (embedded files)
│   │       ├── migrations/           # 000NNN_name.{up,down}.sql
│   │       ├── queries/              # <module>.sql annotated for sqlc
│   │       └── sqlc/                 # Generated (gitignored) — run `make sqlc`
│   └── repository/                   # Repository implementations (the only place that runs SQL)
└── transport/
    ├── rest/
    │   ├── router.go                 # Mounts every module at /api/v1
    │   ├── handler/                  # Cross-cutting handlers (health, readiness)
    │   ├── middleware/               # Auth/authorization contract
    │   ├── response/                 # JSON helpers
    │   └── <module>/
    │       ├── router.go             # RegisterRoutes(r, app)
    │       ├── handler/              # HTTP handlers
    │       └── dto/                  # Request/response structs
    └── graphql/                      # Optional adapter (gqlgen). Delete the folder to drop it.
```

## Layer rules

| Layer | May import | Must not import |
|---|---|---|
| `domain` | stdlib | transport, infrastructure, drivers/ORMs |
| `application` | domain, infrastructure (composition root only) | transport |
| `infrastructure` | domain, stdlib, drivers | transport |
| `transport` | application | domain internals*, infrastructure |
| `cmd/*` | all layers | — |

\* Transport may read domain entities returned by use cases to map them into DTOs / GraphQL models; it must never call domain services or repositories.

## Request flow

`HTTP → chi router → module router → handler → use case (Execute) → domain service → repository → sqlc → PostgreSQL`

## Composition

- `application.New(db)` is the only place that wires repositories → services → modules.
- Each module has its own `app.go` that lists its use cases. Adding a use case touches that file only.
- Transports receive `*application.Application` and call `app.<Module>.<UseCase>.Execute(...)`.

## Optional pieces

| Piece | Status | How to remove |
|---|---|---|
| GraphQL (`transport/graphql`, `gqlgen.yml`) | Included | Delete both, remove `graph.RegisterRoutes` from `cmd/api/main.go`, run `go mod tidy` |
| Background jobs | Documented only — see [JOBS.md](JOBS.md) | — |
| Authentication | Contract only — see [AUTH_JWT.md](AUTH_JWT.md) | — |
