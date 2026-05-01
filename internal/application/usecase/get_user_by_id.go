package usecase

import (
	"backend-service/internal/domain/user"
	"context"
)

// GetUserByID retrieves a user by its identifier.
// Place cross-domain orchestration here when this operation grows to involve
// more than one domain service (e.g. fetching permissions, audit logging).
type GetUserByID struct {
	users *user.Service
}

func NewGetUserByID(users *user.Service) *GetUserByID {
	return &GetUserByID{users: users}
}

func (uc *GetUserByID) Execute(ctx context.Context, id string) (*user.User, error) {
	return uc.users.FindByID(ctx, id)
}