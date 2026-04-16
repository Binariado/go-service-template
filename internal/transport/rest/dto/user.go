package dto

// UserResponse is the REST output representation of a user.
// It is intentionally separate from the domain entity so the API contract
// can evolve independently from the domain model.
type UserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
