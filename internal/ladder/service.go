package ladder

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/me/level-up-hub/apperr"
	"github.com/me/level-up-hub/internal/repository"
)

// Service provides business logic for career ladder operations.
type Service struct {
	repo *repository.Queries
}

// NewService creates a new ladder service.
func NewService(repo *repository.Queries) *Service {
	return &Service{repo: repo}
}

// CreateLadderLevel creates a new career ladder level.
func (s *Service) CreateLadderLevel(ctx context.Context, params repository.CreateLadderLevelParams) error {
	_, err := s.repo.CreateLadderLevel(ctx, params)
	if err != nil {
		slog.Error("failed to create ladder level",
			slog.String("error", err.Error()),
		)
		return apperr.MessageError(fmt.Sprintf(apperr.ErrCreate, apperr.LadderLevelPT), err)
	}

	return nil
}
