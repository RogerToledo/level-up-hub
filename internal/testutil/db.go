package testutil

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/me/level-up-hub/config"
	"github.com/me/level-up-hub/internal/database"
)

// SetupTestDB creates a test database connection
// Only use this for integration tests
func SetupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	cfg := &config.Config{
		Env:       "test",
		DbURLDev:  "postgres://postgres:postgres@localhost:5432/leveluphub_test?sslmode=disable",
		MaxConns:  5,
		MinConns:  1,
		MaxConnLifetime: 3600,
		MaxConnIdleTime: 1800,
		HealthCheckPeriod: 60,
		ConnectTimeout: 5,
	}

	pool, err := database.NewPostgresPool(context.Background(), cfg)
	if err != nil {
		t.Skipf("Skipping test: database not available: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

// CleanupTestData removes test data from database
func CleanupTestData(t *testing.T, pool *pgxpool.Pool, tables ...string) {
	t.Helper()

	ctx := context.Background()
	for _, table := range tables {
		_, err := pool.Exec(ctx, "TRUNCATE TABLE "+table+" CASCADE")
		if err != nil {
			t.Logf("Warning: failed to cleanup table %s: %v", table, err)
		}
	}
}
