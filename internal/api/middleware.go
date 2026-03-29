package api

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		log.Printf("[%s] %d %s %s", c.Request.Method, c.Writer.Status(), c.Request.URL.Path, latency)
	}
}