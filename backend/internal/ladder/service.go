package ladder

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/me/level-up-hub/backend/apperr"
	"github.com/me/level-up-hub/backend/internal/repository"
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

// ListAllLadders returns all available career ladder levels.
func (s *Service) ListAllLadders(ctx context.Context) ([]repository.CareerLadder, error) {
	ladders, err := s.repo.ListEntireLadder(ctx)
	if err != nil {
		slog.Error("failed to list ladder levels",
			slog.String("error", err.Error()),
		)
		return nil, apperr.MessageError(apperr.ErrListLadderLevels, err)
	}

	return ladders, nil
}
