# Background jobs (example)

The template has **no** job system; this is a suggested approach for services that need one, modeled on a production setup using [Asynq](https://github.com/hibiken/asynq) with Redis/Valkey.

## Layout

```
cmd/worker/main.go                       # Second entrypoint: config → DB → app → job server
internal/infrastructure/jobs/
├── client.go                            # Enqueues tasks (used by the API)
├── server.go                            # Runs the worker (asynq.Server)
└── registry.go                          # Maps task type → handler, rejects duplicates
internal/application/<module>/jobs.go    # Task type constants + payload structs + handler logic
```

Rules:

- Job handlers are thin: decode payload → call a use case. Business logic stays in the application layer, so the same use case serves REST, GraphQL and workers.
- The use case depends on a small **port** (interface) like `JobEnqueuer`, implemented by `infrastructure/jobs`. Application code never imports Asynq.
- Handlers must be **idempotent** (tasks can be retried). Return `asynq.SkipRetry`-wrapped errors for permanent failures.
- `JOBS_URL` unset → API runs without a client (enqueue becomes a no-op or an error, your choice); the worker requires it.

## Dependencies

```bash
go get github.com/hibiken/asynq
```

## Client and server sketch

```go
// internal/infrastructure/jobs/client.go
type Client struct{ c *asynq.Client }

func NewClient(redisURL string) (*Client, error) {
    opt, err := asynq.ParseRedisURI(redisURL)
    if err != nil { return nil, err }
    return &Client{c: asynq.NewClient(opt)}, nil
}

func (c *Client) Enqueue(ctx context.Context, taskType string, payload []byte, opts ...asynq.Option) error {
    _, err := c.c.EnqueueContext(ctx, asynq.NewTask(taskType, payload), opts...)
    return err
}

func (c *Client) Close() error { return c.c.Close() }
```

```go
// internal/infrastructure/jobs/server.go
func NewServer(redisURL string, concurrency int, mux *asynq.ServeMux) (*asynq.Server, error) {
    opt, err := asynq.ParseRedisURI(redisURL)
    if err != nil { return nil, err }
    return asynq.NewServer(opt, asynq.Config{
        Concurrency:     concurrency,
        ShutdownTimeout: 30 * time.Second,
        ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, t *asynq.Task, err error) {
            slog.Error("job failed", "task_type", t.Type(), "error", err)
        }),
    }), nil
}
```

`cmd/worker/main.go` builds the same `application.New(db)`, registers handlers in the registry, starts the server and shuts it down on `SIGINT`/`SIGTERM`.

## Configuration

Add to `config/config.go`: `JOBS_URL` (e.g. `redis://valkey:6379/1`) and `JOBS_CONCURRENCY` (default 10, must be > 0). Provide a `LoadWorker()` variant that does not require `PORT`.

## Docker Compose

```yaml
  valkey:
    image: valkey/valkey:8-alpine

  worker:
    build: .
    command: ["/app/worker"]      # build ./cmd/worker into /app/worker in the Dockerfile
    env_file: .env
    environment:
      JOBS_URL: redis://valkey:6379/1
    depends_on: [valkey, db]
```

Add `build-worker` / `run-worker` targets to the Makefile (`go build -o worker ./cmd/worker`).
