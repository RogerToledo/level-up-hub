package task

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/me/level-up-hub/backend/apperr"
	"github.com/me/level-up-hub/backend/internal/repository"
)

type Service struct {
	repo *repository.Queries
	pool *pgxpool.Pool
}

func NewService(repo *repository.Queries, pool *pgxpool.Pool) *Service {
	return &Service{repo: repo, pool: pool}
}

func (s *Service) Create(ctx context.Context, input CreateTaskDTO) (repository.Task, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return repository.Task{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	repoWithTx := s.repo.WithTx(tx)

	params := repository.CreateTaskParams{
		InitiativeID:       input.InitiativeID,
		LadderID:           input.LadderID,
		Title:              input.Title,
		Execution:          pgtype.Text{String: input.Execution, Valid: input.Execution != ""},
		ImpactSummary:      pgtype.Text{String: input.ImpactSummary, Valid: input.ImpactSummary != ""},
		ProgressPercentage: input.ProgressPercentage,
	}

	task, err := repoWithTx.CreateTask(ctx, params)
	if err != nil {
		slog.Error("failed to create task", slog.String("error", err.Error()))
		return repository.Task{}, apperr.MessageError(fmt.Sprintf(apperr.ErrCreate, "task"), err)
	}

	for _, p := range input.Pillars {
		_, err = repoWithTx.CreateTaskPillar(ctx, repository.CreateTaskPillarParams{
			TaskID: task.ID,
			Pillar: p,
		})
		if err != nil {
			slog.Error("failed to create task pillar", slog.String("error", err.Error()))
			return repository.Task{}, apperr.MessageError(fmt.Sprintf(apperr.ErrCreate, "pilar"), err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return repository.Task{}, err
	}

	// Recalculate initiative progress
	if err := s.recalculateInitiativeProgress(ctx, input.InitiativeID); err != nil {
		slog.Error("failed to recalculate progress", slog.String("error", err.Error()))
	}

	return task, nil
}

func (s *Service) Update(ctx context.Context, taskID uuid.UUID, input UpdateTaskDTO) (repository.Task, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return repository.Task{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	repoWithTx := s.repo.WithTx(tx)

	params := repository.UpdateTaskParams{
		ID:                 taskID,
		Title:              input.Title,
		LadderID:           input.LadderID,
		Execution:          pgtype.Text{String: input.Execution, Valid: input.Execution != ""},
		ImpactSummary:      pgtype.Text{String: input.ImpactSummary, Valid: input.ImpactSummary != ""},
		ProgressPercentage: input.ProgressPercentage,
	}

	task, err := repoWithTx.UpdateTask(ctx, params)
	if err != nil {
		slog.Error("failed to update task", slog.String("error", err.Error()))
		return repository.Task{}, apperr.MessageError(fmt.Sprintf(apperr.ErrUpdate, "task"), err)
	}

	// Sync pillars if provided
	if len(input.Pillars) > 0 {
		if err := repoWithTx.DeleteTaskPillars(ctx, taskID); err != nil {
			return repository.Task{}, apperr.MessageError(fmt.Sprintf(apperr.ErrUpdate, "pilar"), err)
		}
		for _, p := range input.Pillars {
			_, err = repoWithTx.CreateTaskPillar(ctx, repository.CreateTaskPillarParams{
				TaskID: taskID,
				Pillar: p,
			})
			if err != nil {
				return repository.Task{}, apperr.MessageError(fmt.Sprintf(apperr.ErrUpdate, "pilar"), err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return repository.Task{}, err
	}

	// Recalculate initiative progress
	if err := s.recalculateInitiativeProgress(ctx, task.InitiativeID); err != nil {
		slog.Error("failed to recalculate progress", slog.String("error", err.Error()))
	}

	return task, nil
}

func (s *Service) Delete(ctx context.Context, taskID uuid.UUID) error {
	task, err := s.repo.FindTaskByID(ctx, taskID)
	if err != nil {
		return apperr.MessageError(fmt.Sprintf(apperr.ErrDelete, "task"), err)
	}

	err = s.repo.DeleteTask(ctx, taskID)
	if err != nil {
		return apperr.MessageError(fmt.Sprintf(apperr.ErrDelete, "task"), err)
	}

	if err := s.recalculateInitiativeProgress(ctx, task.InitiativeID); err != nil {
		slog.Error("failed to recalculate progress", slog.String("error", err.Error()))
	}

	return nil
}

func (s *Service) ListByInitiative(ctx context.Context, initiativeID uuid.UUID) ([]repository.ListTasksByInitiativeRow, error) {
	return s.repo.ListTasksByInitiative(ctx, initiativeID)
}

func (s *Service) GetTaskPillars(ctx context.Context, taskID uuid.UUID) ([]repository.Pillar, error) {
	return s.repo.GetTaskPillars(ctx, taskID)
}

func (s *Service) AddEvidence(ctx context.Context, taskID uuid.UUID, url string, description string) (repository.TaskEvidence, error) {
	params := repository.CreateTaskEvidenceParams{
		TaskID:      taskID,
		EvidenceUrl: url,
		Description: pgtype.Text{String: description, Valid: description != ""},
	}

	evidence, err := s.repo.CreateTaskEvidence(ctx, params)
	if err != nil {
		slog.Error("failed to add task evidence", slog.String("error", err.Error()))
		return repository.TaskEvidence{}, apperr.MessageError(fmt.Sprintf(apperr.ErrCreate, apperr.EvidencePT), err)
	}

	return evidence, nil
}

func (s *Service) ListEvidences(ctx context.Context, taskID uuid.UUID) ([]repository.TaskEvidence, error) {
	return s.repo.ListEvidencesByTask(ctx, taskID)
}

func (s *Service) DeleteEvidence(ctx context.Context, evidenceID uuid.UUID) error {
	return s.repo.DeleteTaskEvidence(ctx, evidenceID)
}

func (s *Service) recalculateInitiativeProgress(ctx context.Context, initiativeID uuid.UUID) error {
	avg, err := s.repo.CalculateInitiativeProgress(ctx, initiativeID)
	if err != nil {
		return err
	}

	_, err = s.pool.Exec(ctx,
		"UPDATE initiatives SET progress_percentage = $1, updated_at = CURRENT_DATE WHERE id = $2",
		avg, initiativeID,
	)
	return err
}
