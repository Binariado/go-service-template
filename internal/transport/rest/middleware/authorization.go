package middleware

import (
	"context"
	"net/http"
	"slices"

	"github.com/binariado/go-service-template/internal/transport/rest/response"
)

// Permission constants follow Perm{Module}{Entity}{Action} and the value
// "{module}.{entity}.{action}". Declare yours here, e.g.:
//
//	const PermUserUsersRead = "user.users.read"

// HasPermission reports whether the authenticated caller holds perm.
func HasPermission(ctx context.Context, perm string) bool {
	claims := GetClaims(ctx)
	return claims != nil && slices.Contains(claims.Permissions, perm)
}

// RequirePermission rejects requests whose caller does not hold perm.
// It must run after AuthMiddleware.
//
//	r.With(middleware.RequirePermission(PermUserUsersRead)).Get("/users/{id}", h.GetByID)
func RequirePermission(perm string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if GetClaims(r.Context()) == nil {
				response.Unauthorized(w, "authentication required")
				return
			}
			if !HasPermission(r.Context(), perm) {
				response.Forbidden(w, "missing permission: "+perm)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
