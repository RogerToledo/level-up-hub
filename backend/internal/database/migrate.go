package database

import (
	"errors"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/me/level-up-hub/backend/config"
)

// RunMigrations executes all pending database migrations.
// It uses the database URL and migrations path from the config.
// Returns nil if migrations succeed or if there are no new migrations to apply.
func RunMigrations(cfg *config.Config) error {
	var dbURL string
	if cfg.Env == "prod" {
		dbURL = cfg.DbURLProd
	} else {
		dbURL = cfg.DbURLDev
	}

	// golang-migrate pgx5 driver expects "pgx5://" scheme
	migrateURL := "pgx5" + dbURL[len("postgres"):]

	slog.Info("running database migrations",
		slog.String("path", cfg.MigrationsPath),
		slog.String("env", cfg.Env),
	)

	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		migrateURL,
	)
	if err != nil {
		slog.Error("failed to create migrate instance", slog.String("error", err.Error()))
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("migrations: no new migrations to apply")
			return nil
		}
		slog.Error("failed to run migrations", slog.String("error", err.Error()))
		return err
	}

	version, dirty, _ := m.Version()
	slog.Info("migrations applied successfully",
		slog.Uint64("version", uint64(version)),
		slog.Bool("dirty", dirty),
	)

	return nil
}
