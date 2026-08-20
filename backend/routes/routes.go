package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/me/level-up-hub/backend/config"
	"github.com/me/level-up-hub/backend/internal/account"
	"github.com/me/level-up-hub/backend/internal/ai"
	"github.com/me/level-up-hub/backend/internal/api"
	"github.com/me/level-up-hub/backend/internal/database"
	"github.com/me/level-up-hub/backend/internal/initiative"
	"github.com/me/level-up-hub/backend/internal/ladder"
	"github.com/me/level-up-hub/backend/internal/leveltarget"
	"github.com/me/level-up-hub/backend/internal/task"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterConfig struct {
	UserHandler        *account.Handler
	LadderHandler      *ladder.LadderHandler
	InitiativeHandler  *initiative.InitiativeHandler
	TaskHandler        *task.TaskHandler
	LevelTargetHandler *leveltarget.Handler
	AIHandler          *ai.Handler
}

func NewRouter(cfg RouterConfig, dbPool *pgxpool.Pool, appCfg *config.Config) *gin.Engine {
	r := gin.New()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     appCfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(
		api.LoggerMiddleware(),
		gin.Recovery(),
	)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		if err := database.HealthCheck(c.Request.Context(), dbPool); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "down",
				"database": "unreachable",
				"error":    err.Error(),
			})
			return
		}

		stats := database.GetPoolStats(dbPool)
		c.JSON(http.StatusOK, gin.H{
			"status":   "up",
			"database": "ok",
			"pool_stats": gin.H{
				"total_conns": stats.TotalConns,
				"idle_conns":  stats.IdleConns,
				"max_conns":   stats.MaxConns,
			},
		})
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := r.Group("/v1")

	// Public routes
	v1.POST("/login", cfg.UserHandler.Login)
	v1.POST("/register", cfg.UserHandler.Register)

	// Protected routes (any authenticated user)
	protected := v1.Group("/")
	protected.Use(api.AuthMiddleware(appCfg.JWTSecret))
	protected.POST("/logout", cfg.UserHandler.Logout)
	protected.GET("/users/:id", cfg.UserHandler.FindByID)
	protected.PUT("/users/:id", cfg.UserHandler.UpdateOwnProfile)

	// Initiatives
	protected.POST("/initiatives", cfg.InitiativeHandler.Create)
	protected.GET("/initiatives", cfg.InitiativeHandler.List)
	protected.PUT("/initiatives/:id", cfg.InitiativeHandler.Update)
	protected.DELETE("/initiatives/:id", cfg.InitiativeHandler.Delete)

	// Tasks (within initiatives)
	protected.POST("/tasks", cfg.TaskHandler.Create)
	protected.PUT("/tasks/:id", cfg.TaskHandler.Update)
	protected.DELETE("/tasks/:id", cfg.TaskHandler.Delete)
	protected.GET("/initiatives/:id/tasks", cfg.TaskHandler.ListByInitiative)
	protected.GET("/tasks/:id/pillars", cfg.TaskHandler.GetPillars)

	// Task evidences
	protected.POST("/tasks/:id/evidence", cfg.TaskHandler.AddEvidence)
	protected.GET("/tasks/:id/evidences", cfg.TaskHandler.ListEvidences)
	protected.PUT("/evidences/:id", cfg.TaskHandler.UpdateEvidence)
	protected.DELETE("/evidences/:id", cfg.TaskHandler.DeleteEvidence)

	// AI Classification
	protected.POST("/tasks/classify", cfg.AIHandler.Classify)

	// Dashboard & Reports
	protected.GET("/dashboard", cfg.InitiativeHandler.GetDashboard)
	protected.GET("/report", cfg.InitiativeHandler.GetDetailedReport)
	protected.GET("/gap-analysis", cfg.InitiativeHandler.GetGapAnalysis)
	protected.GET("/career-radar", cfg.InitiativeHandler.GetReadinessCheck)
	protected.GET("/cycle-comparison", cfg.InitiativeHandler.GetCycleComparison)
	protected.GET("/report/pdf", cfg.InitiativeHandler.DownloadReportPDF)
	protected.POST("/report/send-to-manager", cfg.InitiativeHandler.SendReportToManager)

	// Ladder (read for all users)
	protected.GET("/ladders", cfg.LadderHandler.List)
	protected.GET("/ladder/:id", cfg.LadderHandler.FindByID)

	// Level targets (user-scoped)
	protected.POST("/level-target", cfg.LevelTargetHandler.SetTarget)
	protected.GET("/level-target", cfg.LevelTargetHandler.GetTarget)
	protected.GET("/level-targets", cfg.LevelTargetHandler.ListTargets)
	protected.DELETE("/level-targets/:id", cfg.LevelTargetHandler.DeleteTarget)

	// Admin-only routes
	admin := protected.Group("/")
	admin.Use(api.AdminOnly())
	admin.DELETE("/users/:id", cfg.UserHandler.Delete)
	admin.GET("/users", cfg.UserHandler.FindAll)
	admin.POST("/users", cfg.UserHandler.Register)
	admin.PATCH("/users/:id", cfg.UserHandler.Update)

	admin.POST("/ladder", cfg.LadderHandler.Create)
	admin.PUT("/ladder/:id", cfg.LadderHandler.Update)
	admin.DELETE("/ladder/:id", cfg.LadderHandler.Delete)

	return r
}
