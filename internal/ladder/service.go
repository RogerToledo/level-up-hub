package ladder

import (
	"context"

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
		return  err
	}

	return nil
}
