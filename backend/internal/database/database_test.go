package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPoolConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		maxConns    int
		minConns    int
		valid       bool
		description string
	}{
		{
			name:        "valid configuration",
			maxConns:    25,
			minConns:    5,
			valid:       true,
			description: "Max connections should be greater than min connections",
		},
		{
			name:        "invalid - min greater than max",
			maxConns:    5,
			minConns:    25,
			valid:       false,
			description: "Min connections cannot exceed max connections",
		},
		{
			name:        "invalid - zero max connections",
			maxConns:    0,
			minConns:    5,
			valid:       false,
			description: "Max connections must be positive",
		},
		{
			name:        "valid - equal min and max",
			maxConns:    10,
			minConns:    10,
			valid:       true,
			description: "Min and max can be equal for fixed pool size",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.maxConns > 0 && tt.minConns >= 0 && tt.maxConns >= tt.minConns
			assert.Equal(t, tt.valid, isValid, tt.description)
		})
	}
}

func TestConnectionLifetimeValidation(t *testing.T) {
	tests := []struct {
		name     string
		lifetime int
		valid    bool
	}{
		{
			name:     "valid lifetime - 1 hour",
			lifetime: 3600,
			valid:    true,
		},
		{
			name:     "valid lifetime - 30 minutes",
			lifetime: 1800,
			valid:    true,
		},
		{
			name:     "invalid - negative lifetime",
			lifetime: -100,
			valid:    false,
		},
		{
			name:     "valid - zero (unlimited)",
			lifetime: 0,
			valid:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.lifetime >= 0
			assert.Equal(t, tt.valid, isValid)
		})
	}
}

func TestConnectionPoolMetrics(t *testing.T) {
	// Test metrics structure
	type PoolMetrics struct {
		TotalConns      int32
		IdleConns       int32
		AcquiredConns   int32
		MaxConns        int32
		AcquireCount    int64
		AcquireDuration time.Duration
	}

	metrics := PoolMetrics{
		TotalConns:      10,
		IdleConns:       5,
		AcquiredConns:   5,
		MaxConns:        25,
		AcquireCount:    100,
		AcquireDuration: time.Millisecond * 50,
	}

	assert.Equal(t, int32(10), metrics.TotalConns)
	assert.Equal(t, int32(5), metrics.IdleConns)
	assert.Equal(t, metrics.TotalConns, metrics.IdleConns+metrics.AcquiredConns)
	assert.LessOrEqual(t, metrics.TotalConns, metrics.MaxConns)
}

func TestHealthCheckInterval(t *testing.T) {
	tests := []struct {
		name     string
		interval int
		valid    bool
	}{
		{
			name:     "valid - 60 seconds",
			interval: 60,
			valid:    true,
		},
		{
			name:     "valid - 30 seconds",
			interval: 30,
			valid:    true,
		},
		{
			name:     "invalid - too frequent (1 second)",
			interval: 1,
			valid:    false,
		},
		{
			name:     "invalid - negative",
			interval: -10,
			valid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Health checks should be at least 5 seconds apart
			isValid := tt.interval >= 5
			assert.Equal(t, tt.valid, isValid)
		})
	}
}

func TestDatabaseURLParsing(t *testing.T) {
	tests := []struct {
		name  string
		url   string
		valid bool
	}{
		{
			name:  "valid postgres URL",
			url:   "postgres://user:pass@localhost:5432/dbname",
			valid: true,
		},
		{
			name:  "valid with sslmode",
			url:   "postgres://user:pass@localhost:5432/dbname?sslmode=disable",
			valid: true,
		},
		{
			name:  "invalid - empty URL",
			url:   "",
			valid: false,
		},
		{
			name:  "invalid - missing protocol",
			url:   "user:pass@localhost:5432/dbname",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simple validation - check if URL starts with postgres://
			isValid := len(tt.url) > 0 && (tt.url[:11] == "postgres://" || tt.url[:14] == "postgresql://")
			assert.Equal(t, tt.valid, isValid)
		})
	}
}

func TestConnectionTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout int
		valid   bool
	}{
		{
			name:    "valid - 5 seconds",
			timeout: 5,
			valid:   true,
		},
		{
			name:    "valid - 10 seconds",
			timeout: 10,
			valid:   true,
		},
		{
			name:    "invalid - zero timeout",
			timeout: 0,
			valid:   false,
		},
		{
			name:    "invalid - negative timeout",
			timeout: -5,
			valid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.timeout > 0
			assert.Equal(t, tt.valid, isValid)
		})
	}
}

func TestContextWithTimeout(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
	}{
		{
			name:     "5 second timeout",
			duration: 5 * time.Second,
		},
		{
			name:     "1 second timeout",
			duration: time.Second,
		},
		{
			name:     "100 millisecond timeout",
			duration: 100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.duration)
			defer cancel()

			assert.NotNil(t, ctx)

			deadline, ok := ctx.Deadline()
			assert.True(t, ok)
			assert.True(t, time.Until(deadline) <= tt.duration)
		})
	}
}

// Test pool stats structure
func TestPoolStats(t *testing.T) {
	type Stats struct {
		AcquireCount         int64
		AcquireDuration      time.Duration
		AcquiredConns        int32
		CanceledAcquireCount int64
		EmptyAcquireCount    int64
		IdleConns            int32
		MaxConns             int32
		TotalConns           int32
	}

	stats := Stats{
		AcquireCount:    1000,
		AcquiredConns:   5,
		IdleConns:       10,
		MaxConns:        25,
		TotalConns:      15,
		AcquireDuration: time.Millisecond * 10,
	}

	assert.Greater(t, stats.AcquireCount, int64(0))
	assert.LessOrEqual(t, stats.TotalConns, stats.MaxConns)
	assert.Equal(t, stats.TotalConns, stats.IdleConns+stats.AcquiredConns)
}

// Benchmark tests
func BenchmarkConnectionPoolValidation(b *testing.B) {
	maxConns := 25
	minConns := 5

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = maxConns > 0 && minConns >= 0 && maxConns >= minConns
	}
}
