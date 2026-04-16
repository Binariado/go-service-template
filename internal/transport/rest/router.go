package rest

import (
	"backend-service/internal/application"
	"backend-service/internal/transport/rest/handler"

	"github.com/go-chi/chi/v5"
)

// NewRouter mounts all versioned REST routes under /api/v1.
// It receives the Application so each handler gets only the service it needs.
func NewRouter(app *application.Application) *chi.Mux {
	r := chi.NewRouter()

	// Users
	userHandler := handler.NewUserHandler(app.User)
	r.Get("/users/{id}", userHandler.GetByID)

	return r
}
