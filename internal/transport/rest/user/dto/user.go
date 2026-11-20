package userdto

// UserResponse is the REST output representation of a user.
// It is intentionally separate from the domain entity so the API contract
// can evolve independently from the domain model.
type UserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateUserRequest is the REST input for creating a user.
type CreateUserRequest struct {
	Name string `json:"name"`
}
