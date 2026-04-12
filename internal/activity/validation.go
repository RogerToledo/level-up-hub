package activity

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ValidateActivityID(c *gin.Context) (uuid.UUID, error) {
	id := c.Param("id")
	if id == "" {
		return uuid.UUID{}, errors.New("id is required")
	}
	return uuid.Parse(id)
}
