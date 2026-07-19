package ladder

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
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

// FindByID returns a specific career ladder level by ID.
func (s *Service) FindByID(ctx context.Context, id uuid.UUID) (repository.CareerLadder, error) {
	ladder, err := s.repo.FindLadderLevel(ctx, id)
	if err != nil {
		slog.Error("failed to find ladder level",
			slog.String("error", err.Error()),
			slog.String("id", id.String()),
		)
		return repository.CareerLadder{}, err
	}
	return ladder, nil
}

// UpdateLadderLevel updates an existing career ladder level.
func (s *Service) UpdateLadderLevel(ctx context.Context, params repository.UpdateLadderLevelParams) error {
	err := s.repo.UpdateLadderLevel(ctx, params)
	if err != nil {
		slog.Error("failed to update ladder level",
			slog.String("error", err.Error()),
			slog.String("id", params.ID.String()),
		)
		return apperr.MessageError(fmt.Sprintf(apperr.ErrUpdate, apperr.LadderLevelPT), err)
	}
	return nil
}

// DeleteLadderLevel deletes a career ladder level by ID.
func (s *Service) DeleteLadderLevel(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteLadderLevel(ctx, id)
	if err != nil {
		slog.Error("failed to delete ladder level",
			slog.String("error", err.Error()),
			slog.String("id", id.String()),
		)
		return apperr.MessageError(fmt.Sprintf(apperr.ErrDelete, apperr.LadderLevelPT), err)
	}
	return nil
}
