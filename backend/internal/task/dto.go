package task

import "github.com/google/uuid"

type CreateTaskDTO struct {
	InitiativeID       uuid.UUID `json:"initiative_id" binding:"required"`
	Title              string    `json:"title" binding:"required"`
	Execution          string    `json:"execution"`
	ProgressPercentage int32     `json:"progress_percentage" binding:"gte=0,lte=100"`
}

type UpdateTaskDTO struct {
	Title              string `json:"title" binding:"required"`
	Execution          string `json:"execution"`
	ProgressPercentage int32  `json:"progress_percentage" binding:"gte=0,lte=100"`
}
