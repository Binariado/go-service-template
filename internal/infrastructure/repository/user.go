package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/binariado/go-service-template/internal/domain/user"
	"github.com/binariado/go-service-template/internal/infrastructure/database/pg/sqlc"

	"github.com/google/uuid"
)

// SQLUserRepository implements user.Repository using PostgreSQL through sqlc.
type SQLUserRepository struct {
	queries *sqlc.Queries
}

func NewSQLUserRepository(db *sql.DB) *SQLUserRepository {
	return &SQLUserRepository{queries: sqlc.New(db)}
}

func (r *SQLUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		// A malformed identifier can never match a stored user.
		return nil, nil
	}

	row, err := r.queries.GetUserByID(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user.User{ID: row.ID.String(), Name: row.Name}, nil
}

func (r *SQLUserRepository) Create(ctx context.Context, name string) (*user.User, error) {
	row, err := r.queries.CreateUser(ctx, name)
	if err != nil {
		return nil, err
	}

	return &user.User{ID: row.ID.String(), Name: row.Name}, nil
}
