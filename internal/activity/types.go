package activity

type PillarStats struct {
	Achieved   int32   `json:"achieved"`
	Planned    int32   `json:"planned"`
	Percentage float64 `json:"percentage"`
}

type DashboardResponse struct {
    CurrentLevel   string                 `json:"current_level"`
    PdiProgress    map[string]PillarStats `json:"pdi_progress"`
    MaxPdiXp       int32                  `json:"max_pdi_xp"`
    TotalAchieved  int32                  `json:"total_achieved"`
    Overdelivery   map[string]int32       `json:"overdelivery"`
}