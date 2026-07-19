package leveltarget

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/me/level-up-hub/backend/apperr"
	"github.com/me/level-up-hub/backend/internal/repository"
)

type Service struct {
	repo *repository.Queries
}

func NewService(repo *repository.Queries) *Service {
	return &Service{repo: repo}
}

func toPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func (s *Service) Upsert(ctx context.Context, userID uuid.UUID, target repository.LadderLevel, year int32) (repository.LevelTarget, error) {
	lt, err := s.repo.UpsertLevelTarget(ctx, repository.UpsertLevelTargetParams{
		UserID: toPgUUID(userID),
		Target: target,
		Year:   year,
	})
	if err != nil {
		slog.Error("failed to upsert level target", slog.String("error", err.Error()))
		return repository.LevelTarget{}, apperr.MessageError(fmt.Sprintf(apperr.ErrCreate, "meta de nivel"), err)
	}
	return lt, nil
}

func (s *Service) FindByUserAndYear(ctx context.Context, userID uuid.UUID, year int32) (repository.FindLevelTargetByUserAndYearRow, error) {
	return s.repo.FindLevelTargetByUserAndYear(ctx, repository.FindLevelTargetByUserAndYearParams{
		UserID: toPgUUID(userID),
		Year:   year,
	})
}

func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID) ([]repository.ListLevelTargetsByUserRow, error) {
	return s.repo.ListLevelTargetsByUser(ctx, toPgUUID(userID))
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	err := s.repo.DeleteLevelTarget(ctx, repository.DeleteLevelTargetParams{
		ID:     id,
		UserID: toPgUUID(userID),
	})
	if err != nil {
		slog.Error("failed to delete level target", slog.String("error", err.Error()))
		return apperr.MessageError(fmt.Sprintf(apperr.ErrDelete, "meta de nivel"), err)
	}
	return nil
}
