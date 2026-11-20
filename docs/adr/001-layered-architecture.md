# ADR-001: Layered architecture, layer-first and module-second

**Status:** Accepted

**Decision:** Code is organized first by layer (`domain`, `application`, `infrastructure`, `transport`) and then by module (`user`, `order`, ...). Dependencies point inward: transport → application → domain; infrastructure implements domain contracts. Each module has three domain files (`entity.go`, `repository.go`, `service.go`), its own `app.go` in `application/<module>`, and its own `transport/rest/<module>` package with `router.go`, `handler/` and `dto/`.

**Why:** Boundaries are enforced by the compiler (transport only sees `*application.Application`), modules can be added without touching existing ones, and it scales from one domain to many without restructuring.

**Trade-offs:** More files and packages than a flat layout; accepted for predictability.
