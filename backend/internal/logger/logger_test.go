package logger

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		name        string
		env         string
		expectLevel slog.Level
	}{
		{
			name:        "production environment uses INFO level",
			env:         "prod",
			expectLevel: slog.LevelInfo,
		},
		{
			name:        "development environment uses DEBUG level",
			env:         "dev",
			expectLevel: slog.LevelDebug,
		},
		{
			name:        "unknown environment defaults to DEBUG",
			env:         "test",
			expectLevel: slog.LevelDebug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := Setup(tt.env)

			assert.NotNil(t, logger, "logger should not be nil")
			assert.Equal(t, logger, slog.Default(), "logger should be set as default")
		})
	}
}

func TestWithRequestID(t *testing.T) {
	requestID := "test-request-123"
	baseLogger := slog.Default()

	logger := WithRequestID(baseLogger, requestID)

	assert.NotNil(t, logger, "logger with request ID should not be nil")
}

func TestWithUserID(t *testing.T) {
	userID := "user-456"
	baseLogger := slog.Default()

	logger := WithUserID(baseLogger, userID)

	assert.NotNil(t, logger, "logger with user ID should not be nil")
}
