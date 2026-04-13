package api

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/me/level-up-hub/backend/internal/logger"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		requestLogger := logger.WithRequestID(slog.Default(), requestID)

		if userID, exists := c.Get("user_id"); exists {
			if uid, ok := userID.(uuid.UUID); ok {
				requestLogger = logger.WithUserID(requestLogger, uid.String())
			}
		}

		c.Set("logger", requestLogger)

		requestLogger.Info("incoming request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("query", c.Request.URL.RawQuery),
			slog.String("ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
		)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logLevel := slog.LevelInfo
		if status >= 500 {
			logLevel = slog.LevelError
		} else if status >= 400 {
			logLevel = slog.LevelWarn
		}

		requestLogger.Log(c.Request.Context(), logLevel, "request completed",
			slog.Int("status", status),
			slog.Duration("latency", latency),
			slog.Int("response_size", c.Writer.Size()),
		)
	}
}
