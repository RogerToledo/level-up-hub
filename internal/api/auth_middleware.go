package api

import (
	"log/slog"
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
			slog.Warn("unauthorized request - missing token",
				slog.String("path", c.Request.URL.Path),
				slog.String("ip", c.ClientIP()),
			)
			rest.Error(c.Writer, 401, apperr.ErrRequiredToken, nil)
			c.Abort()

			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			slog.Warn("unauthorized request - invalid token format",
				slog.String("path", c.Request.URL.Path),
				slog.String("ip", c.ClientIP()),
			)
			rest.Error(c.Writer, 401, apperr.ErrInvalidToken, nil)
			c.Abort()

			return
		}

		tokenString := parts[1]

		claims, err := auth.ValidateToken(tokenString, secret)
		if err != nil {
			slog.Warn("unauthorized request - invalid token",
				slog.String("path", c.Request.URL.Path),
				slog.String("ip", c.ClientIP()),
			)
			rest.Error(c.Writer, 401, apperr.ErrInvalidToken, nil)
			c.Abort()

			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()

	}
}
