// Package middleware provides reusable HTTP middlewares for the REST transport.
//
// The template ships only the authentication/authorization contract; there is
// no token implementation. Provide a TokenValidator (JWT, sessions, an external
// provider, ...) and register AuthMiddleware on the routes that need it.
// See docs/AUTH_JWT.md for a JWT example.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/binariado/go-service-template/internal/transport/rest/response"
)

type contextKey string

const claimsKey contextKey = "auth_claims"

// Claims is the identity of the authenticated caller.
// Extend it with whatever your authentication system provides.
type Claims struct {
	Subject     string
	Permissions []string
}

// TokenValidator is the contract an authentication implementation must satisfy.
type TokenValidator interface {
	ValidateToken(token string) (*Claims, error)
}

// AuthMiddleware validates the "Authorization: Bearer <token>" header and
// stores the resulting Claims in the request context.
func AuthMiddleware(validator TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				response.Unauthorized(w, "missing authorization header")
				return
			}

			scheme, token, ok := strings.Cut(header, " ")
			if !ok || scheme != "Bearer" || token == "" {
				response.Unauthorized(w, "invalid authorization header format")
				return
			}

			claims, err := validator.ValidateToken(token)
			if err != nil || claims == nil {
				response.Unauthorized(w, "invalid or expired token")
				return
			}

			next.ServeHTTP(w, r.WithContext(WithClaims(r.Context(), claims)))
		})
	}
}

// WithClaims returns a context carrying claims.
func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// GetClaims returns the claims stored by AuthMiddleware, or nil when absent.
func GetClaims(ctx context.Context) *Claims {
	claims, _ := ctx.Value(claimsKey).(*Claims)
	return claims
}
