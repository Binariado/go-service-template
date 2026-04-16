package repository

import (
	"backend-service/internal/domain/user"
	"database/sql"
)

// SQLUserRepository implements user.Repository using PostgreSQL.
type SQLUserRepository struct {
	db *sql.DB
}

func NewSQLUserRepository(db *sql.DB) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) FindByID(userID string) (*user.User, error) {
	// TODO: replace with a real SQL query
	// Example: SELECT id, name FROM users WHERE id = $1
	if userID == "1" {
		return &user.User{ID: "1", Name: "Template User"}, nil
	}
	return nil, nil
}
