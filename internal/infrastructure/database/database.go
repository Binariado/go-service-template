package database

import (
	"context"
	"database/sql"

	"github.com/binariado/go-service-template/internal/infrastructure/database/pg"
)

// DB groups the data stores used by the service. Add new stores here
// (e.g. a cache client) so the rest of the code receives a single handle.
type DB struct {
	Postgres *sql.DB
}

// Connect opens every configured data store.
func Connect(ctx context.Context, postgresURL string) *DB {
	return &DB{Postgres: pg.Connect(ctx, postgresURL)}
}
