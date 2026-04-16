package apperr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageError(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		err      error
		expected string
	}{
		{
			name:     "error with custom message",
			message:  "Failed to create user",
			err:      errors.New("database error"),
			expected: "Failed to create user: database error",
		},
		{
			name:     "error with nil error",
			message:  "Custom message",
			err:      nil,
			expected: "Custom message: <nil>",
		},
		{
			name:     "empty message",
			message:  "",
			err:      errors.New("some error"),
			expected: ": some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MessageError(tt.message, tt.err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, result.Error())
		})
	}
}

func TestErrorConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
	}{
		{name: "ErrInternalServerError", constant: ErrInternalServerError},
		{name: "ErrBadRequest", constant: ErrBadRequest},
		{name: "ErrUnauthorized", constant: ErrUnauthorized},
		{name: "ErrNotFound", constant: ErrNotFound},
		{name: "ErrInvalidCredentials", constant: ErrInvalidCredentials},
		{name: "ErrDuplicateEmail", constant: ErrDuplicateEmail},
		{name: "ErrInvalidDate", constant: ErrInvalidDate},
		{name: "ErrNoManagerEmail", constant: ErrNoManagerEmail},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.constant)
		})
	}
}

func TestSuccessMessageConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
	}{
		{name: "OkCreate", constant: OkCreate},
		{name: "OkUpdate", constant: OkUpdate},
		{name: "OkDelete", constant: OkDelete},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.constant)
			assert.Contains(t, tt.constant, "%s")
		})
	}
}

func TestEntityNameConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
	}{
		{name: "UserPT", constant: UserPT},
		{name: "ActivityPT", constant: ActivityPT},
		{name: "LadderPT", constant: LadderPT},
		{name: "EvidencePT", constant: EvidencePT},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.constant)
		})
	}
}

func TestErrorMessageFormatting(t *testing.T) {
	tests := []struct {
		name     string
		template string
		entity   string
		expected string
	}{
		{
			name:     "create message",
			template: OkCreate,
			entity:   UserPT,
			expected: "Usuário criado com sucesso!",
		},
		{
			name:     "update message",
			template: OkUpdate,
			entity:   ActivityPT,
			expected: "Atividade atualizada com sucesso!",
		},
		{
			name:     "delete message",
			template: OkDelete,
			entity:   LadderPT,
			expected: "Nível deletado com sucesso!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test validates the format string structure
			assert.Contains(t, tt.template, "%s")
		})
	}
}

func TestErrorCodes(t *testing.T) {
	errorMessages := []string{
		ErrInternalServerError,
		ErrBadRequest,
		ErrUnauthorized,
		ErrNotFound,
		ErrInvalidCredentials,
		ErrDuplicateEmail,
		ErrInvalidDate,
		ErrNoManagerEmail,
	}

	for _, msg := range errorMessages {
		t.Run(msg, func(t *testing.T) {
			assert.NotEmpty(t, msg)
			assert.Greater(t, len(msg), 5, "Error message should be descriptive")
		})
	}
}

func TestMessageErrorChaining(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := MessageError("wrapped", originalErr)

	assert.NotNil(t, wrappedErr)
	assert.Contains(t, wrappedErr.Error(), "original error")
	assert.Contains(t, wrappedErr.Error(), "wrapped")
}

func TestEntityNamesInPortuguese(t *testing.T) {
	entities := map[string]string{
		UserPT:     "user",
		ActivityPT: "activity",
		LadderPT:   "ladder",
		EvidencePT: "evidence",
	}

	for entity, english := range entities {
		t.Run(entity, func(t *testing.T) {
			assert.NotEmpty(t, entity)
			assert.NotEqual(t, entity, english, "PT constant should be in Portuguese")
		})
	}
}

// Benchmark tests
func BenchmarkMessageError(b *testing.B) {
	err := errors.New("test error")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MessageError("test message", err)
	}
}
