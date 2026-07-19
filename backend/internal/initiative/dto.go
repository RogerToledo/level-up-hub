package initiative

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/me/level-up-hub/backend/internal/repository"
)

type CreateInitiativeDTO struct {
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description *string   `json:"description"`
	IsPdiTarget bool      `json:"is_pdi_target"`
}

func (dto *CreateInitiativeDTO) ToRepositoryParams() repository.CreateInitiativeParams {
	params := repository.CreateInitiativeParams{
		UserID:      dto.UserID,
		Title:       dto.Title,
		IsPdiTarget: dto.IsPdiTarget,
	}

	if dto.Description != nil {
		params.Description = pgtype.Text{String: *dto.Description, Valid: true}
	} else {
		params.Description = pgtype.Text{Valid: false}
	}

	return params
}

type UpdateInitiativeDTO struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	IsPdiTarget bool    `json:"is_pdi_target"`
}

func (dto *UpdateInitiativeDTO) ToRepositoryParams(initiativeID uuid.UUID, userID uuid.UUID) repository.UpdateInitiativeParams {
	params := repository.UpdateInitiativeParams{
		ID:          initiativeID,
		UserID:      userID,
		Title:       dto.Title,
		IsPdiTarget: dto.IsPdiTarget,
	}

	if dto.Description != nil {
		params.Description = pgtype.Text{String: *dto.Description, Valid: true}
	} else {
		params.Description = pgtype.Text{Valid: false}
	}

	return params
}

type GapAnalysisResponse struct {
	Level      string `json:"level"`
	Pillar     string `json:"pillar"`
	Target     int32  `json:"target"`
	Achieved   int32  `json:"achieved"`
	Gap        int32  `json:"gap"`
	Status     string `json:"status"`
	Percentage int32  `json:"percentage"`
}
