package application

import (
	"backend-service/internal/domain/user"
	"backend-service/internal/infrastructure/repository"
	"database/sql"
)

// Application is the composition root for all domain services and use cases.
// It is the single dependency injected into transport layers (REST, GraphQL, workers).
// Add new domain services here as the project grows.
type Application struct {
	User *user.Service
}

// New wires all dependencies and returns a ready Application.
// db is the only external dependency at this layer; infrastructure adapters are
// instantiated internally so transport layers never know about them.
func New(db *sql.DB) *Application {
	userRepo := repository.NewSQLUserRepository(db)
	userService := user.NewUserService(userRepo)

	return &Application{
		User: userService,
	}
}
