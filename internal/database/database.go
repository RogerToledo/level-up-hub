package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/me/level-up-hub/config"
)

func NewPostgresPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	var dbUrl string
	if cfg.Env == "prod" {
		dbUrl = cfg.DbUrlProd
	} else {
		dbUrl = cfg.DbUrlDev
	}

	slog.Info("configuring database connection pool",
		slog.String("env", cfg.Env),
		slog.Int("max_conns", cfg.MaxConns),
		slog.Int("min_conns", cfg.MinConns),
		slog.Int("max_conn_lifetime_sec", cfg.MaxConnLifetime),
		slog.Int("max_conn_idle_time_sec", cfg.MaxConnIdleTime),
	)

	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		slog.Error("failed to parse database URL", slog.String("error", err.Error()))
		return nil, err
	}

	// Pool configuration
	poolConfig.MaxConns = int32(cfg.MaxConns)
	poolConfig.MinConns = int32(cfg.MinConns)
	poolConfig.MaxConnLifetime = time.Duration(cfg.MaxConnLifetime) * time.Second
	poolConfig.MaxConnIdleTime = time.Duration(cfg.MaxConnIdleTime) * time.Second
	poolConfig.HealthCheckPeriod = time.Duration(cfg.HealthCheckPeriod) * time.Second

	// Connection timeout
	poolConfig.ConnConfig.ConnectTimeout = time.Duration(cfg.ConnectTimeout) * time.Second

	slog.Debug("attempting to connect to database")

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		slog.Error("failed to create connection pool",
			slog.String("error", err.Error()),
			slog.String("env", cfg.Env),
		)
		return nil, err
	}

	// Test connection
	slog.Debug("pinging database")
	if err := pool.Ping(ctx); err != nil {
		slog.Error("database ping failed", slog.String("error", err.Error()))
		pool.Close()
		return nil, err
	}

	// Success log with statistics
	stats := pool.Stat()
	slog.Info("database connection pool established",
		slog.Int("total_conns", int(stats.TotalConns())),
		slog.Int("idle_conns", int(stats.IdleConns())),
	)

	return pool, nil
}
