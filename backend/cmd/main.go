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
	_ "github.com/me/level-up-hub/backend/docs"
	"github.com/me/level-up-hub/backend/internal/account"
	"github.com/me/level-up-hub/backend/internal/database"
	"github.com/me/level-up-hub/backend/internal/email"
	"github.com/me/level-up-hub/backend/internal/initiative"
	"github.com/me/level-up-hub/backend/internal/ladder"
	"github.com/me/level-up-hub/backend/internal/logger"
	"github.com/me/level-up-hub/backend/internal/repository"
	"github.com/me/level-up-hub/backend/internal/task"
	"github.com/me/level-up-hub/backend/routes"
)

// @title           Level Up Hub API
// @version         1.0
// @description     API for career management and professional development with XP system, initiatives and reports.
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

	// Run database migrations
	if migErr := database.RunMigrations(cfg); migErr != nil {
		log.Error("failed to run migrations", slog.String("error", migErr.Error()))
		panic(migErr)
	}

	repo := repository.New(dbPool)

	// Services
	accountService := account.NewService(repo)
	ladderService := ladder.NewService(repo)
	initiativeService := initiative.NewService(repo, dbPool)
	taskService := task.NewService(repo, dbPool)
	emailService := email.NewService(cfg)

	// Handlers
	accountHandler := account.NewHandler(accountService, cfg)
	ladderHandler := ladder.NewHandler(ladderService, cfg)
	initiativeHandler := initiative.NewHandler(initiativeService, cfg, emailService)
	taskHandler := task.NewHandler(taskService, cfg)

	r := routes.NewRouter(routes.RouterConfig{
		UserHandler:       accountHandler,
		LadderHandler:     ladderHandler,
		InitiativeHandler: initiativeHandler,
		TaskHandler:       taskHandler,
	}, dbPool, cfg)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutdown signal received, initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", slog.String("error", err.Error()))
	} else {
		log.Info("server shutdown completed successfully")
	}

	log.Info("closing database connection pool...")
	dbPool.Close()
	log.Info("database connections closed")

	log.Info("application stopped gracefully")
}
