package rest
package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	tests := []struct {
		name           string
		message        interface{}
		expectedStatus int
	}{
		{
			name:           "success with string message",
			message:        "Operation successful",
			expectedStatus: http.StatusOK,
		},
		{
			name: "success with struct message",
			message: map[string]interface{}{
				"id":   "123",
				"name": "Test",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "created status",
			message:        "Resource created",
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			Send(w, tt.message, tt.expectedStatus)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		})
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		errorMsg       string
		details        interface{}
		expectedStatus int
	}{
		{
			name:           "bad request error",
			status:         http.StatusBadRequest,
			errorMsg:       "Invalid input",
			details:        "Field 'name' is required",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "unauthorized error",
			status:         http.StatusUnauthorized,
			errorMsg:       "Unauthorized access",
			details:        nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "internal server error",
			status:         http.StatusInternalServerError,
			errorMsg:       "Internal error",
			details:        errors.New("database connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "not found error",
			status:         http.StatusNotFound,
			errorMsg:       "Resource not found",
			details:        "User with id 123 not found",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			Error(w, tt.status, tt.errorMsg, tt.details)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		})
	}
}

func TestStatusCodeMapping(t *testing.T) {
	statusCodes := map[string]int{
		"OK":                    http.StatusOK,
		"Created":               http.StatusCreated,
		"Bad Request":           http.StatusBadRequest,
		"Unauthorized":          http.StatusUnauthorized,
		"Forbidden":             http.StatusForbidden,
		"Not Found":             http.StatusNotFound,
		"Internal Server Error": http.StatusInternalServerError,
	}

	for name, code := range statusCodes {
		t.Run(name, func(t *testing.T) {
			assert.GreaterOrEqual(t, code, 200)
			assert.LessOrEqual(t, code, 599)
		})
	}
}

func TestJSONEncoding(t *testing.T) {
	tests := []struct {
		name    string
		message interface{}
	}{
		{
			name:    "string message",
			message: "test",
		},
		{
			name:    "number message",
			message: 123,
		},
		{
			name: "map message",
			message: map[string]string{
				"key": "value",
			},
		},
		{
			name: "nested structure",
			message: map[string]interface{}{
				"user": map[string]string{
					"name":  "John",
					"email": "john@example.com",
				},
				"status": "active",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			Send(w, tt.message, http.StatusOK)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.NotEmpty(t, w.Body.String())
		})
	}
}

// Test error response structure
func TestErrorResponseStructure(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, http.StatusBadRequest, "Validation error", "Invalid field")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
	assert.Contains(t, w.Body.String(), "Validation error")
}

// Benchmark tests
func BenchmarkSend(b *testing.B) {
	message := map[string]string{
		"message": "Success",
		"status":  "ok",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		Send(w, message, http.StatusOK)
	}
}

func BenchmarkError(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		Error(w, http.StatusBadRequest, "Error message", "Error details")
	}
}
