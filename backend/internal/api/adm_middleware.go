package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/me/level-up-hub/backend/apperr"
	"github.com/me/level-up-hub/backend/internal/rest"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			rest.Error(c.Writer, http.StatusForbidden, apperr.ErrAdminOnly, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
