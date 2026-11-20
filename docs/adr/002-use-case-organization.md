# ADR-002: One use case per file

**Status:** Accepted

**Decision:** Every application operation is a struct `<Verb><Noun>UseCase` in `application/<module>/<verb>_<noun>.go` with a single `Execute(ctx, ...)`. Transports never call domain services. When a module grows, group use cases in sub-packages per entity (`application/<module>/<entity>/`, import alias `<entity>uc`) while keeping `application/<module>/app.go` as the module composition root.

**Why:** Open/Closed (adding an operation does not modify existing files), isolated tests, a uniform signature, and a natural home for cross-domain orchestration.

**Trade-offs:** The module `app.go` grows linearly with the number of use cases.
