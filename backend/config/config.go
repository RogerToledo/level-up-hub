// Package config provides application configuration loading and management.
package config

import (
	"log/slog"
	"sync"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// Config holds all application configuration settings loaded from environment variables.
type Config struct {
	DbURLDev  string `env:"DB_URL_DEV" envDefault:"postgres://localhost:5432/dev_db"`
	DbURLProd string `env:"DB_URL_PROD" envDefault:"postgres://localhost:5432/prod_db"`
	Port      string `env:"PORT" envDefault:"8081"`
	Env       string `env:"ENV" envDefault:"dev"`

	// Connection Pool Settings
	MaxConns          int `env:"MAX_CONNS" envDefault:"25"`            // Maximum connections in pool
	MinConns          int `env:"MIN_CONNS" envDefault:"5"`             // Minimum connections maintained
	MaxConnLifetime   int `env:"MAX_CONN_LIFETIME" envDefault:"3600"`  // Maximum connection lifetime (seconds)
	MaxConnIdleTime   int `env:"MAX_CONN_IDLE_TIME" envDefault:"1800"` // Maximum idle time before closing (seconds)
	HealthCheckPeriod int `env:"HEALTH_CHECK_PERIOD" envDefault:"60"`  // Period between health checks (seconds)
	ConnectTimeout    int `env:"CONNECT_TIMEOUT" envDefault:"5"`       // Timeout to connect (seconds)

	JWTSecret string `env:"JWT_SECRET" envDefault:"supersecretkey"`
}

var (
	cfg  *Config
	onde sync.Once
)

// LoadConfig loads and returns the application configuration.
// It uses sync.Once to ensure configuration is loaded only once.
func LoadConfig() *Config {
	onde.Do(func() {
		err := godotenv.Load()
		if err != nil {
			slog.Warn("no .env file found, using environment variables")
		}

		cfg = &Config{}
		if err := env.Parse(cfg); err != nil {
			slog.Error("failed to parse config", slog.String("error", err.Error()))
			panic(err)
		}

	})
	return cfg
}
