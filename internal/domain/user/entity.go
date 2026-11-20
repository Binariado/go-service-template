package user

import "errors"

// ErrInvalidName is returned when a user name is empty or too long.
var ErrInvalidName = errors.New("name must be between 1 and 100 characters")

// User is the domain entity. It has no dependency on any transport layer.
type User struct {
	ID   string
	Name string
}
