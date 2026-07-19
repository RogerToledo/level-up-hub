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
	params := repository.CreateTaskParams{
		InitiativeID:       input.InitiativeID,
		Title:              input.Title,
		Execution:          pgtype.Text{String: input.Execution, Valid: input.Execution != ""},
		ProgressPercentage: input.ProgressPercentage,
	}

	task, err := s.repo.CreateTask(ctx, params)
	if err != nil {
		slog.Error("failed to create task",
			slog.String("error", err.Error()),
			slog.String("initiative_id", input.InitiativeID.String()),
		)
		return repository.Task{}, apperr.MessageError(fmt.Sprintf(apperr.ErrCreate, "task"), err)
	}

	// Recalculate initiative progress
	if err := s.recalculateInitiativeProgress(ctx, input.InitiativeID); err != nil {
		slog.Error("failed to recalculate initiative progress", slog.String("error", err.Error()))
	}

	return task, nil
}

func (s *Service) Update(ctx context.Context, taskID uuid.UUID, input UpdateTaskDTO) (repository.Task, error) {
	params := repository.UpdateTaskParams{
		ID:                 taskID,
		Title:              input.Title,
		Execution:          pgtype.Text{String: input.Execution, Valid: input.Execution != ""},
		ProgressPercentage: input.ProgressPercentage,
	}

	task, err := s.repo.UpdateTask(ctx, params)
	if err != nil {
		slog.Error("failed to update task",
			slog.String("error", err.Error()),
			slog.String("task_id", taskID.String()),
		)
		return repository.Task{}, apperr.MessageError(fmt.Sprintf(apperr.ErrUpdate, "task"), err)
	}

	// Recalculate initiative progress
	if err := s.recalculateInitiativeProgress(ctx, task.InitiativeID); err != nil {
		slog.Error("failed to recalculate initiative progress", slog.String("error", err.Error()))
	}

	return task, nil
}

func (s *Service) Delete(ctx context.Context, taskID uuid.UUID) error {
	// Get the task first to know the initiative_id
	task, err := s.repo.FindTaskByID(ctx, taskID)
	if err != nil {
		return apperr.MessageError(fmt.Sprintf(apperr.ErrDelete, "task"), err)
	}

	err = s.repo.DeleteTask(ctx, taskID)
	if err != nil {
		return apperr.MessageError(fmt.Sprintf(apperr.ErrDelete, "task"), err)
	}

	// Recalculate initiative progress
	if err := s.recalculateInitiativeProgress(ctx, task.InitiativeID); err != nil {
		slog.Error("failed to recalculate initiative progress", slog.String("error", err.Error()))
	}

	return nil
}

func (s *Service) ListByInitiative(ctx context.Context, initiativeID uuid.UUID) ([]repository.ListTasksByInitiativeRow, error) {
	return s.repo.ListTasksByInitiative(ctx, initiativeID)
}

func (s *Service) AddEvidence(ctx context.Context, taskID uuid.UUID, url string, description string) (repository.TaskEvidence, error) {
	params := repository.CreateTaskEvidenceParams{
		TaskID:      taskID,
		EvidenceUrl: url,
		Description: pgtype.Text{String: description, Valid: description != ""},
	}

	evidence, err := s.repo.CreateTaskEvidence(ctx, params)
	if err != nil {
		slog.Error("failed to add task evidence",
			slog.String("error", err.Error()),
			slog.String("task_id", taskID.String()),
		)
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

	// Find the initiative owner to pass user_id
	_, err = s.pool.Exec(ctx,
		"UPDATE initiatives SET progress_percentage = $1, updated_at = CURRENT_DATE WHERE id = $2",
		avg, initiativeID,
	)
	return err
}
