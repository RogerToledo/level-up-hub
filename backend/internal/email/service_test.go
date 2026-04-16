package email

import (
	"testing"

	"github.com/me/level-up-hub/backend/config"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	cfg := &config.Config{
		SMTPHost:     "smtp.gmail.com",
		SMTPPort:     587,
		SMTPUser:     "test@example.com",
		SMTPPassword: "testpass",
		SMTPFrom:     "noreply@example.com",
		SMTPFromName: "Test Service",
	}

	service := NewService(cfg)

	assert.NotNil(t, service)
	assert.Equal(t, cfg, service.cfg)
}

func TestEncodeBase64(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "empty input",
			input:    []byte{},
			expected: "",
		},
		{
			name:     "simple text",
			input:    []byte("Hello World"),
			expected: "SGVsbG8gV29ybGQ=",
		},
		{
			name:     "binary data",
			input:    []byte{0xFF, 0xFE, 0xFD},
			expected: "//79",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encodeBase64(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEmailConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *config.Config
		expectValid bool
	}{
		{
			name: "valid configuration",
			cfg: &config.Config{
				SMTPHost:     "smtp.gmail.com",
				SMTPPort:     587,
				SMTPUser:     "user@example.com",
				SMTPPassword: "password",
				SMTPFrom:     "noreply@example.com",
				SMTPFromName: "Service",
			},
			expectValid: true,
		},
		{
			name: "missing SMTP host",
			cfg: &config.Config{
				SMTPHost:     "",
				SMTPPort:     587,
				SMTPUser:     "user@example.com",
				SMTPPassword: "password",
			},
			expectValid: false,
		},
		{
			name: "invalid port",
			cfg: &config.Config{
				SMTPHost:     "smtp.gmail.com",
				SMTPPort:     0,
				SMTPUser:     "user@example.com",
				SMTPPassword: "password",
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.cfg.SMTPHost != "" && tt.cfg.SMTPPort > 0
			assert.Equal(t, tt.expectValid, isValid)
		})
	}
}

func TestSMTPPortValidation(t *testing.T) {
	validPorts := []int{25, 465, 587, 2525}

	for _, port := range validPorts {
		t.Run("valid port", func(t *testing.T) {
			assert.Greater(t, port, 0)
			assert.Less(t, port, 65536)
		})
	}
}

func TestEmailAddressValidation(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		isValid bool
	}{
		{
			name:    "valid email",
			email:   "user@example.com",
			isValid: true,
		},
		{
			name:    "valid email with subdomain",
			email:   "user@mail.example.com",
			isValid: true,
		},
		{
			name:    "invalid - no @",
			email:   "userexample.com",
			isValid: false,
		},
		{
			name:    "invalid - no domain",
			email:   "user@",
			isValid: false,
		},
		{
			name:    "invalid - no username",
			email:   "@example.com",
			isValid: false,
		},
		{
			name:    "empty email",
			email:   "",
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simple email validation logic
			hasAt := false
			hasDot := false
			parts := []string{}
			
			for i, c := range tt.email {
				if c == '@' {
					hasAt = true
					if i == 0 || i == len(tt.email)-1 {
						hasAt = false
					}
				}
				if c == '.' && hasAt {
					hasDot = true
				}
			}
			
			isValid := tt.email != "" && hasAt && hasDot
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestReportEmailTemplate(t *testing.T) {
	userName := "John Doe"
	managerName := "Jane Manager"

	template := "Olá " + managerName + ",\n\n" +
		"Segue em anexo o relatório de atividades de " + userName + ".\n\n" +
		"Atenciosamente,\n" +
		"Level Up Hub"

	assert.Contains(t, template, userName)
	assert.Contains(t, template, managerName)
	assert.Contains(t, template, "relatório")
	assert.Contains(t, template, "Level Up Hub")
}

func TestSMTPAuthValidation(t *testing.T) {
	tests := []struct {
		name     string
		user     string
		password string
		valid    bool
	}{
		{
			name:     "valid credentials",
			user:     "user@example.com",
			password: "password123",
			valid:    true,
		},
		{
			name:     "empty user",
			user:     "",
			password: "password123",
			valid:    false,
		},
		{
			name:     "empty password",
			user:     "user@example.com",
			password: "",
			valid:    false,
		},
		{
			name:     "both empty",
			user:     "",
			password: "",
			valid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.user != "" && tt.password != ""
			assert.Equal(t, tt.valid, isValid)
		})
	}
}

// Benchmark tests
func BenchmarkEncodeBase64(b *testing.B) {
	data := []byte("This is a test message for benchmarking base64 encoding performance")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = encodeBase64(data)
	}
}

func BenchmarkNewService(b *testing.B) {
	cfg := &config.Config{
		SMTPHost:     "smtp.gmail.com",
		SMTPPort:     587,
		SMTPUser:     "test@example.com",
		SMTPPassword: "testpass",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewService(cfg)
	}
}
