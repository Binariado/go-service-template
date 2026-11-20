package application

import (
	"github.com/binariado/go-service-template/internal/application/user"
	domainuser "github.com/binariado/go-service-template/internal/domain/user"
	"github.com/binariado/go-service-template/internal/infrastructure/database"
	"github.com/binariado/go-service-template/internal/infrastructure/repository"
)

// Application is the top-level composition root.
// It is the single dependency injected into transport layers (REST, GraphQL, workers).
// Each field is a module App; add new modules here as the project grows.
type Application struct {
	User *user.App
}

// New wires all dependencies and returns a ready Application.
func New(db *database.DB) *Application {
	userRepo := repository.NewSQLUserRepository(db.Postgres)
	userService := domainuser.NewUserService(userRepo)

	return &Application{
		User: user.NewApp(userService),
	}
}
