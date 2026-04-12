package identity

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ParseID(id string) (uuid.UUID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, err
	}
	return parsedID, nil
}

func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return uuid.UUID{}, errors.New("user not found in token")
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("invalid user_id type")
	}

	return userID, nil
}

func ValidateIDParam(c *gin.Context) (uuid.UUID, error) {
	idString := c.Param("id")
	if idString == "" {
		return uuid.UUID{}, errors.New("id is required")
	}

	id, err := ParseID(idString)
	if err != nil {
		return uuid.UUID{}, errors.New("invalid id format")
	}

	return id, nil	
}