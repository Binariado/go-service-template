package user

// User is the domain entity. It has no dependency on any transport layer.
type User struct {
	ID   string
	Name string
}
