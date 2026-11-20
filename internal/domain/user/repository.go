package user

import "context"

// Repository defines the persistence contract for the user domain.
// Implementations live in the infrastructure layer.
type Repository interface {
	// FindByID returns nil, nil when the user does not exist.
	FindByID(ctx context.Context, id string) (*User, error)

	// Create persists a new user and returns it with its generated ID.
	Create(ctx context.Context, name string) (*User, error)
}
