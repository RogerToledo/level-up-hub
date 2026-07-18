package database

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/me/level-up-hub/backend/config"
)

// RunMigrations executes all pending database migrations.
// It uses the database URL and migrations path from the config.
// Returns nil if migrations succeed or if there are no new migrations to apply.
func RunMigrations(cfg *config.Config) error {
	dbURL := cfg.GetDatabaseURL()

	// golang-migrate pgx5 driver expects "pgx5://" scheme
	// Handle both "postgres://" and "postgresql://" prefixes
	var migrateURL string
	switch {
	case strings.HasPrefix(dbURL, "postgresql://"):
		migrateURL = "pgx5://" + strings.TrimPrefix(dbURL, "postgresql://")
	case strings.HasPrefix(dbURL, "postgres://"):
		migrateURL = "pgx5://" + strings.TrimPrefix(dbURL, "postgres://")
	default:
		migrateURL = dbURL
	}

	// Resolve migrations path to absolute for reliability in containers
	migrationsPath := cfg.MigrationsPath
	if !filepath.IsAbs(migrationsPath) {
		wd, _ := os.Getwd()
		migrationsPath = filepath.Join(wd, migrationsPath)
	}

	slog.Info("running database migrations",
		slog.String("path", migrationsPath),
		slog.String("env", cfg.Env),
	)

	m, err := migrate.New(
		"file://"+migrationsPath,
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
