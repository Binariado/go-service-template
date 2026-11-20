package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeValidator struct{}

func (fakeValidator) ValidateToken(token string) (*Claims, error) {
	switch token {
	case "admin":
		return &Claims{Subject: "1", Permissions: []string{"x.read"}}, nil
	case "reader":
		return &Claims{Subject: "2"}, nil
	}
	return nil, errors.New("invalid")
}

func TestAuthAndPermissions(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	h := AuthMiddleware(fakeValidator{})(RequirePermission("x.read")(ok))

	tests := []struct {
		name   string
		header string
		want   int
	}{
		{"no header", "", http.StatusUnauthorized},
		{"bad format", "Token abc", http.StatusUnauthorized},
		{"invalid token", "Bearer nope", http.StatusUnauthorized},
		{"missing permission", "Bearer reader", http.StatusForbidden},
		{"allowed", "Bearer admin", http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tt.want {
				t.Errorf("status = %d, want %d", rec.Code, tt.want)
			}
		})
	}
}
