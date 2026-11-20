package pg

import (
	"database/sql"
	"errors"
	"log/slog"
	"os"

	"github.com/binariado/go-service-template/internal/infrastructure/database/pg/migrations"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// RunMigrations applies all pending migrations. Migration files are embedded in
// the binary, so no files need to be shipped with the runtime image.
func RunMigrations(db *sql.DB) {
	source, err := iofs.New(migrations.FS, ".")
	if err != nil {
		slog.Error("failed to load embedded migrations", "error", err)
		os.Exit(1)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error("failed to create migration driver", "error", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		slog.Error("failed to create migrator", "error", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	slog.Info("database migrations applied successfully")
}
