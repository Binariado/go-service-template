package database

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// ConnectPg opens a PostgreSQL connection, verifies reachability and configures
// the connection pool. The caller is responsible for closing the returned *sql.DB.
func ConnectPg(ctx context.Context) *sql.DB {
	dataSource := os.Getenv("POSTGRES_URL")
	if dataSource == "" {
		panic("required environment variable POSTGRES_URL is not set")
	}

	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		slog.Error("failed to open database connection", "error", err)
		os.Exit(1)
	}

	// Connection pool tuning — adjust for production workloads.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		slog.Error("failed to reach database", "error", err)
		os.Exit(1)
	}

	slog.Info("database connection established")
	return db
}
