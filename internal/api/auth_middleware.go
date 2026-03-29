package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/me/level-up-hub/apperr"
	"github.com/me/level-up-hub/auth"
	"github.com/me/level-up-hub/internal/rest"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			rest.Error(c.Writer, 401, apperr.ErrRequiredToken, nil)
			c.Abort()

			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			rest.Error(c.Writer, 401, apperr.ErrInvalidToken, nil)
			c.Abort()

			return
		}

		tokenString := parts[1]

		claims, err := auth.ValidateToken(tokenString, secret)
		if err != nil {
			rest.Error(c.Writer, 401, apperr.ErrInvalidToken, nil)
			c.Abort()
			
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()

	}
}
