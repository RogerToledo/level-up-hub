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

	// Force migration version if configured (used to recover from dirty state)
	if cfg.ForceMigrationVersion > 0 {
		slog.Warn("forcing migration version", slog.Int("version", cfg.ForceMigrationVersion))
		if err := m.Force(cfg.ForceMigrationVersion); err != nil {
			slog.Error("failed to force migration version", slog.String("error", err.Error()))
			return err
		}
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("migrations: no new migrations to apply")
			return nil
		}
		// Handle dirty state: force the dirty version to clean it, then retry
		if dirtyVersion, ok := isDirtyError(err); ok {
			slog.Warn("migrations: dirty state detected, attempting recovery",
				slog.Int("dirty_version", dirtyVersion),
			)
			// Force the dirty version (marks it as not dirty so Up can proceed)
			if forceErr := m.Force(dirtyVersion); forceErr != nil {
				slog.Error("failed to force migration version", slog.String("error", forceErr.Error()))
				return forceErr
			}
			// Try Up again — it will skip the forced version and apply the next ones
			if retryErr := m.Up(); retryErr != nil && !errors.Is(retryErr, migrate.ErrNoChange) {
				// If it still fails (e.g. "already exists"), the DB is ahead of schema_migrations
				// Try forcing to the latest available migration
				slog.Warn("migrations: retry failed, DB may be ahead of schema_migrations",
					slog.String("error", retryErr.Error()),
				)
				return retryErr
			}
			slog.Info("migrations: recovered from dirty state")
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

// isDirtyError checks if the error is a dirty database error and returns the version
func isDirtyError(err error) (int, bool) {
	if err == nil {
		return 0, false
	}
	// golang-migrate returns ErrDirty with the version number
	var dirtyErr *migrate.ErrDirty
	if errors.As(err, &dirtyErr) {
		return dirtyErr.Version, true
	}
	return 0, false
}
