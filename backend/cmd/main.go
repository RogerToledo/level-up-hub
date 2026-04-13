package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/me/level-up-hub/backend/config"
	_ "github.com/me/level-up-hub/backend/docs" // Swagger docs
	"github.com/me/level-up-hub/backend/internal/account"
	"github.com/me/level-up-hub/backend/internal/activity"
	"github.com/me/level-up-hub/backend/internal/database"
	"github.com/me/level-up-hub/backend/internal/ladder"
	"github.com/me/level-up-hub/backend/internal/logger"
	"github.com/me/level-up-hub/backend/internal/repository"
	"github.com/me/level-up-hub/backend/routes"
)

// @title           Level Up Hub API
// @version         1.0
// @description     API para gerenciamento de carreira e desenvolvimento profissional com sistema de XP, atividades e relatórios.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@leveluphub.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8081
// @BasePath  /v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Digite 'Bearer' seguido do token JWT

func main() {
	cfg := config.LoadConfig()

	// Configura logger estruturado
	log := logger.Setup(cfg.Env)

	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Info("application starting",
		slog.String("env", cfg.Env),
		slog.String("port", cfg.Port),
	)

	dbPool, err := database.NewPostgresPool(context.Background(), cfg)
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		panic(err)
	}
	defer dbPool.Close()

	log.Info("database connected",
		slog.String("env", cfg.Env),
		slog.Int("max_conns", cfg.MaxConns),
		slog.Int("min_conns", cfg.MinConns),
	)

	repo := repository.New(dbPool)
	service := account.NewService(repo)
	handler := account.NewHandler(service, cfg)
	ladderService := ladder.NewService(repo)
	ladderHandler := ladder.NewHandler(ladderService, cfg)
	activityService := activity.NewService(repo, dbPool)
	activityHandler := activity.NewHandler(activityService, cfg)

	r := routes.NewRouter(routes.RouterConfig{
		UserHandler:     handler,
		LadderHandler:   ladderHandler,
		ActivityHandler: activityHandler,
	}, dbPool, cfg)

	// Configure HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
		// Timeout configurations for security
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info("server starting",
			slog.String("port", cfg.Port),
			slog.String("env", cfg.Env),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed to start", slog.String("error", err.Error()))
			panic(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// SIGINT (Ctrl+C) and SIGTERM (kill) are captured
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutdown signal received, initiating graceful shutdown...")

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", slog.String("error", err.Error()))
	} else {
		log.Info("server shutdown completed successfully")
	}

	// Close database pool
	log.Info("closing database connection pool...")
	dbPool.Close()
	log.Info("database connections closed")

	log.Info("application stopped gracefully")
}
