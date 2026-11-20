# Backend Service Template · Go

A general-purpose Go backend template for building production-ready services.
REST is first-class. GraphQL is an optional adapter. Authentication, background
jobs and other concerns are documented, so you add only what your service needs.

---

## What's included

- **REST** — versioned router (`/api/v1`), one package per module (router, handlers, DTOs), JSON response helpers
- **GraphQL** — optional adapter (`transport/graphql`, gqlgen) wired through the application layer
- **Health endpoints** — `/health` (liveness) and `/readiness` (DB ping)
- **Use cases** — one file per operation; transport never calls the domain directly
- **Clean domain** — entities with zero dependency on transport or drivers
- **PostgreSQL** — connection pool, **sqlc** typed queries and **golang-migrate** migrations (embedded in the binary, applied on startup)
- **Auth contract** — `TokenValidator`, `AuthMiddleware`, `RequirePermission` (no implementation; JWT example in docs)
- **Structured logs** — `log/slog`, request IDs on every call
- **Production bootstrap** — HTTP timeouts, graceful shutdown (30 s), CORS, panic recovery
- **Tooling** — Makefile, Docker (multi-stage, distroless) and Docker Compose with PostgreSQL
- **Docs & standards** — architecture, coding standards, ADRs, review checklist, `AGENT.md`

---

## Architecture

```
.
├── cmd/api/main.go                 # Bootstrap: config → DB → migrations → app → router → server
├── config/                         # Environment-driven config (single source of truth)
├── docs/                           # Architecture, standards, ADRs, guides
├── internal/
│   ├── domain/<module>/            # entity.go, repository.go, service.go — no transport, no drivers
│   ├── application/
│   │   ├── app.go                  # application.New(db): top-level composition root
│   │   ├── common/listing/         # Pagination / sort normalization
│   │   └── <module>/
│   │       ├── app.go              # Module composition root (lists its use cases)
│   │       └── <verb>_<noun>.go    # One use case per file
│   ├── infrastructure/
│   │   ├── database/               # DB handle, transaction context
│   │   │   └── pg/                 # connect.go, migrate.go, migrations/, queries/, sqlc/ (generated)
│   │   └── repository/             # Repository implementations (the only code that runs SQL)
│   └── transport/
│       ├── rest/
│       │   ├── router.go           # Mounts every module at /api/v1
│       │   ├── handler/            # Health / readiness
│       │   ├── middleware/         # Auth / authorization contract
│       │   ├── response/           # JSON, NotFound, BadRequest, Unauthorized, Forbidden, ...
│       │   └── <module>/           # router.go, handler/, dto/
│       └── graphql/                # Optional gqlgen adapter (schema, resolvers, generated code)
├── gqlgen.yml · sqlc.yaml · Makefile · Dockerfile · docker-compose.yml
└── AGENT.md
```

### Layer rules

| Layer | Allowed imports | Forbidden imports |
|---|---|---|
| `domain` | stdlib only | transport, infrastructure, ORMs |
| `application` | domain, infrastructure (composition root) | transport |
| `infrastructure` | domain, stdlib, drivers | transport |
| `transport` | application | domain services, infrastructure |
| `cmd/*` | all layers | — |

> **Key rule:** handlers receive `*application.Application` and call use cases via
> `app.<Module>.<UseCase>.Execute(...)`. They never call domain services or repositories.

More in [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) and [docs/CODING_STANDARDS.md](docs/CODING_STANDARDS.md).

---

## Getting started

### 1. Clone and rename the module

```bash
git clone <this-repo> my-service
cd my-service

make rename MODULE=github.com/your-org/my-service
```

This replaces the module path in Go files, `go.mod`, YAML, Markdown and the Dockerfile, then runs `go mod tidy`. After renaming you may remove the `rename` target from the `Makefile`. Also rename the `container_name`s in `docker-compose.yml` if you wish.

### 2. Configure environment

```bash
cp .env.example .env
```

| Variable | Required | Description | Example |
|---|---|---|---|
| `PORT` | yes | HTTP port | `3000` |
| `POSTGRES_URL` | yes | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable` |
| `ORIGIN_ALLOWED` | no | Comma-separated CORS origins | `http://localhost:3000,https://myapp.com` |
| `GRAPHQL_PLAYGROUND_ENABLED` | no | Serve the GraphQL playground at `/` (default `false`) | `true` |

### 3. Run

```bash
# With Docker (API + PostgreSQL)
docker compose up --build

# Or locally (needs a reachable PostgreSQL)
make run
```

Migrations are applied on startup. Generated sqlc code is not committed: `make run|build|test|vet` and the Dockerfile generate it for you (`make sqlc`).

---

## Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Liveness — 200 if the process is alive |
| `GET` | `/readiness` | Readiness — 200 only if the DB is reachable |
| `GET` | `/api/v1/users/{id}` | Example REST endpoint (`id` is a UUID) |
| `POST` | `/api/v1/users` | Example write endpoint, body `{"name": "Ada"}` → 201 |
| `POST` | `/query` | GraphQL endpoint (`userById` query, `createUser` mutation) |
| `GET` | `/` | GraphQL playground (only with `GRAPHQL_PLAYGROUND_ENABLED=true`) |

Try it:

```bash
psql "$POSTGRES_URL" -c "INSERT INTO users (name) VALUES ('Ada') RETURNING id"
curl localhost:3000/api/v1/users/<id>
```

---

## How to extend the template

Follow this order (also in [AGENT.md](AGENT.md)):

1. **Migration** — `internal/infrastructure/database/pg/migrations/000002_create_orders_table.{up,down}.sql`
2. **Query** — `internal/infrastructure/database/pg/queries/order.sql`, then `make sqlc`
3. **Domain** — `internal/domain/order/{entity,repository,service}.go`
4. **Repository** — `internal/infrastructure/repository/order.go`
5. **Use case** — `internal/application/order/create_order.go`

```go
type CreateOrderUseCase struct {
    orders *domainorder.Service
    users  *domainuser.Service // cross-domain orchestration belongs here
}

func NewCreateOrderUseCase(orders *domainorder.Service, users *domainuser.Service) *CreateOrderUseCase { ... }

func (uc *CreateOrderUseCase) Execute(ctx context.Context, req CreateOrderRequest) (*domainorder.Order, error) { ... }
```

6. **Module app + root wiring**

```go
// internal/application/order/app.go
type App struct{ CreateOrder *CreateOrderUseCase }

func NewApp(orders *domainorder.Service, users *domainuser.Service) *App {
    return &App{CreateOrder: NewCreateOrderUseCase(orders, users)}
}

// internal/application/app.go
type Application struct {
    User  *user.App
    Order *order.App
}
// in New(db): orderRepo := repository.NewSQLOrderRepository(db.Postgres) ... Order: order.NewApp(orderSvc, userSvc)
```

7. **REST module** — `internal/transport/rest/order/{router.go,handler/,dto/}`

```go
// internal/transport/rest/order/router.go
func RegisterRoutes(r chi.Router, app *application.Application) {
    h := orderhandler.NewOrderHandler(app)
    r.Post("/orders", h.Create)
    r.Get("/orders/{id}", h.GetByID)
}

// internal/transport/rest/router.go
user.RegisterRoutes(r, app)
order.RegisterRoutes(r, app)
```

8. **GraphQL (optional)** — edit `internal/transport/graphql/schema.graphqls`, run `make gqlgen`, then implement the resolver by calling the use case and mapping the entity to the GraphQL model.

### Protecting routes

The template ships only the contract. Implement `middleware.TokenValidator` with the mechanism you prefer and use `AuthMiddleware` / `RequirePermission` — see [docs/AUTH_JWT.md](docs/AUTH_JWT.md).

### Removing GraphQL

Delete `internal/transport/graphql/` and `gqlgen.yml`, remove the `graph.RegisterRoutes(...)` line and its import from `cmd/api/main.go`, drop the `gqlgen` target from the Makefile, and run `go mod tidy`.

### Background jobs

Not included. See [docs/JOBS.md](docs/JOBS.md) for a suggested Asynq + Valkey setup.

---

## Infrastructure

- **PostgreSQL pool** (`pg/connect.go`): `MaxOpenConns 25`, `MaxIdleConns 5`, `ConnMaxLifetime 5 min`.
- **HTTP timeouts** (`cmd/api/main.go`): `ReadTimeout 15 s`, `WriteTimeout 15 s`, `IdleTimeout 60 s`.
- **Graceful shutdown**: `SIGINT` / `SIGTERM` give in-flight requests up to 30 s.
- **Migrations**: embedded SQL, applied on startup — see [docs/DATABASE.md](docs/DATABASE.md).

---

## Commands

```bash
make help        # list targets
make run         # generate sqlc + run the API
make build       # generate sqlc + build ./service
make test        # generate sqlc + go test ./...
make vet         # generate sqlc + go vet ./...
make lint        # golangci-lint
make sqlc        # regenerate typed queries
make gqlgen      # regenerate GraphQL code
make compose-up  # Docker Compose (API + PostgreSQL)
make tidy        # go mod tidy
```

---

## Documentation

| File | Content |
|---|---|
| [AGENT.md](AGENT.md) | Rules and workflow for AI agents / contributors |
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Tree, layer rules, request flow |
| [docs/CODING_STANDARDS.md](docs/CODING_STANDARDS.md) | Naming, aliases, error mapping, tests |
| [docs/DATABASE.md](docs/DATABASE.md) | Migrations, sqlc, transactions |
| [docs/API_GUIDELINES.md](docs/API_GUIDELINES.md) | Routes, responses, pagination |
| [docs/AUTH_JWT.md](docs/AUTH_JWT.md) | Example JWT authentication |
| [docs/JOBS.md](docs/JOBS.md) | Example background jobs |
| [docs/REVIEW_CHECKLIST.md](docs/REVIEW_CHECKLIST.md) | PR checklist |
| [docs/adr/](docs/adr/) | Architecture decision records |

---

## Design decisions

- **REST is first-class, GraphQL is optional.** A service can use only REST without touching any GraphQL file.
- **Transport receives `*application.Application`, nothing else.** The layer boundary is enforced by the compiler.
- **Layer-first, module-second.** Each domain is a module in every layer; adding one never modifies the others.
- **One use case per file.** Cross-domain orchestration lives in use cases — not in handlers, and not in a domain service that imports another.
- **`application.New(db)` is the single composition root.** Each module wires its use cases in its own `app.go`.
- **Domain has zero transport dependency.** Mappers (entity → DTO / GraphQL model) live in the transport layer.
- **Nothing vendor- or business-specific.** No auth provider, no specific middleware; only reusable platform patterns.
