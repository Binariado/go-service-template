package handler

import (
	"database/sql"
	"net/http"

	"github.com/binariado/go-service-template/internal/transport/rest/response"
)

// Health handles GET /health — always returns 200 if the process is alive.
func Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Readiness handles GET /readiness — returns 200 only if the DB is reachable.
func Readiness(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			response.JSON(w, http.StatusServiceUnavailable, map[string]string{
				"status": "unavailable",
				"reason": "database unreachable",
			})
			return
		}
		response.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
	}
}
