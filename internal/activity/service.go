package activity

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/me/level-up-hub/apperr"
	"github.com/me/level-up-hub/internal/repository"
)

type Service struct {
	repo *repository.Queries
	pool *pgxpool.Pool
}

func NewService(repo *repository.Queries, pool *pgxpool.Pool) *Service {
	return &Service{repo: repo, pool: pool}
}

func (s *Service) CreateCompleteActivity(ctx context.Context, input CreateActivityDTO) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		slog.Error("failed to begin transaction",
			slog.String("error", err.Error()),
			slog.String("user_id", input.UserID.String()),
		)
		return err
	}
	defer tx.Rollback(ctx)

	repoWithTx := s.repo.WithTx(tx)

	activity, err := repoWithTx.CreateActivity(ctx, input.ToRepositoryParams())
	if err != nil {
		slog.Error("failed to create activity",
			slog.String("error", err.Error()),
			slog.String("user_id", input.UserID.String()),
		)
		return apperr.MessageError(fmt.Sprintf(apperr.ErrCreate, apperr.ActivityPT), err)
	}

	for _, p := range input.Pillars {
		_, err = repoWithTx.CreateActivityPillar(ctx, repository.CreateActivityPillarParams{
			ActivityID: activity.ID,
			Pillar:     p,
		})
		if err != nil {
			slog.Error("failed to create activity pillar",
				slog.String("error", err.Error()),
				slog.String("user_id", input.UserID.String()),
				slog.String("activity_id", activity.ID.String()),
				slog.String("pillar", string(p)),
			)
			return apperr.MessageError(fmt.Sprintf(apperr.ErrCreate, apperr.PillarPT), err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("transaction commit failed",
            slog.String("error", err.Error()),
            slog.String("user_id", input.UserID.String()),
            slog.String("activity_id", activity.ID.String()),
        )
		return apperr.MessageError(fmt.Sprintf(apperr.ErrCreate, apperr.ActivityPT), err)
	}

	return nil
}

func (s *Service) CreateActivity(ctx context.Context, params repository.CreateActivityParams) error {
	_, err := s.repo.CreateActivity(ctx, params)
	if err != nil {
		return apperr.MessageError(fmt.Sprintf(apperr.ErrCreate, apperr.ActivityPT), err)
	}

	return nil
}

func (s *Service) AddEvidence(ctx context.Context, activityID uuid.UUID, userID uuid.UUID, url string, description string) (repository.ActivityEvidence, error) {
	_, err := s.repo.FindActivityByID(ctx, repository.FindActivityByIDParams{
		ID:     activityID,
		UserID: userID,
	})
	if err != nil {
		slog.Warn("unauthorized activity access attempt",
            slog.String("activity_id", activityID.String()),
            slog.String("user_id", userID.String()),
            slog.String("error", err.Error()),
        )
		return repository.ActivityEvidence{}, apperr.MessageError(apperr.ErrActivityNotFound, err)
	}

	evidence, err := s.repo.AddEvidence(ctx, repository.AddEvidenceParams{
		ActivityID:  activityID,
		EvidenceUrl: url,
		Description: pgtype.Text{String: description, Valid: description != ""},
	})
	if err != nil {
		slog.Error("failed to add evidence",
			slog.String("error", err.Error()),
			slog.String("activity_id", activityID.String()),
			slog.String("user_id", userID.String()),
		)
		return repository.ActivityEvidence{}, apperr.MessageError(fmt.Sprintf(apperr.ErrCreate, apperr.EvidencePT), err)
	}

	return evidence, nil
}

func (s *Service) UpdateProgress(ctx context.Context, activityID uuid.UUID, userID uuid.UUID, progress int32) error {
	if progress < 0 || progress > 100 {
		return apperr.MessageError(apperr.ErrInvalidProgress, nil)
	}

	_, err := s.repo.UpdateActivityProgress(ctx, repository.UpdateActivityProgressParams{
		ID:                 activityID,
		ProgressPercentage: progress,
		UserID:             userID,
	})

	return apperr.MessageError(fmt.Sprintf(apperr.ErrUpdate, apperr.ActivityPT), err)
}

func (s *Service) Delete(ctx context.Context, activityID uuid.UUID, userID uuid.UUID) error {
	err := s.repo.DeleteActivity(ctx, repository.DeleteActivityParams{
		ID:     activityID,
		UserID: userID,
	})
	if err != nil {
		return apperr.MessageError(fmt.Sprintf(apperr.ErrDelete, apperr.ActivityPT), err)
	}
	return nil
}

func (s *Service) GetCareerDashboard(ctx context.Context, userID uuid.UUID) (*DashboardResponse, error) {
	rows, err := s.repo.FindPdiDashboard(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := &DashboardResponse{
		PdiProgress:  make(map[string]PillarStats),
		Overdelivery: make(map[string]int32),
	}

	for _, row := range rows {
		pillarName := row.Pillar

		resp.CurrentLevel = string(row.Level)
		resp.MaxPdiXp += row.TotalPdiPlanned
		resp.TotalAchieved += row.TotalAchieved

		if row.TotalPdiPlanned > 0 {
			percentage := (float64(row.TotalAchieved) / float64(row.TotalPdiPlanned)) * 100

			resp.PdiProgress[pillarName] = PillarStats{
				Achieved:   row.TotalAchieved,
				Planned:    row.TotalPdiPlanned,
				Percentage: percentage,
			}
		}

		if row.OverdeliveryXp > 0 {
			resp.Overdelivery[pillarName] += row.OverdeliveryXp
		}
	}

	return resp, nil

}

func (s *Service) ListActivities(ctx context.Context, userID uuid.UUID) ([]repository.ListUserActivitiesRow, error) {
	return s.repo.ListUserActivities(ctx, userID)
}

func (s *Service) GetActivitiesEvidence(ctx context.Context, userID uuid.UUID) ([]repository.ListUserActivitiesWithEvidencesRow, error) {
	return s.repo.ListUserActivitiesWithEvidences(ctx, userID)
}

func (s *Service) GetDetailedReport(ctx context.Context, userID uuid.UUID) ([]repository.FindDetailedActivityReportRow, error) {
	return s.repo.FindDetailedActivityReport(ctx, userID)
}

func (s *Service) GetGapAnalysis(ctx context.Context, userID uuid.UUID, year int) ([]GapAnalysisResponse, error) {
	rows, err := s.repo.FindGapAnalysis(ctx, repository.FindGapAnalysisParams{
		UserID: userID,
		Year:   int32(year),
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
	rows, err := s.repo.FindActivityComposition(ctx, userID)
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
			VolumePercent: math.Round(volPct*100) / 100, // Arredonda para 2 casas decimais
			XpPercent:     math.Round(xpPct*100) / 100,
		})
	}

	return radar, nil
}

func (s *Service) GetCycleComparison(ctx context.Context, userID uuid.UUID) (*ComparisonReport, error) {
	// 1. Descobre os ciclos no banco
	current, _ := s.repo.FindCurrentCycle(ctx)
	previous, _ := s.repo.FindPreviousCycle(ctx, current.StartDate)

	// 2. Busca performance de ambos
	currentPerf, _ := s.repo.FindPerformanceByPeriod(ctx, repository.FindPerformanceByPeriodParams{
		UserID: userID, CompletedAt: current.StartDate, CompletedAt_2: current.EndDate,
	})
	prevPerf, _ := s.repo.FindPerformanceByPeriod(ctx, repository.FindPerformanceByPeriodParams{
		UserID: userID, CompletedAt: previous.StartDate, CompletedAt_2: previous.EndDate,
	})

	// 3. Lógica de De-Para (Map para busca rápida)
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
