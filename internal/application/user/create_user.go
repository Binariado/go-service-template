package user

import (
	"context"

	domainuser "github.com/binariado/go-service-template/internal/domain/user"
)

// CreateUserUseCase registers a new user.
// Place cross-domain orchestration here (e.g. assigning a default role,
// sending a welcome notification) as this operation grows.
type CreateUserUseCase struct {
	users *domainuser.Service
}

func NewCreateUserUseCase(users *domainuser.Service) *CreateUserUseCase {
	return &CreateUserUseCase{users: users}
}

// Execute returns domainuser.ErrInvalidName when the name is not valid.
func (uc *CreateUserUseCase) Execute(ctx context.Context, name string) (*domainuser.User, error) {
	return uc.users.Create(ctx, name)
}
