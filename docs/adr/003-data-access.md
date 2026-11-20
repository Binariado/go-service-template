# ADR-003: PostgreSQL with sqlc and embedded migrations

**Status:** Accepted

**Decision:** PostgreSQL through `database/sql` + `lib/pq`. Schema is versioned with golang-migrate; the SQL files are embedded in the binary and applied on startup. Queries are written in SQL and compiled to typed Go by sqlc (generated code is gitignored and produced by `make sqlc`). Raw `QueryContext`/`ExecContext` is the fallback for dynamic queries. Repositories are the only place that executes SQL.

**Why:** Type-safe queries without an ORM, reproducible schema, and a self-contained binary/image.

**Trade-offs:** Generated code must exist before building (`make` targets and the Dockerfile run `make sqlc`); applying migrations on startup assumes a single deployment unit or a migration-safe rollout.
