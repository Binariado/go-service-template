// Package middleware provides reusable HTTP middlewares for the service.
//
// To add authentication middleware, implement it here following the chi middleware
// pattern and register it in server.go or in the specific router that needs it.
//
// Example skeleton:
//
//	func RequireAuth(next http.Handler) http.Handler {
//	    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	        token := r.Header.Get("Authorization")
//	        // validate token ...
//	        next.ServeHTTP(w, r)
//	    })
//	}
package middleware
