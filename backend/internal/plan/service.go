package plan

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

func (s *Service) Create(ctx context.Context, userID uuid.UUID, input CreatePlanDTO) error {
	_, err := s.repo.FindPlanByTitle(ctx, repository.FindPlanByTitleParams{
		UserID: userID,
		Title:  input.Title,
	})
	if err == nil {
		return apperr.MessageError(fmt.Sprintf(apperr.ErrAlredyExists, "plano"), fmt.Errorf("plano com titulo '%s' ja existe", input.Title))
	}

	_, err = s.repo.CreatePlan(ctx, input.ToRepositoryParams(userID))
	if err != nil {
		slog.Error("failed to create plan", slog.String("error", err.Error()))
		return apperr.MessageError(fmt.Sprintf(apperr.ErrCreate, "plano"), err)
	}
	return nil
}

func (s *Service) Update(ctx context.Context, planID uuid.UUID, userID uuid.UUID, dto UpdatePlanDTO) error {
	params := dto.ToRepositoryParams(planID, userID)
	_, err := s.repo.UpdatePlan(ctx, params)
	if err != nil {
		slog.Error("failed to update plan", slog.String("error", err.Error()))
		return apperr.MessageError(fmt.Sprintf(apperr.ErrUpdate, "plano"), err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, planID uuid.UUID, userID uuid.UUID) error {
	err := s.repo.DeletePlan(ctx, repository.DeletePlanParams{
		ID:     planID,
		UserID: userID,
	})
	if err != nil {
		return apperr.MessageError(fmt.Sprintf(apperr.ErrDelete, "plano"), err)
	}
	return nil
}

func (s *Service) FindByID(ctx context.Context, planID uuid.UUID, userID uuid.UUID) (repository.Plan, error) {
	plan, err := s.repo.FindPlanByID(ctx, repository.FindPlanByIDParams{
		ID:     planID,
		UserID: userID,
	})
	if err != nil {
		return repository.Plan{}, apperr.MessageError(fmt.Sprintf(apperr.ErrIsNotFound, "plano"), err)
	}
	return plan, nil
}

func (s *Service) ListPlans(ctx context.Context, userID uuid.UUID) ([]repository.ListUserPlansRow, error) {
	return s.repo.ListUserPlans(ctx, userID)
}

func (s *Service) Reorder(ctx context.Context, planID uuid.UUID, userID uuid.UUID, newPosition int32) error {
	// Get current plan
	plan, err := s.repo.FindPlanByID(ctx, repository.FindPlanByIDParams{
		ID:     planID,
		UserID: userID,
	})
	if err != nil {
		return apperr.MessageError(fmt.Sprintf(apperr.ErrIsNotFound, "plano"), err)
	}

	currentPosition := plan.Position

	// Get total plans for this user
	count, err := s.repo.CountUserPlans(ctx, userID)
	if err != nil {
		return apperr.MessageError(fmt.Sprintf(apperr.ErrFindAll, "planos"), err)
	}

	// Validate new position
	if newPosition < 0 || newPosition >= int32(count) {
		return apperr.MessageError(apperr.ErrBadRequest, fmt.Errorf("posicao invalida"))
	}

	// If same position, do nothing
	if currentPosition == newPosition {
		return nil
	}

	// Update positions of affected plans
	if currentPosition < newPosition {
		// Moving down: shift plans up
		for pos := currentPosition + 1; pos <= newPosition; pos++ {
			otherPlan, err := s.repo.FindPlanByPosition(ctx, repository.FindPlanByPositionParams{
				UserID:   userID,
				Position: pos,
			})
			if err != nil {
				continue
			}
			_ = s.repo.ReorderPlan(ctx, repository.ReorderPlanParams{
				ID:       otherPlan.ID,
				UserID:   userID,
				Position: pos - 1,
			})
		}
	} else {
		// Moving up: shift plans down
		for pos := currentPosition - 1; pos >= newPosition; pos-- {
			otherPlan, err := s.repo.FindPlanByPosition(ctx, repository.FindPlanByPositionParams{
				UserID:   userID,
				Position: pos,
			})
			if err != nil {
				continue
			}
			_ = s.repo.ReorderPlan(ctx, repository.ReorderPlanParams{
				ID:       otherPlan.ID,
				UserID:   userID,
				Position: pos + 1,
			})
		}
	}

	// Set the plan to its new position
	return s.repo.ReorderPlan(ctx, repository.ReorderPlanParams{
		ID:       planID,
		UserID:   userID,
		Position: newPosition,
	})
}