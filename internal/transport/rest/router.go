package rest

import (
	"github.com/binariado/go-service-template/internal/application"
	"github.com/binariado/go-service-template/internal/transport/rest/user"

	"github.com/go-chi/chi/v5"
)

// NewRouter mounts every module's routes. The result is mounted at /api/v1.
// Each module exposes RegisterRoutes; add new modules here.
func NewRouter(app *application.Application) *chi.Mux {
	r := chi.NewRouter()

	user.RegisterRoutes(r, app)

	return r
}
