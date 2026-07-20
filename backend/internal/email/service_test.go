package email

import (
	"testing"

	"github.com/me/level-up-hub/backend/config"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	cfg := &config.Config{
		SMTPFrom:     "noreply@example.com",
		SMTPFromName: "Test Service",
	}

	service := NewService(cfg)

	assert.NotNil(t, service)
	assert.Equal(t, cfg, service.cfg)
	assert.Nil(t, service.client) // No API key = no client
}

func TestNewServiceWithAPIKey(t *testing.T) {
	cfg := &config.Config{
		ResendAPIKey: "re_test_123",
		SMTPFrom:     "noreply@example.com",
		SMTPFromName: "Test Service",
	}

	service := NewService(cfg)

	assert.NotNil(t, service)
	assert.NotNil(t, service.client)
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"John Doe", "John_Doe"},
		{"user@email.com", "useremailcom"},
		{"José Silva", "Jos_Silva"},
		{"simple", "simple"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeFilename(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
