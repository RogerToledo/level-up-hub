package main

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/me/level-up-hub/config"
	"github.com/me/level-up-hub/internal/account"
	"github.com/me/level-up-hub/internal/activity"
	"github.com/me/level-up-hub/internal/database"
	"github.com/me/level-up-hub/internal/ladder"
	"github.com/me/level-up-hub/internal/logger"
	"github.com/me/level-up-hub/internal/repository"
	"github.com/me/level-up-hub/routes"
)

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

	log.Info("server starting", slog.String("port", cfg.Port))

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Error("server failed to start", slog.String("error", err.Error()))
		panic(err)
	}
}
