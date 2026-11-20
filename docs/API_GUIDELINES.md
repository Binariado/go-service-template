# API Guidelines

## Routes

- Versioned under `/api/v1`; plural nouns; ids in the path: `/users/{id}`.
- Each module registers its routes in `transport/rest/<module>/router.go` (`RegisterRoutes`), and `transport/rest/router.go` only mounts modules.
- Validate UUID path/query params with `uuid.Parse` and answer `400` when invalid.
- Group protected routes and apply the auth contract:

```go
r.Group(func(r chi.Router) {
    r.Use(middleware.AuthMiddleware(validator))
    r.With(middleware.RequirePermission(PermUserUsersRead)).Get("/users/{id}", h.GetByID)
})
```

## Write endpoints

`POST /api/v1/users` is the reference implementation (`transport/rest/user/handler/user.go`):

- Decode with `json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))`; answer `400` on malformed JSON.
- Validation lives in the domain service, which returns a sentinel error (`user.ErrInvalidName`); the handler maps it with `errors.Is` to `400`. Unknown errors are logged and answered with `500`.
- Answer `201` with the created resource as a response DTO.
- The same use case backs the GraphQL `createUser` mutation, so rules are never duplicated per transport.

## Status codes

`200` get/update · `201` create · `204` delete · `400` `401` `403` `404` `409` `429` `500` errors.

## Responses

- Always use `response.*` helpers, never `http.Error`.
- Errors use the envelope `{"error": "not_found", "code": "optional.stable_code", "message": "..."}`.
- Response DTOs are separate from domain entities and mapped explicitly.

## Lists

Query params: `page` (default 1), `limit` (default 20, max 100), `sort`, `order` (`asc|desc`), plus module filters (`?search=`, `?is_active=true`).
Use `application/common/listing.Normalize` inside the use case; it validates sort fields against an allow-list (safe to use in `ORDER BY`).

Suggested envelope: `{ "data": [...], "total": 0, "page": 1, "limit": 20 }`.

## Soft delete

Add `deleted_at TIMESTAMPTZ`, filter `WHERE deleted_at IS NULL` in list queries and return `204` from delete.

## GraphQL

Optional adapter at `/query`. Resolvers call use cases and map entities to GraphQL models explicitly. The playground at `/` is enabled only with `GRAPHQL_PLAYGROUND_ENABLED=true`. Regenerate with `make gqlgen`.
