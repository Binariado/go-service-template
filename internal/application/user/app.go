package user

import (
	domainuser "github.com/binariado/go-service-template/internal/domain/user"
)

// App is the composition root of the user module.
// Add one field per use case and wire it in NewApp.
type App struct {
	GetUserByID *GetUserByIDUseCase
	CreateUser  *CreateUserUseCase
}

func NewApp(svc *domainuser.Service) *App {
	return &App{
		GetUserByID: NewGetUserByIDUseCase(svc),
		CreateUser:  NewCreateUserUseCase(svc),
	}
}
