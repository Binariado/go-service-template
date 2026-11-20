package user

import (
	"context"
	"strings"
	"unicode/utf8"
)

const maxNameLength = 100

// Service handles business logic related to users.
// Mutations and complex orchestration belong here.
// Simple queries are delegated through the application layer.
type Service struct {
	repo Repository
}

// NewUserService creates a new instance of Service.
func NewUserService(repo Repository) *Service {
	return &Service{repo: repo}
}

// FindByID returns a user by its ID, or nil when it does not exist.
func (s *Service) FindByID(ctx context.Context, userID string) (*User, error) {
	return s.repo.FindByID(ctx, userID)
}

// Create validates the name and persists a new user.
func (s *Service) Create(ctx context.Context, name string) (*User, error) {
	name = strings.TrimSpace(name)
	if n := utf8.RuneCountInString(name); n == 0 || n > maxNameLength {
		return nil, ErrInvalidName
	}
	return s.repo.Create(ctx, name)
}
