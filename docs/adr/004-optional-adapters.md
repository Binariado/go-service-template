# ADR-004: REST first; GraphQL, jobs and auth are optional

**Status:** Accepted

**Decision:** REST (`transport/rest`) is the primary transport. GraphQL lives in `transport/graphql` as an adapter wired by a single `graph.RegisterRoutes` call and can be removed by deleting the folder. Background jobs are documented ([JOBS.md](../JOBS.md)), not included. Authentication ships only as a contract (`TokenValidator`, `AuthMiddleware`, `RequirePermission`); an implementation example is in [AUTH_JWT.md](../AUTH_JWT.md).

**Why:** The template must be a clean base: no vendor-specific or business-specific code, and nothing to rip out for teams that do not need it.
