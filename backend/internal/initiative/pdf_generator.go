package initiative

import (
	"bytes"
	"fmt"
	"time"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/me/level-up-hub/backend/internal/repository"
)

// ReportData encapsula os dados necessarios para gerar o PDF
type ReportData struct {
	Initiatives  []repository.FindDetailedInitiativeReportRow
	UserName     string
	UserEmail    string
	CurrentLevel string
}

var (
	primaryColor   = color.Color{Red: 25, Green: 118, Blue: 210}
	secondaryColor = color.Color{Red: 67, Green: 160, Blue: 71}
	warningColor   = color.Color{Red: 251, Green: 140, Blue: 0}
	lightGray      = color.Color{Red: 245, Green: 245, Blue: 245}
	darkGray       = color.Color{Red: 97, Green: 97, Blue: 97}
)

func GenerateDetailedDossierPDF(data ReportData) (*bytes.Buffer, error) {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)

	buildCoverPage(m, data)
	buildExecutiveSummary(m, data.Initiatives)
	buildInitiativesSection(m, data.Initiatives)

	m.SetPageMargins(10, 15, 15)

	buffer, err := m.Output()
	if err != nil {
		return nil, err
	}
	return &buffer, nil
}

func buildCoverPage(m pdf.Maroto, data ReportData) {
	m.Row(60, func() {
		m.Col(12, func() { m.Text("", props.Text{}) })
	})

	m.Row(25, func() {
		m.Col(12, func() {
			m.Text("Relatório de Evoluç˜o de Carreira", props.Text{
				Size: 24, Style: consts.Bold, Align: consts.Center, Color: primaryColor,
			})
		})
	})

	m.Row(15, func() {
		m.Col(12, func() {
			m.Text("Relatorio Consolidado de Desempenho", props.Text{
				Size: 14, Align: consts.Center, Color: darkGray,
			})
		})
	})

	m.Row(40, func() {})

	if data.UserName != "" {
		m.Row(12, func() {
			m.ColSpace(3)
			m.Col(6, func() {
				m.Text("Colaborador:", props.Text{Size: 11, Style: consts.Bold, Align: consts.Center})
			})
			m.ColSpace(3)
		})
		m.Row(10, func() {
			m.ColSpace(3)
			m.Col(6, func() {
				m.Text(data.UserName, props.Text{Size: 12, Align: consts.Center, Color: primaryColor})
			})
			m.ColSpace(3)
		})
	}

	if data.UserEmail != "" {
		m.Row(10, func() {
			m.ColSpace(3)
			m.Col(6, func() {
				m.Text(data.UserEmail, props.Text{Size: 9, Align: consts.Center, Color: darkGray})
			})
			m.ColSpace(3)
		})
	}

	m.Row(30, func() {})

	m.Row(10, func() {
		m.ColSpace(3)
		m.Col(6, func() {
			m.Text(fmt.Sprintf("Gerado em: %s", time.Now().Format("02/01/2006 as 15:04")), props.Text{
				Size: 10, Align: consts.Center, Color: darkGray,
			})
		})
		m.ColSpace(3)
	})

	m.Row(20, func() {
		m.ColSpace(3)
		m.Col(6, func() {
			m.Text(fmt.Sprintf("Total de Iniciativas: %d", len(data.Initiatives)), props.Text{
				Size: 11, Style: consts.Bold, Align: consts.Center,
			})
		})
		m.ColSpace(3)
	})
}

func buildExecutiveSummary(m pdf.Maroto, initiatives []repository.FindDetailedInitiativeReportRow) {
	m.AddPage()

	m.Row(15, func() {
		m.Col(12, func() {
			m.Text("Resumo Executivo", props.Text{Size: 18, Style: consts.Bold, Color: primaryColor})
		})
	})

	m.Row(2, func() {
		m.Col(12, func() { m.Line(1, props.Line{Color: primaryColor}) })
	})

	m.Row(8, func() {})

	stats := calculateStatistics(initiatives)

	m.Row(25, func() {
		m.Col(4, func() {
			m.Text("Total de Iniciativas", props.Text{Size: 8, Color: darkGray})
			m.Text(fmt.Sprintf("%d", stats.TotalCount), props.Text{Top: 5, Size: 14, Style: consts.Bold, Color: primaryColor})
		})
		m.Col(4, func() {
			m.Text("Concluidas", props.Text{Size: 8, Color: darkGray})
			m.Text(fmt.Sprintf("%d", stats.CompletedCount), props.Text{Top: 5, Size: 14, Style: consts.Bold, Color: secondaryColor})
		})
		m.Col(4, func() {
			m.Text("XP Total Conquistado", props.Text{Size: 8, Color: darkGray})
			m.Text(fmt.Sprintf("%d XP", stats.TotalXP), props.Text{Top: 5, Size: 14, Style: consts.Bold, Color: warningColor})
		})
	})
}

func buildInitiativesSection(m pdf.Maroto, initiatives []repository.FindDetailedInitiativeReportRow) {
	m.AddPage()

	m.Row(15, func() {
		m.Col(12, func() {
			m.Text("Detalhamento das Tasks", props.Text{Size: 18, Style: consts.Bold, Color: primaryColor})
		})
	})

	m.Row(2, func() {
		m.Col(12, func() { m.Line(1, props.Line{Color: primaryColor}) })
	})

	m.Row(8, func() {})

	for idx, init := range initiatives {
		// Title
		m.Row(12, func() {
			m.Col(12, func() {
				m.Text(fmt.Sprintf("%d. %s", idx+1, init.Title), props.Text{
					Size: 11, Style: consts.Bold, Color: primaryColor,
				})
			})
		})

		// Level, XP, Progress, PDI
		m.Row(8, func() {
			m.Col(3, func() { m.Text("Nivel:", props.Text{Size: 9, Style: consts.Bold}) })
			m.Col(3, func() { m.Text(string(init.Level), props.Text{Size: 9}) })
			m.Col(3, func() { m.Text("XP:", props.Text{Size: 9, Style: consts.Bold}) })
			m.Col(3, func() {
				m.Text(fmt.Sprintf("%d pontos", init.XpReward), props.Text{Size: 9, Color: warningColor})
			})
		})

		m.Row(8, func() {
			m.Col(3, func() { m.Text("Progresso:", props.Text{Size: 9, Style: consts.Bold}) })
			m.Col(3, func() {
				m.Text(fmt.Sprintf("%d%%", init.ProgressPercentage), props.Text{
					Size: 9, Style: consts.Bold, Color: getColorByProgress(init.ProgressPercentage),
				})
			})
			m.Col(3, func() { m.Text("PDI:", props.Text{Size: 9, Style: consts.Bold}) })
			m.Col(3, func() {
				pdiText := "Nao"
				if init.IsPdiTarget {
					pdiText = "Sim"
				}
				m.Text(pdiText, props.Text{Size: 9})
			})
		})

		// Execution
		if init.Execution.Valid && init.Execution.String != "" {
			m.Row(4, func() {})
			m.Row(8, func() {
				m.Col(12, func() {
					m.Text("Execucao:", props.Text{Size: 9, Style: consts.Bold})
				})
			})
			m.Row(8, func() {
				m.Col(12, func() {
					m.Text(init.Execution.String, props.Text{Size: 9, Color: darkGray})
				})
			})
		}

		// Impact Summary
		if init.ImpactSummary.Valid && init.ImpactSummary.String != "" {
			m.Row(4, func() {})
			m.Row(8, func() {
				m.Col(12, func() {
					m.Text("Impacto:", props.Text{Size: 9, Style: consts.Bold})
				})
			})
			m.Row(8, func() {
				m.Col(12, func() {
					m.Text(init.ImpactSummary.String, props.Text{Size: 9, Color: darkGray})
				})
			})
		}

		// Separator
		m.Row(8, func() {})
		m.Row(1, func() {
			m.Col(12, func() { m.Line(0.3, props.Line{Color: lightGray}) })
		})
		m.Row(8, func() {})
	}

	if len(initiatives) == 0 {
		m.Row(20, func() {
			m.Col(12, func() {
				m.Text("Nenhuma task registrada", props.Text{
					Size: 12, Align: consts.Center, Color: darkGray, Style: consts.Italic,
				})
			})
		})
	}
}

type statistics struct {
	TotalCount     int
	CompletedCount int
	TotalXP        int32
}

func calculateStatistics(initiatives []repository.FindDetailedInitiativeReportRow) statistics {
	var stats statistics
	for _, init := range initiatives {
		stats.TotalCount++
		if init.ProgressPercentage == 100 {
			stats.CompletedCount++
			stats.TotalXP += init.XpReward
		}
	}
	return stats
}

func getColorByProgress(progress int32) color.Color {
	if progress == 100 {
		return secondaryColor
	} else if progress >= 50 {
		return warningColor
	}
	return darkGray
}

// suppress unused import
var _ = lightGray
