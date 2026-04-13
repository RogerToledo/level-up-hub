package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PoolStats contains connection pool statistics
type PoolStats struct {
	TotalConns int32 `json:"total_conns"`
	IdleConns  int32 `json:"idle_conns"`
	MaxConns   int32 `json:"max_conns"`
}

// GetPoolStats returns current pool statistics
func GetPoolStats(pool *pgxpool.Pool) PoolStats {
	stats := pool.Stat()
	return PoolStats{
		TotalConns: stats.TotalConns(),
		IdleConns:  stats.IdleConns(),
		MaxConns:   stats.MaxConns(),
	}
}

// LogPoolStats logs pool statistics
func LogPoolStats(pool *pgxpool.Pool) {
	stats := GetPoolStats(pool)

	slog.Info("connection pool stats",
		slog.Int("total_conns", int(stats.TotalConns)),
		slog.Int("idle_conns", int(stats.IdleConns)),
		slog.Int("max_conns", int(stats.MaxConns)),
		slog.Float64("usage_percent", float64(stats.TotalConns)/float64(stats.MaxConns)*100),
	)

	// Alert if pool is near the limit
	usagePercent := float64(stats.TotalConns) / float64(stats.MaxConns) * 100
	if usagePercent > 80 {
		slog.Warn("connection pool usage high",
			slog.Float64("usage_percent", usagePercent),
			slog.Int("total_conns", int(stats.TotalConns)),
			slog.Int("max_conns", int(stats.MaxConns)),
			slog.String("recommendation", "consider increasing MAX_CONNS"),
		)
	}
}

// StartPoolMonitor starts periodic pool monitoring
// Returns a channel to stop monitoring
func StartPoolMonitor(pool *pgxpool.Pool, interval time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				LogPoolStats(pool)
			case <-stop:
				slog.Info("stopping pool monitor")
				return
			}
		}
	}()

	return stop
}

// HealthCheck verifies pool connection health
func HealthCheck(ctx context.Context, pool *pgxpool.Pool) error {
	// 5 second timeout for health check
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("database health check failed", slog.String("error", err.Error()))
		return err
	}

	stats := GetPoolStats(pool)

	// Check if there are available connections
	if stats.TotalConns == 0 {
		slog.Warn("no database connections available")
	}

	slog.Debug("database health check passed",
		slog.Int("total_conns", int(stats.TotalConns)),
		slog.Int("idle_conns", int(stats.IdleConns)),
	)

	return nil
}
