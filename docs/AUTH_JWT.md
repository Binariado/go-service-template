# Authentication with JWT (example)

The template does **not** ship an authentication system. It ships a contract in `internal/transport/rest/middleware`:

```go
type Claims struct { Subject string; Permissions []string }
type TokenValidator interface { ValidateToken(token string) (*Claims, error) }

AuthMiddleware(validator)        // 401 unless "Authorization: Bearer <valid token>"
RequirePermission("module.entity.action") // 403 unless Claims.Permissions has it
```

Any mechanism (JWT, opaque sessions, an external IdP) fits by implementing `TokenValidator`. Below is one way to do it with JWT.

## 1. Dependency

```bash
go get github.com/golang-jwt/jwt/v5
```

## 2. Validator

`internal/transport/rest/middleware/jwt.go` (token validation is an HTTP concern, so it lives next to the middleware and can use `Claims` without breaking layer rules):

```go
package middleware

import (
    "errors"

    "github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
    Permissions []string `json:"permissions"`
    jwt.RegisteredClaims
}

type JWTValidator struct {
    key      any    // []byte for HS256, *rsa.PublicKey for RS256
    method   string // "HS256" or "RS256"
    issuer   string
    audience string
}

func NewJWTValidator(key any, method, issuer, audience string) *JWTValidator {
    return &JWTValidator{key: key, method: method, issuer: issuer, audience: audience}
}

func (v *JWTValidator) ValidateToken(token string) (*Claims, error) {
    parsed, err := jwt.ParseWithClaims(token, &jwtClaims{},
        func(*jwt.Token) (any, error) { return v.key, nil },
        jwt.WithValidMethods([]string{v.method}), // never accept the "none" algorithm or algorithm switching
        jwt.WithIssuer(v.issuer),
        jwt.WithAudience(v.audience),
        jwt.WithExpirationRequired(),
    )
    if err != nil {
        return nil, err
    }
    c, ok := parsed.Claims.(*jwtClaims)
    if !ok || c.Subject == "" {
        return nil, errors.New("invalid claims")
    }
    return &Claims{Subject: c.Subject, Permissions: c.Permissions}, nil
}
```

## 3. Configuration

Add the values to `config/config.go` (never read `os.Getenv` elsewhere), e.g. `JWT_SECRET` or `JWT_PUBLIC_KEY_PATH`, `JWT_ISSUER`, `JWT_AUDIENCE`, `ACCESS_TOKEN_EXPIRY_MINUTES`. Fail fast at startup when they are missing.

## 4. Wiring

Build the validator in `cmd/api/main.go`, pass it to `rest.NewRouter(app, validator)`, forward it to the module `RegisterRoutes`, and protect routes:

```go
r.Group(func(r chi.Router) {
    r.Use(middleware.AuthMiddleware(validator))
    r.With(middleware.RequirePermission(PermUserUsersRead)).Get("/users/{id}", h.GetByID)
})
```

## 5. Issuing tokens (sign-in)

Model it as a use case (`application/auth/signin.go`) that verifies credentials through a domain `auth` module and asks a token issuer for tokens. The issuer is an infrastructure component (`infrastructure/security/jwt_issuer.go`) that signs with the private key; the HTTP handler only maps the request/response.

## Recommendations

- **Prefer RS256** (private key signs, public key validates and can be exposed through `/.well-known/jwks.json`) when other services must validate tokens. HS256 is fine for a single service.
- Short-lived access tokens (≈15 min) plus opaque, hashed, revocable **refresh tokens** stored in the database.
- Hash passwords with bcrypt/argon2 (`golang.org/x/crypto`).
- Keep keys out of the repo (`config/keys/` is a good gitignored location) and never log tokens.
- Add rate limiting to sign-in and refresh endpoints.
