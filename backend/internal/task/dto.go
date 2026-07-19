package task

import (
	"github.com/google/uuid"
	"github.com/me/level-up-hub/backend/internal/repository"
)

type CreateTaskDTO struct {
	InitiativeID       uuid.UUID           `json:"initiative_id" binding:"required"`
	LadderID           uuid.UUID           `json:"ladder_id" binding:"required"`
	Pillars            []repository.Pillar `json:"pillars" binding:"required"`
	Title              string              `json:"title" binding:"required"`
	Execution          string              `json:"execution"`
	ImpactSummary      string              `json:"impact_summary"`
	ProgressPercentage int32               `json:"progress_percentage" binding:"gte=0,lte=100"`
}

type UpdateTaskDTO struct {
	LadderID           uuid.UUID           `json:"ladder_id" binding:"required"`
	Pillars            []repository.Pillar `json:"pillars"`
	Title              string              `json:"title" binding:"required"`
	Execution          string              `json:"execution"`
	ImpactSummary      string              `json:"impact_summary"`
	ProgressPercentage int32               `json:"progress_percentage" binding:"gte=0,lte=100"`
}
