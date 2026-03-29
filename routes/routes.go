package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/me/level-up-hub/config"
	"github.com/me/level-up-hub/internal/account"
	"github.com/me/level-up-hub/internal/api"
)

type RouterConfig struct {
	UserHandler *account.Handler
}

func NewRouter(cfg RouterConfig, dbPool *pgxpool.Pool, appCfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(
		api.LoggerMiddleware(),
		gin.Recovery(),
	)

	// Health check endpoint to verify database connectivity
	r.GET("/health", func(c *gin.Context) {
		err := dbPool.Ping(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down", "database": "unreachable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "up", "database": "ok"})
	})

	// API v1 routes

	// Account routes
	v1 := r.Group("/v1")
	
	// Public routes
	v1.POST("/login", cfg.UserHandler.Login)
	v1.POST("/register", cfg.UserHandler.Register)

	// Protected routes
	protected := v1.Group("/")
	protected.Use(api.AuthMiddleware(appCfg.JWTSecret))
	protected.PUT("/users/:id", cfg.UserHandler.Update)
	protected.GET("/users/:id", cfg.UserHandler.FindByID)

	// Admin-only routes
	admin := protected.Group("/")
	admin.Use(api.AdminOnly())
	admin.DELETE("/users/:id", cfg.UserHandler.Delete)
	admin.GET("/users", cfg.UserHandler.FindAll)

	return r
}
