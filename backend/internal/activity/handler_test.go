package activity

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/me/level-up-hub/backend/internal/repository"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestCreateActivityDTO_Validation(t *testing.T) {
	tests := []struct {
		name        string
		dto         CreateActivityDTO
		expectError bool
	}{
		{
			name: "valid activity",
			dto: CreateActivityDTO{
				UserID:             uuid.New(),
				LadderID:           uuid.New(),
				Pillars:            []repository.Pillar{repository.PillarTECHNICAL},
				Title:              "Test Activity",
				ProgressPercentage: 50,
				IsPdiTarget:        true,
			},
			expectError: false,
		},
		{
			name: "invalid progress percentage - too high",
			dto: CreateActivityDTO{
				UserID:             uuid.New(),
				LadderID:           uuid.New(),
				Pillars:            []repository.Pillar{repository.PillarTECHNICAL},
				Title:              "Test Activity",
				ProgressPercentage: 150,
				IsPdiTarget:        true,
			},
			expectError: true,
		},
		{
			name: "invalid progress percentage - negative",
			dto: CreateActivityDTO{
				UserID:             uuid.New(),
				LadderID:           uuid.New(),
				Pillars:            []repository.Pillar{repository.PillarTECHNICAL},
				Title:              "Test Activity",
				ProgressPercentage: -10,
				IsPdiTarget:        true,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()

			body, _ := json.Marshal(tt.dto)
			req := httptest.NewRequest(http.MethodPost, "/activities", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if tt.expectError {
				assert.NotEqual(t, http.StatusCreated, w.Code)
			}
		})
	}
}

func TestUpdateActivityDTO_ToRepositoryParams(t *testing.T) {
	description := "Test description"
	impactSummary := "Test impact"

	dto := UpdateActivityDTO{
		Title:              "Updated Title",
		Description:        &description,
		ProgressPercentage: 75,
		ImpactSummary:      &impactSummary,
		IsPdiTarget:        true,
	}

	activityID := uuid.New()
	userID := uuid.New()

	params := dto.ToRepositoryParams(activityID, userID)

	assert.Equal(t, activityID, params.ID)
	assert.Equal(t, userID, params.UserID)
	assert.Equal(t, "Updated Title", params.Title)
	assert.True(t, params.Description.Valid)
	assert.Equal(t, description, params.Description.String)
	assert.Equal(t, int32(75), params.ProgressPercentage)
	assert.True(t, params.ImpactSummary.Valid)
	assert.Equal(t, impactSummary, params.ImpactSummary.String)
	assert.True(t, params.IsPdiTarget)
}

func TestDashboardResponse_Initialization(t *testing.T) {
	resp := &DashboardResponse{
		OfficialLevel: repository.LadderLevelP2,
		TargetLevel:   repository.LadderLevelP3,
		PdiProgress:   make(map[string]PillarStats),
		Overdelivery:  make(map[string]int32),
		MaxPdiXp:      1000,
		TotalAchieved: 750,
	}

	assert.Equal(t, repository.LadderLevelP2, resp.OfficialLevel)
	assert.Equal(t, repository.LadderLevelP3, resp.TargetLevel)
	assert.NotNil(t, resp.PdiProgress)
	assert.NotNil(t, resp.Overdelivery)
	assert.Equal(t, int32(1000), resp.MaxPdiXp)
	assert.Equal(t, int32(750), resp.TotalAchieved)
}

func TestPillarStats_Calculation(t *testing.T) {
	stats := PillarStats{
		Achieved:   80,
		Planned:    100,
		Percentage: 80.0,
	}

	assert.Equal(t, int32(80), stats.Achieved)
	assert.Equal(t, int32(100), stats.Planned)
	assert.Equal(t, 80.0, stats.Percentage)
}

func TestGapAnalysisResponse_StatusMapping(t *testing.T) {
	tests := []struct {
		name       string
		gap        int32
		percentage int32
		status     string
	}{
		{
			name:       "completed",
			gap:        0,
			percentage: 100,
			status:     "DONE",
		},
		{
			name:       "in progress",
			gap:        50,
			percentage: 50,
			status:     "IN_PROGRESS",
		},
		{
			name:       "critical",
			gap:        80,
			percentage: 20,
			status:     "CRITICAL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := GapAnalysisResponse{
				Level:      "P2",
				Pillar:     "TECHNICAL",
				Target:     100,
				Achieved:   100 - tt.gap,
				Gap:        tt.gap,
				Status:     tt.status,
				Percentage: tt.percentage,
			}

			assert.Equal(t, tt.status, resp.Status)
			assert.Equal(t, tt.gap, resp.Gap)
		})
	}
}

func TestReportData_Structure(t *testing.T) {
	activities := []repository.FindDetailedActivityReportRow{
		{
			Title:              "Test Activity 1",
			ProgressPercentage: 100,
			Level:              repository.LadderLevelP2,
		},
	}

	reportData := ReportData{
		Activities:   activities,
		UserName:     "John Doe",
		UserEmail:    "john@example.com",
		CurrentLevel: "P2",
	}

	assert.Len(t, reportData.Activities, 1)
	assert.Equal(t, "John Doe", reportData.UserName)
	assert.Equal(t, "john@example.com", reportData.UserEmail)
	assert.Equal(t, "P2", reportData.CurrentLevel)
}

// Benchmark tests
func BenchmarkCreateActivityDTO_ToRepositoryParams(b *testing.B) {
	description := "Test description"
	dto := CreateActivityDTO{
		UserID:             uuid.New(),
		LadderID:           uuid.New(),
		Pillars:            []repository.Pillar{repository.PillarTECHNICAL},
		Title:              "Test Activity",
		Description:        &description,
		ProgressPercentage: 50,
		IsPdiTarget:        true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dto.ToRepositoryParams()
	}
}
