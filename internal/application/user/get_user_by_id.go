package user

import (
	"context"

	domainuser "github.com/binariado/go-service-template/internal/domain/user"
)

// GetUserByIDUseCase retrieves a user by its identifier.
// Place cross-domain orchestration here when this operation grows to involve
// more than one domain service (e.g. fetching permissions, audit logging).
type GetUserByIDUseCase struct {
	users *domainuser.Service
}

func NewGetUserByIDUseCase(users *domainuser.Service) *GetUserByIDUseCase {
	return &GetUserByIDUseCase{users: users}
}

// Execute returns nil, nil when the user does not exist.
func (uc *GetUserByIDUseCase) Execute(ctx context.Context, id string) (*domainuser.User, error) {
	return uc.users.FindByID(ctx, id)
}
