package identity

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUserIDFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns user ID when present and valid UUID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		expectedID := uuid.New()
		c.Set("user_id", expectedID)

		userID, err := GetUserIDFromContext(c)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, userID)
	})

	t.Run("returns error when user_id not in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		_, err := GetUserIDFromContext(c)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found in token")
	})

	t.Run("returns error when user_id is wrong type", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Set("user_id", "not-a-uuid")

		_, err := GetUserIDFromContext(c)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user_id type")
	})

	t.Run("returns error when user_id is nil", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Set("user_id", nil)

		_, err := GetUserIDFromContext(c)

		assert.Error(t, err)
	})
}

func TestValidateIDParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns valid UUID from path parameter", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		expectedID := uuid.New()
		c.Params = gin.Params{
			{Key: "id", Value: expectedID.String()},
		}

		id, err := ValidateIDParam(c)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
	})

	t.Run("returns error for invalid UUID format", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{
			{Key: "id", Value: "not-a-valid-uuid"},
		}

		_, err := ValidateIDParam(c)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid id format")
	})

	t.Run("returns error for empty ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{
			{Key: "id", Value: ""},
		}

		_, err := ValidateIDParam(c)

		assert.Error(t, err)
	})

	t.Run("returns error when ID param not present", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{}

		_, err := ValidateIDParam(c)

		assert.Error(t, err)
	})
}
