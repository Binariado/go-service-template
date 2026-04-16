package user

// Repository defines the persistence contract for the user domain.
// Implementations live in the infrastructure layer.
type Repository interface {
	FindByID(id string) (*User, error)
}
