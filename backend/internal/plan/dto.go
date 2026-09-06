package plan

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/me/level-up-hub/backend/internal/repository"
)

type CreatePlanDTO struct {
	Title        string  `json:"title" binding:"required"`
	Description  *string `json:"description"`
	InitiativeID *string `json:"initiative_id"`
	LevelTarget  *string `json:"level_target"`
}

func (dto *CreatePlanDTO) ToRepositoryParams(userID uuid.UUID) repository.CreatePlanParams {
	params := repository.CreatePlanParams{
		UserID:  userID,
		Title:   dto.Title,
		Status:  "active",
	}

	if dto.Description != nil {
		params.Description = pgtype.Text{String: *dto.Description, Valid: true}
	} else {
		params.Description = pgtype.Text{Valid: false}
	}

	if dto.InitiativeID != nil {
		id, err := uuid.Parse(*dto.InitiativeID)
		if err == nil {
			params.InitiativeID = pgtype.UUID{Bytes: id, Valid: true}
		}
	}

	if dto.LevelTarget != nil {
		params.LevelTarget = pgtype.Text{String: *dto.LevelTarget, Valid: true}
	} else {
		params.LevelTarget = pgtype.Text{Valid: false}
	}

	return params
}

type UpdatePlanDTO struct {
	Title        string  `json:"title" binding:"required"`
	Description  *string `json:"description"`
	InitiativeID *string `json:"initiative_id"`
	LevelTarget  *string `json:"level_target"`
	Status       *string `json:"status"`
}

func (dto *UpdatePlanDTO) ToRepositoryParams(planID uuid.UUID, userID uuid.UUID) repository.UpdatePlanParams {
	params := repository.UpdatePlanParams{
		ID:     planID,
		UserID: userID,
		Title:  dto.Title,
		Status: "active",
	}

	if dto.Description != nil {
		params.Description = pgtype.Text{String: *dto.Description, Valid: true}
	} else {
		params.Description = pgtype.Text{Valid: false}
	}

	if dto.InitiativeID != nil {
		id, err := uuid.Parse(*dto.InitiativeID)
		if err == nil {
			params.InitiativeID = pgtype.UUID{Bytes: id, Valid: true}
		}
	}

	if dto.LevelTarget != nil {
		params.LevelTarget = pgtype.Text{String: *dto.LevelTarget, Valid: true}
	} else {
		params.LevelTarget = pgtype.Text{Valid: false}
	}

	if dto.Status != nil {
		params.Status = *dto.Status
	}

	return params
}

type PlanResponse struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	Description     *string `json:"description"`
	InitiativeID    *string `json:"initiative_id"`
	InitiativeTitle *string `json:"initiative_title"`
	LevelTarget     *string `json:"level_target"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}