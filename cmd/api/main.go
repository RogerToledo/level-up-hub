package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/me/level-up-hub/internal/account"
	"github.com/me/level-up-hub/config"
	"github.com/me/level-up-hub/internal/database"
	"github.com/me/level-up-hub/internal/repository"
	"github.com/me/level-up-hub/routes"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	fmt.Printf("Connected to %s mode\n", cfg.Env)

	dbPool, err := database.NewPostgresPool(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	fmt.Printf("Connected to %s pool\n", cfg.Env)

	repo := repository.New(dbPool)
	service := account.NewService(repo)
	handler := account.NewHandler(service, cfg)

	r := routes.NewRouter(routes.RouterConfig{
		UserHandler: handler,
	}, dbPool, cfg)

	fmt.Printf("🚀 Server starting on port %s\n", cfg.Port)
	
	r.Run(":" + cfg.Port)
}
