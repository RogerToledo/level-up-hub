package config

import (
	"log/slog"
	"sync"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	DbUrlDev  string `env:"DB_URL_DEV" envDefault:"postgres://localhost:5432/dev_db"`
	DbUrlProd string `env:"DB_URL_PROD" envDefault:"postgres://localhost:5432/prod_db"`
	Port      string `env:"PORT" envDefault:"8081"`
	Env       string `env:"ENV" envDefault:"dev"`
	MaxConns  int    `env:"MAX_CONNS" envDefault:"10"`
	MinConns  int    `env:"MIN_CONNS" envDefault:"1"`
	JWTSecret string `env:"JWT_SECRET" envDefault:"supersecretkey"`
}

var (
	cfg  *Config
	onde sync.Once
)

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
