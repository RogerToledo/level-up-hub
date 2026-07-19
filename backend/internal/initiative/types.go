package initiative

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
	Diff      int32  `json:"diff"`
}

type ComparisonReport struct {
	CurrentCycleName  string            `json:"current_cycle"`
	PreviousCycleName string            `json:"previous_cycle"`
	GrowthXP          int32             `json:"growth_xp"`
	PercentChange     float64           `json:"percent_change"`
	LevelEvolution    []LevelComparison `json:"level_evolution"`
}
