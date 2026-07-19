package leveltarget

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/me/level-up-hub/backend/apperr"
	"github.com/me/level-up-hub/backend/internal/repository"
)

type Service struct {
	repo *repository.Queries
}

func NewService(repo *repository.Queries) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, params repository.CreateLevelTargetParams) (repository.LevelTarget, error) {
	lt, err := s.repo.CreateLevelTarget(ctx, params)
	if err != nil {
		slog.Error("failed to create level target", slog.String("error", err.Error()))
		return repository.LevelTarget{}, apperr.MessageError(fmt.Sprintf(apperr.ErrCreate, "meta de nivel"), err)
	}
	return lt, nil
}

func (s *Service) FindByID(ctx context.Context, id uuid.UUID) (repository.FindLevelTargetByIDRow, error) {
	return s.repo.FindLevelTargetByID(ctx, id)
}

func (s *Service) FindByYear(ctx context.Context, year int32) (repository.FindLevelTargetByYearRow, error) {
	return s.repo.FindLevelTargetByYear(ctx, year)
}

func (s *Service) List(ctx context.Context) ([]repository.ListLevelTargetsRow, error) {
	return s.repo.ListLevelTargets(ctx)
}

func (s *Service) Update(ctx context.Context, params repository.UpdateLevelTargetParams) (repository.LevelTarget, error) {
	lt, err := s.repo.UpdateLevelTarget(ctx, params)
	if err != nil {
		slog.Error("failed to update level target", slog.String("error", err.Error()))
		return repository.LevelTarget{}, apperr.MessageError(fmt.Sprintf(apperr.ErrUpdate, "meta de nivel"), err)
	}
	return lt, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteLevelTarget(ctx, id)
	if err != nil {
		slog.Error("failed to delete level target", slog.String("error", err.Error()))
		return apperr.MessageError(fmt.Sprintf(apperr.ErrDelete, "meta de nivel"), err)
	}
	return nil
}
