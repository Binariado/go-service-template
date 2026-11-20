package user

import (
	"github.com/binariado/go-service-template/internal/application"
	userhandler "github.com/binariado/go-service-template/internal/transport/rest/user/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers the user module routes on r.
func RegisterRoutes(r chi.Router, app *application.Application) {
	userHandler := userhandler.NewUserHandler(app)

	r.Post("/users", userHandler.Create)
	r.Get("/users/{id}", userHandler.GetByID)
}
