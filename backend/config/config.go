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
	DbURLDev    string `env:"DB_URL_DEV" envDefault:"postgres://localhost:5432/dev_db"`
	DbURLProd   string `env:"DB_URL_PROD" envDefault:"postgres://localhost:5432/prod_db"`
	DatabaseURL string `env:"DATABASE_URL"` // Railway injects this automatically
	Port        string `env:"PORT" envDefault:"8081"`
	Env         string `env:"ENV" envDefault:"dev"`

	// Connection Pool Settings
	MaxConns          int `env:"MAX_CONNS" envDefault:"25"`            // Maximum connections in pool
	MinConns          int `env:"MIN_CONNS" envDefault:"5"`             // Minimum connections maintained
	MaxConnLifetime   int `env:"MAX_CONN_LIFETIME" envDefault:"3600"`  // Maximum connection lifetime (seconds)
	MaxConnIdleTime   int `env:"MAX_CONN_IDLE_TIME" envDefault:"1800"` // Maximum idle time before closing (seconds)
	HealthCheckPeriod int `env:"HEALTH_CHECK_PERIOD" envDefault:"60"`  // Period between health checks (seconds)
	ConnectTimeout    int `env:"CONNECT_TIMEOUT" envDefault:"5"`       // Timeout to connect (seconds)

	JWTSecret      string   `env:"JWT_SECRET" envDefault:"supersecretkey"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS" envSeparator:"," envDefault:"http://localhost:3000"`

	// Migration Settings
	MigrationsPath        string `env:"MIGRATIONS_PATH" envDefault:"backend/db/migrations"`
	ForceMigrationVersion int    `env:"FORCE_MIGRATION_VERSION" envDefault:"0"`

	// Email/SMTP Settings
	SMTPHost     string `env:"SMTP_HOST" envDefault:"smtp.gmail.com"`
	SMTPPort     int    `env:"SMTP_PORT" envDefault:"587"`
	SMTPUser     string `env:"SMTP_USER" envDefault:""`
	SMTPPassword string `env:"SMTP_PASSWORD" envDefault:""`
	SMTPFrom     string `env:"SMTP_FROM" envDefault:"noreply@leveluphub.com"`
	SMTPFromName string `env:"SMTP_FROM_NAME" envDefault:"Level Up Hub"`

	// Resend API (used instead of SMTP in production)
	ResendAPIKey string `env:"RESEND_API_KEY" envDefault:""`

	// Gemini AI
	GeminiAPIKey string `env:"GEMINI_API_KEY" envDefault:""`
}

var (
	cfg  *Config
	onde sync.Once
)

// LoadConfig loads and returns the application configuration.
// It uses sync.Once to ensure configuration is loaded only once.
func LoadConfig() *Config {
	onde.Do(func() {
		err := godotenv.Load("backend/.env")
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

// GetDatabaseURL returns the appropriate database URL.
// Priority: DATABASE_URL (Railway) > DB_URL_PROD (prod) > DB_URL_DEV (dev)
func (c *Config) GetDatabaseURL() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	if c.Env == "prod" {
		return c.DbURLProd
	}
	return c.DbURLDev
}
