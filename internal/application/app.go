package application

import (
	"backend-service/internal/application/usecase"
	"backend-service/internal/domain/user"
	"backend-service/internal/infrastructure/repository"
	"database/sql"
)

// Application is the composition root for all use cases.
// It is the single dependency injected into transport layers (REST, GraphQL, workers).
// Add new use cases here as the project grows.
type Application struct {
	GetUserByID *usecase.GetUserByID
}

// New wires all dependencies and returns a ready Application.
func New(db *sql.DB) *Application {
	userRepo := repository.NewSQLUserRepository(db)
	userService := user.NewUserService(userRepo)

	return &Application{
		GetUserByID: usecase.NewGetUserByID(userService),
	}
}