package logger

import (
	"context"
	"log/slog"
	"os"
)

// Setup configura o logger estruturado baseado no ambiente
func Setup(env string) *slog.Logger {
	var handler slog.Handler

	// Em produção usa JSON, em dev usa texto legível
	if env == "prod" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

// FromContext extrai o logger do contexto, ou retorna o default
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value("logger").(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// WithRequestID adiciona request_id ao logger
func WithRequestID(logger *slog.Logger, requestID string) *slog.Logger {
	return logger.With(slog.String("request_id", requestID))
}

// WithUserID adiciona user_id ao logger
func WithUserID(logger *slog.Logger, userID string) *slog.Logger {
	return logger.With(slog.String("user_id", userID))
}
