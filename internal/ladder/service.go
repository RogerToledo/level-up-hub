package ladder

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/me/level-up-hub/apperr"
	"github.com/me/level-up-hub/internal/repository"
)

type Service struct {
	repo *repository.Queries
}

func NewService(repo *repository.Queries) *Service {
	return &Service{repo: repo}
}

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
