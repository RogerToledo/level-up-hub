package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/me/level-up-hub/internal/config"
	"github.com/me/level-up-hub/internal/database"
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

	r := gin.Default()
	r.Run(":" + cfg.Port)
}
