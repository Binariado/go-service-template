# Pull Request Review Checklist

## Architecture
- [ ] Handler has no business logic
- [ ] SQL exists only in repositories (via sqlc)
- [ ] DTOs stay in transport and never reach the domain
- [ ] Domain imports no transport, infrastructure or HTTP/SQL packages
- [ ] Transport does not import infrastructure
- [ ] Dependencies are injected via constructors and wired in `application.New` / the module `app.go`
- [ ] One use case = one file = one `Execute`

## Transport & API
- [ ] Routes are plural nouns and registered in the module `router.go`
- [ ] Correct status codes (200/201/204/400/401/403/404/429/500)
- [ ] Responses use `response.*` helpers; errors use `ErrorBody`
- [ ] Response DTOs are mapped from entities; timestamps are RFC3339
- [ ] Required fields and UUID params are validated
- [ ] Protected routes use `AuthMiddleware` / `RequirePermission`
- [ ] Lists use `listing.Normalize` (defaults page=1, limit=20, max=100) with a sort allow-list

## Data
- [ ] Migration has matching `.up.sql` and `.down.sql`
- [ ] `make sqlc` was run and queries compile
- [ ] No applied migration was edited

## Quality
- [ ] `ctx` is the first parameter and is propagated
- [ ] Errors are logged once with `slog.*Context`; internal errors are not leaked to clients
- [ ] No `os.Getenv` outside `config/`; new variables are in `.env.example`
- [ ] Tests added/updated; `make vet` and `make test` pass
- [ ] Docs updated when structure or conventions change
