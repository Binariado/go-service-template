# Backend Service Template · Go

A general-purpose Go backend template for building production-ready services.
REST is first-class. GraphQL is an optional adapter. Workers and webhooks can be
added without rewriting the foundation.

---

## What's included

- **REST** — versioned router (`/api/v1`), handlers, DTOs, JSON response helpers
- **GraphQL** — optional adapter using `gqlgen`, wired through the application layer
- **Health endpoints** — `/health` (liveness) and `/readiness` (DB ping)
- **Application layer** — single composition root; transport layers never touch repositories directly
- **Clean domain** — domain entities have zero dependency on transport or ORMs
- **Structured logs** — `log/slog` with request IDs on every call
- **Production bootstrap** — HTTP timeouts, graceful shutdown (30 s), CORS, panic recovery
- **PostgreSQL** — connection pool configured out of the box
- **Docker** — multi-stage build with a minimal distroless runtime image

---

## Architecture

```
.
├── config/                         # Environment-driven config (single source of truth)
├── internal/
│   ├── domain/                     # Business rules — no transport, no ORM imports
│   │   └── user/
│   │       ├── entity.go           # Domain entity
│   │       ├── repository.go       # Repository interface (contract)
│   │       └── service.go          # Business logic
│   │
│   ├── application/                # Composition root — wires repos + services
│   │   └── app.go                  # application.New(db) → injected into transports
│   │
│   ├── infrastructure/             # Implementations of domain contracts
│   │   ├── database/               # PostgreSQL connection + pool
│   │   ├── repository/             # SQL implementations of domain repositories
│   │   ├── middleware/             # HTTP middleware skeletons (auth, etc.)
│   │   └── api/graph/              # gqlgen generated code + resolvers
│   │       └── model/              # GraphQL transport models (generated)
│   │
│   └── transport/
│       └── rest/                   # REST transport (first-class)
│           ├── router.go           # Versioned router mounted at /api/v1
│           ├── handler/            # HTTP handlers (one file per domain)
│           ├── dto/                # Request / response structs
│           └── response/           # JSON helpers: JSON, NotFound, BadRequest, InternalError
│
└── server.go                       # Bootstrap: config → DB → app → router → HTTP server
```

### Layer rules

| Layer | Allowed imports | Forbidden imports |
|---|---|---|
| `domain` | stdlib only | transport, infrastructure, ORMs |
| `application` | domain, infrastructure | transport |
| `infrastructure` | domain, stdlib, drivers | transport |
| `transport` | application, dto | domain internals, infrastructure |
| `server.go` | all layers | nothing specific from domain internals |

---

## Getting started

### 1. Clone and rename the module

```bash
git clone <this-repo> my-service
cd my-service

# Replace the module name everywhere
find . -type f -name "*.go" | xargs sed -i '' 's|backend-service|github.com/your-org/my-service|g'
sed -i '' 's|backend-service|github.com/your-org/my-service|g' go.mod
```

### 2. Configure environment

```bash
cp .env.example .env   # or create .env manually
```

Required variables:

| Variable | Description | Example |
|---|---|---|
| `PORT` | Port the HTTP server listens on | `3000` |
| `POSTGRES_URL` | PostgreSQL connection string | `postgres://user:pass@localhost:5432/mydb?sslmode=disable` |
| `ORIGIN_ALLOWED` | Comma-separated allowed CORS origins | `http://localhost:3000,https://myapp.com` |

### 3. Run

```bash
# Development
go run server.go

# With Docker / Podman
docker compose up --build
```

---

## Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Liveness — returns 200 if the process is alive |
| `GET` | `/readiness` | Readiness — returns 200 only if the DB is reachable |
| `GET` | `/api/v1/users/:id` | Example REST endpoint |
| `POST` | `/query` | GraphQL endpoint |
| `GET` | `/` | GraphQL Playground (development only) |

---

## How to extend the template

### Add a new domain

```
internal/domain/order/
    entity.go        # Order struct
    repository.go    # Repository interface
    service.go       # Business logic
internal/infrastructure/repository/order.go   # SQL implementation
```

Register the new service in `internal/application/app.go`:

```go
type Application struct {
    User  *user.Service
    Order *order.Service   // add here
}

func New(db *sql.DB) *Application {
    orderRepo := repository.NewSQLOrderRepository(db)
    return &Application{
        User:  user.NewUserService(repository.NewSQLUserRepository(db)),
        Order: order.NewOrderService(orderRepo),
    }
}
```

### Add a REST handler

1. Create `internal/transport/rest/handler/order.go`
2. Register the route in `internal/transport/rest/router.go`:

```go
orderHandler := handler.NewOrderHandler(app.Order)
r.Get("/orders/{id}", orderHandler.GetByID)
r.Post("/orders", orderHandler.Create)
```

### Add a GraphQL resolver

1. Update `internal/infrastructure/api/graph/schema.graphqls`
2. Regenerate:

```bash
go run github.com/99designs/gqlgen generate
```

3. Implement the resolver in `schema.resolvers.go` — map domain entity → GraphQL model:

```go
func (r *queryResolver) OrderByID(ctx context.Context, id string) (*model.Order, error) {
    domainOrder, err := r.App.Order.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return &model.Order{ID: domainOrder.ID}, nil  // explicit mapper
}
```

---

## Infrastructure

### PostgreSQL connection pool

Configured in `internal/infrastructure/database/connect.go`:

```
MaxOpenConns:    25
MaxIdleConns:     5
ConnMaxLifetime:  5 min
```

Adjust these values based on your workload and database plan.

### HTTP timeouts

Configured in `server.go`:

```
ReadTimeout:   15 s
WriteTimeout:  15 s
IdleTimeout:   60 s
```

### Graceful shutdown

The server listens for `SIGINT` / `SIGTERM` and gives in-flight requests up to
30 seconds to complete before exiting.

---

## Commands

```bash
# Run the service
go run server.go

# Tidy dependencies
go mod tidy

# Regenerate GraphQL code after schema changes
go run github.com/99designs/gqlgen generate

# Build binary
go build -o service .

# Docker
docker compose up --build
docker compose down
```

---

## Design decisions

- **REST is first-class, GraphQL is optional.** A new service can be built using
  only REST without touching any GraphQL file.
- **Domain has zero transport dependency.** Entities and service interfaces never
  import `model.User` (GraphQL) or REST DTOs. Mappers live in the transport layer.
- **`application.New(db)` is the single composition root.** Transport layers receive
  the `Application` struct and nothing else. Adding a new dependency means changing
  one file (`app.go`), not every handler.
- **Infrastructure specifics stay out of the template.** No Firebase, no specific
  auth provider, no business-specific middleware. Only reusable platform patterns.
