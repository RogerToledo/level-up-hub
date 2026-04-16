package activity

import "github.com/me/level-up-hub/backend/internal/repository"

type PillarStats struct {
	Achieved   int32   `json:"achieved"`
	Planned    int32   `json:"planned"`
	Percentage float64 `json:"percentage"`
}

type DashboardResponse struct {
	OfficialLevel repository.LadderLevel `json:"official_level"`
	TargetLevel   repository.LadderLevel `json:"target_level"`
	TotalAchieved int32                  `json:"total_achieved"`
	MaxPdiXp      int32                  `json:"max_pdi_xp"`
	PdiProgress   map[string]PillarStats `json:"pdi_progress"`
	Overdelivery  map[string]int32       `json:"overdelivery"`
}
