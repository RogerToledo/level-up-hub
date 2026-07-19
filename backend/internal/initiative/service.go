package initiative

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/google/uuid"
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

func (s *Service) Create(ctx context.Context, input CreateInitiativeDTO) error {
	_, err := s.repo.CreateInitiative(ctx, input.ToRepositoryParams())
	if err != nil {
		slog.Error("failed to create initiative", slog.String("error", err.Error()))
		return apperr.MessageError(fmt.Sprintf(apperr.ErrCreate, "iniciativa"), err)
	}
	return nil
}

func (s *Service) Update(ctx context.Context, initiativeID uuid.UUID, userID uuid.UUID, dto UpdateInitiativeDTO) error {
	params := dto.ToRepositoryParams(initiativeID, userID)
	_, err := s.repo.UpdateInitiative(ctx, params)
	if err != nil {
		slog.Error("failed to update initiative", slog.String("error", err.Error()))
		return apperr.MessageError(fmt.Sprintf(apperr.ErrUpdate, "iniciativa"), err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, initiativeID uuid.UUID, userID uuid.UUID) error {
	err := s.repo.DeleteInitiative(ctx, repository.DeleteInitiativeParams{
		ID:     initiativeID,
		UserID: userID,
	})
	if err != nil {
		return apperr.MessageError(fmt.Sprintf(apperr.ErrDelete, "iniciativa"), err)
	}
	return nil
}

func (s *Service) ListInitiatives(ctx context.Context, userID uuid.UUID) ([]repository.ListUserInitiativesRow, error) {
	return s.repo.ListUserInitiatives(ctx, userID)
}

func (s *Service) GetCareerDashboard(ctx context.Context, userID uuid.UUID) (*DashboardResponse, error) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	rows, err := s.repo.FindPdiDashboard(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := &DashboardResponse{
		OfficialLevel: user.CurrentLevel,
		PdiProgress:   make(map[string]PillarStats),
		Overdelivery:  make(map[string]int32),
	}

	var highestTarget repository.LadderLevel = ""
	for _, row := range rows {
		if string(row.Level) > string(highestTarget) {
			highestTarget = row.Level
		}

		percentage := float64(0)
		if row.TotalPdiPlanned > 0 {
			percentage = (float64(row.TotalAchieved) / float64(row.TotalPdiPlanned)) * 100
		}

		resp.MaxPdiXp += row.TotalPdiPlanned
		resp.TotalAchieved += row.TotalAchieved
		resp.PdiProgress[row.Pillar] = PillarStats{
			Achieved:   row.TotalAchieved,
			Planned:    row.TotalPdiPlanned,
			Percentage: percentage,
		}
		levelKey := string(row.Level)
		resp.Overdelivery[levelKey] += row.OverdeliveryXp
	}

	resp.TargetLevel = highestTarget
	return resp, nil
}

func (s *Service) GetDetailedReport(ctx context.Context, userID uuid.UUID) ([]repository.FindDetailedInitiativeReportRow, error) {
	return s.repo.FindDetailedInitiativeReport(ctx, userID)
}

func (s *Service) GetDetailedReportData(ctx context.Context, userID uuid.UUID) (ReportData, error) {
	initiatives, err := s.repo.FindDetailedInitiativeReport(ctx, userID)
	if err != nil {
		return ReportData{}, err
	}

	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return ReportData{
			Initiatives: initiatives,
			UserName:    "Colaborador",
		}, nil
	}

	return ReportData{
		Initiatives:  initiatives,
		UserName:     user.Username,
		UserEmail:    user.Email,
		CurrentLevel: string(user.CurrentLevel),
	}, nil
}

func (s *Service) GetGapAnalysis(ctx context.Context, userID uuid.UUID, year int) ([]GapAnalysisResponse, error) {
	rows, err := s.repo.FindGapAnalysis(ctx, repository.FindGapAnalysisParams{
		UserID:  userID,
		Column2: int32(year),
	})
	if err != nil {
		return nil, err
	}

	var analysis []GapAnalysisResponse
	for _, row := range rows {
		status := "IN_PROGRESS"
		if row.GapXp <= 0 {
			status = "DONE"
		} else if row.CompletionPercentage < 30 {
			status = "CRITICAL"
		}

		analysis = append(analysis, GapAnalysisResponse{
			Level:      string(row.Level),
			Pillar:     row.Pillar,
			Target:     row.TargetXp,
			Achieved:   row.AchievedXp,
			Gap:        row.GapXp,
			Status:     status,
			Percentage: row.CompletionPercentage,
		})
	}
	return analysis, nil
}

func (s *Service) GetCareerRadar(ctx context.Context, userID uuid.UUID) (*CareerRadar, error) {
	rows, err := s.repo.FindInitiativeComposition(ctx, userID)
	if err != nil {
		return nil, err
	}

	radar := &CareerRadar{
		Breakdown: make([]LevelComposition, 0),
	}

	for _, row := range rows {
		radar.TotalActivities += row.TotalActivities
		radar.TotalXP += row.TotalXp
	}

	for _, row := range rows {
		volPct := 0.0
		xpPct := 0.0
		if radar.TotalActivities > 0 {
			volPct = (float64(row.TotalActivities) / float64(radar.TotalActivities)) * 100
		}
		if radar.TotalXP > 0 {
			xpPct = (float64(row.TotalXp) / float64(radar.TotalXP)) * 100
		}

		radar.Breakdown = append(radar.Breakdown, LevelComposition{
			LevelName:     string(row.Level),
			ActivityCount: row.TotalActivities,
			TotalXP:       row.TotalXp,
			VolumePercent: math.Round(volPct*100) / 100,
			XpPercent:     math.Round(xpPct*100) / 100,
		})
	}

	return radar, nil
}

func (s *Service) GetCycleComparison(ctx context.Context, userID uuid.UUID) (*ComparisonReport, error) {
	current, _ := s.repo.FindCurrentCycle(ctx)
	previous, _ := s.repo.FindPreviousCycle(ctx, current.StartDate)

	currentPerf, _ := s.repo.FindPerformanceByPeriod(ctx, repository.FindPerformanceByPeriodParams{
		UserID: userID, CompletedAt: current.StartDate, CompletedAt_2: current.EndDate,
	})
	prevPerf, _ := s.repo.FindPerformanceByPeriod(ctx, repository.FindPerformanceByPeriodParams{
		UserID: userID, CompletedAt: previous.StartDate, CompletedAt_2: previous.EndDate,
	})

	prevMap := make(map[string]int32)
	var totalPrevXP int32
	for _, p := range prevPerf {
		prevMap[string(p.Level)] = p.TotalXp
		totalPrevXP += p.TotalXp
	}

	var totalCurrXP int32
	report := &ComparisonReport{
		CurrentCycleName:  current.Name,
		PreviousCycleName: previous.Name,
		LevelEvolution:    []LevelComparison{},
	}

	for _, c := range currentPerf {
		totalCurrXP += c.TotalXp
		prevXP := prevMap[string(c.Level)]
		report.LevelEvolution = append(report.LevelEvolution, LevelComparison{
			LevelName: string(c.Level),
			CurrentXP: c.TotalXp,
			PrevXP:    prevXP,
			Diff:      c.TotalXp - prevXP,
		})
	}

	report.GrowthXP = totalCurrXP - totalPrevXP
	if totalPrevXP > 0 {
		report.PercentChange = (float64(report.GrowthXP) / float64(totalPrevXP)) * 100
	}

	return report, nil
}
