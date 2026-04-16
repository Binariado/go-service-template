package user

import "context"

// Service handles business logic related to users.
// Mutations and complex orchestration belong here.
// Simple queries are delegated through the application layer.
type Service struct {
	repo Repository
}

// NewUserService creates a new instance of UserService.
func NewUserService(repo Repository) *Service {
	return &Service{repo: repo}
}

// FindByID returns a user by its ID.
func (s *Service) FindByID(ctx context.Context, userID string) (*User, error) {
	return s.repo.FindByID(userID)
}
