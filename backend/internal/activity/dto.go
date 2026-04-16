package activity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/me/level-up-hub/backend/internal/repository"
)

type CreateActivityDTO struct {
	UserID             uuid.UUID           `json:"user_id" binding:"required"`
	LadderID           uuid.UUID           `json:"ladder_id" binding:"required"`
	Pillars            []repository.Pillar `json:"pillars" binding:"required"`
	Title              string              `json:"title" binding:"required"`
	Description        *string             `json:"description"`
	ProgressPercentage int32               `json:"progress_percentage" binding:"required,min=0,max=100"`
	ImpactSummary      *string             `json:"impact_summary"`
	IsPdiTarget        bool                `json:"is_pdi_target"`
}

func (dto *CreateActivityDTO) ToRepositoryParams() repository.CreateActivityParams {
	params := repository.CreateActivityParams{
		UserID:             dto.UserID,
		LadderID:           dto.LadderID,
		Title:              dto.Title,
		ProgressPercentage: dto.ProgressPercentage,
		IsPdiTarget:        dto.IsPdiTarget,
	}

	// Converte ponteiros de string para pgtype.Text
	if dto.Description != nil {
		params.Description = pgtype.Text{String: *dto.Description, Valid: true}
	} else {
		params.Description = pgtype.Text{Valid: false}
	}

	if dto.ImpactSummary != nil {
		params.ImpactSummary = pgtype.Text{String: *dto.ImpactSummary, Valid: true}
	} else {
		params.ImpactSummary = pgtype.Text{Valid: false}
	}

	return params
}

type UpdateActivityDTO struct {
	Title              string  `json:"title" binding:"required"`
	Description        *string `json:"description"`
	ProgressPercentage int32   `json:"progress_percentage" binding:"required,min=0,max=100"`
	ImpactSummary      *string `json:"impact_summary"`
	IsPdiTarget        bool    `json:"is_pdi_target"`
}

func (dto *UpdateActivityDTO) ToRepositoryParams(activityID uuid.UUID, userID uuid.UUID) repository.UpdateActivityParams {
	params := repository.UpdateActivityParams{
		ID:                 activityID,
		UserID:             userID,
		Title:              dto.Title,
		ProgressPercentage: dto.ProgressPercentage,
		IsPdiTarget:        dto.IsPdiTarget,
	}

	// Converte ponteiros de string para pgtype.Text
	if dto.Description != nil {
		params.Description = pgtype.Text{String: *dto.Description, Valid: true}
	} else {
		params.Description = pgtype.Text{Valid: false}
	}

	if dto.ImpactSummary != nil {
		params.ImpactSummary = pgtype.Text{String: *dto.ImpactSummary, Valid: true}
	} else {
		params.ImpactSummary = pgtype.Text{Valid: false}
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

type ReadinessCheck struct {
	IsConsistent bool   `json:"is_consistent"`
	Message      string `json:"message"`
	TargetLevel  string `json:"target_level"`
	TargetCount  int32  `json:"target_count"`
	OthersCount  int32  `json:"others_count"`
}

type LevelComposition struct {
	LevelName     string  `json:"level_name"`
	ActivityCount int32   `json:"activity_count"`
	TotalXP       int32   `json:"total_xp"`
	VolumePercent float64 `json:"volume_percent"`
	XpPercent     float64 `json:"xp_percent"`
}

type CareerRadar struct {
	TotalActivities int32              `json:"total_activities"`
	TotalXP         int32              `json:"total_xp"`
	Breakdown       []LevelComposition `json:"breakdown"`
}

type LevelComparison struct {
	LevelName string `json:"level_name"`
	CurrentXP int32  `json:"current_xp"`
	PrevXP    int32  `json:"prev_xp"`
	Diff      int32  `json:"diff"` // If positive, you accelerated. If negative, you decelerated.
}

type ComparisonReport struct {
	CurrentCycleName  string            `json:"current_cycle"`
	PreviousCycleName string            `json:"previous_cycle"`
	GrowthXP          int32             `json:"growth_xp"`      // Total XP difference
	PercentChange     float64           `json:"percent_change"` // % change in XP
	LevelEvolution    []LevelComparison `json:"level_evolution"`
}
