# Backend Service Template · Go

A general-purpose Go backend template for building production-ready services.
REST is first-class. GraphQL is an optional adapter. Workers and webhooks can be
added without rewriting the foundation.

---

## What's included

- **REST** — versioned router (`/api/v1`), handlers, DTOs, JSON response helpers
- **GraphQL** — optional adapter using `gqlgen`, wired through the application layer
- **Health endpoints** — `/health` (liveness) and `/readiness` (DB ping)
- **Use cases** — application layer with explicit use case structs; transport never calls domain directly
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
│   ├── application/                # Composition root + use cases
│   │   ├── app.go                  # application.New(db) → exposes use cases to transports
│   │   └── usecase/
│   │       └── get_user_by_id.go   # One file per use case
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
| `transport` | application | domain internals, infrastructure |
| `server.go` | all layers | nothing specific from domain internals |

> **Key rule:** handlers receive `*application.Application` and call use cases via
> `app.<UseCase>.Execute(...)`. They never import domain packages or call domain
> services directly.

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

Create the domain files:

```
internal/domain/order/
    entity.go        # Order struct
    repository.go    # Repository interface
    service.go       # Business logic
internal/infrastructure/repository/order.go   # SQL implementation
```

Add a use case in `internal/application/usecase/`:

```go
// internal/application/usecase/create_order.go
type CreateOrder struct {
    orders *order.Service
    users  *user.Service   // cross-domain orchestration belongs here
}

func NewCreateOrder(orders *order.Service, users *user.Service) *CreateOrder {
    return &CreateOrder{orders: orders, users: users}
}

func (uc *CreateOrder) Execute(ctx context.Context, cmd CreateOrderCmd) error {
    // orchestrate across domains without coupling them to each other
}
```

Register everything in `internal/application/app.go`:

```go
type Application struct {
    GetUserByID *usecase.GetUserByID
    CreateOrder *usecase.CreateOrder   // add here
}

func New(db *sql.DB) *Application {
    userRepo  := repository.NewSQLUserRepository(db)
    orderRepo := repository.NewSQLOrderRepository(db)

    userService  := user.NewUserService(userRepo)
    orderService := order.NewOrderService(orderRepo)

    return &Application{
        GetUserByID: usecase.NewGetUserByID(userService),
        CreateOrder: usecase.NewCreateOrder(orderService, userService),
    }
}
```

---

### Add a REST handler

Create the handler and register the route:

```go
// internal/transport/rest/handler/order.go
type OrderHandler struct {
    app *application.Application
}

func NewOrderHandler(app *application.Application) *OrderHandler {
    return &OrderHandler{app: app}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
    // decode DTO → call use case → encode response
    err := h.app.CreateOrder.Execute(r.Context(), cmd)
    ...
}
```

Register the route in `internal/transport/rest/router.go`:

```go
orderHandler := handler.NewOrderHandler(app)
r.Post("/orders", orderHandler.Create)
r.Get("/orders/{id}", orderHandler.GetByID)
```

---

### Scaling REST by domain (recommended as the service grows)

When the number of domains grows, split `transport/rest/` into one sub-package
per domain instead of keeping all handlers and DTOs in a flat folder:

```
internal/transport/rest/
├── router.go
├── response/
│   └── json.go
├── user/                        # ← domain sub-package
│   ├── handler.go               # UserHandler
│   └── dto.go                   # UserResponse, CreateUserRequest, ...
└── order/                       # ← domain sub-package
    ├── handler.go               # OrderHandler
    └── dto.go                   # OrderResponse, CreateOrderRequest, ...
```

Each domain package registers its own routes. `router.go` only orchestrates:

```go
// internal/transport/rest/router.go
func NewRouter(app *application.Application) *chi.Mux {
    r := chi.NewRouter()

    userHandler  := user.NewHandler(app)
    orderHandler := order.NewHandler(app)

    r.Get("/users/{id}", userHandler.GetByID)
    r.Post("/orders",    orderHandler.Create)
    r.Get("/orders/{id}", orderHandler.GetByID)

    return r
}
```

Handlers across all sub-packages follow the same rule: they receive
`*application.Application` and call use cases — never domain services directly.

---

### Scaling use cases by domain (recommended as the service grows)

When use cases multiply, group them into domain sub-packages under `application/`:

```
internal/application/
├── app.go
└── usecase/
    ├── user/
    │   ├── get_by_id.go
    │   └── create_user.go
    └── order/
        ├── create_order.go
        └── cancel_order.go
```

`app.go` stays the single composition root — it just imports from sub-packages:

```go
type Application struct {
    GetUserByID *useruc.GetByID
    CreateUser  *useruc.CreateUser
    CreateOrder *orderuc.CreateOrder
    CancelOrder *orderuc.CancelOrder
}
```

---

### Add a GraphQL resolver

1. Update `internal/infrastructure/api/graph/schema.graphqls`
2. Regenerate:

```bash
go run github.com/99designs/gqlgen generate
```

3. Implement the resolver in `schema.resolvers.go` — call the use case, map to GraphQL model:

```go
func (r *queryResolver) OrderByID(ctx context.Context, id string) (*model.Order, error) {
    o, err := r.App.GetOrderByID.Execute(ctx, id)
    if err != nil {
        return nil, err
    }
    return &model.Order{ID: o.ID}, nil  // explicit mapper: domain entity → transport model
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
- **Transport receives `*application.Application`, nothing else.** Handlers call
  `app.<UseCase>.Execute(...)` and never import domain or infrastructure packages.
  This enforces the layer boundary at the compiler level.
- **Use cases live in `application/usecase/`.** Single-domain operations start there.
  When an operation needs to coordinate multiple domain services, the use case is
  the right place — not the handler, and not a domain service that imports another.
- **`application.New(db)` is the single composition root.** Adding a new dependency
  means changing one file (`app.go`), not every handler.
- **Domain has zero transport dependency.** Entities and service interfaces never
  import `model.User` (GraphQL) or REST DTOs. Mappers live in the transport layer.
- **Infrastructure specifics stay out of the template.** No Firebase, no specific
  auth provider, no business-specific middleware. Only reusable platform patterns.
