package activity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/me/level-up-hub/backend/internal/repository"
)

// ReportData encapsula os dados necessários para gerar o PDF
type ReportData struct {
	Activities   []repository.FindDetailedActivityReportRow
	UserName     string
	UserEmail    string
	CurrentLevel string
}

// Cores do tema
var (
	primaryColor   = color.Color{Red: 25, Green: 118, Blue: 210}   // Azul
	secondaryColor = color.Color{Red: 67, Green: 160, Blue: 71}    // Verde
	warningColor   = color.Color{Red: 251, Green: 140, Blue: 0}    // Laranja
	lightGray      = color.Color{Red: 245, Green: 245, Blue: 245}  // Cinza claro
	darkGray       = color.Color{Red: 97, Green: 97, Blue: 97}     // Cinza escuro
)

// Evidence representa uma evidência extraída do JSON
type Evidence struct {
	URL         string    `json:"url"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// GenerateDossierPDF cria o PDF com base nos dados do banco
func GenerateDossierPDF(activities []repository.FindDetailedActivityReportRow) (*bytes.Buffer, error) {
	data := ReportData{
		Activities: activities,
		UserName:   "Colaborador", // Será substituído quando tivermos acesso ao usuário
		UserEmail:  "",
	}
	return GenerateDetailedDossierPDF(data)
}

// GenerateDetailedDossierPDF cria um PDF profissional e completo
func GenerateDetailedDossierPDF(data ReportData) (*bytes.Buffer, error) {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)

	// Página de Capa
	buildCoverPage(m, data)

	// Página de Resumo Executivo
	buildExecutiveSummary(m, data.Activities)

	// Detalhamento das Atividades
	buildActivitiesSection(m, data.Activities)

	// Rodapé em todas as páginas
	m.SetPageMargins(10, 15, 15)

	buffer, err := m.Output()
	if err != nil {
		return nil, err
	}
	return &buffer, nil
}

// buildCoverPage cria a página de capa do relatório
func buildCoverPage(m pdf.Maroto, data ReportData) {
	m.Row(60, func() {
		m.Col(12, func() {
			m.Text("", props.Text{})
		})
	})

	// Título principal
	m.Row(25, func() {
		m.Col(12, func() {
			m.Text("Dossiê de Evolução de Carreira", props.Text{
				Size:  24,
				Style: consts.Bold,
				Align: consts.Center,
				Color: primaryColor,
			})
		})
	})

	m.Row(15, func() {
		m.Col(12, func() {
			m.Text("Relatório Consolidado de Desempenho", props.Text{
				Size:  14,
				Align: consts.Center,
				Color: darkGray,
			})
		})
	})

	m.Row(40, func() {})

	// Informações do colaborador (se disponíveis)
	if data.UserName != "" {
		m.Row(12, func() {
			m.ColSpace(3)
			m.Col(6, func() {
				m.Text("Colaborador:", props.Text{
					Size:  11,
					Style: consts.Bold,
					Align: consts.Center,
				})
			})
			m.ColSpace(3)
		})
		m.Row(10, func() {
			m.ColSpace(3)
			m.Col(6, func() {
				m.Text(data.UserName, props.Text{
					Size:  12,
					Align: consts.Center,
					Color: primaryColor,
				})
			})
			m.ColSpace(3)
		})
	}

	if data.UserEmail != "" {
		m.Row(10, func() {
			m.ColSpace(3)
			m.Col(6, func() {
				m.Text(data.UserEmail, props.Text{
					Size:  9,
					Align: consts.Center,
					Color: darkGray,
				})
			})
			m.ColSpace(3)
		})
	}

	m.Row(30, func() {})

	// Data de geração
	m.Row(10, func() {
		m.ColSpace(3)
		m.Col(6, func() {
			m.Text(fmt.Sprintf("Gerado em: %s", time.Now().Format("02/01/2006 às 15:04")), props.Text{
				Size:  10,
				Align: consts.Center,
				Color: darkGray,
			})
		})
		m.ColSpace(3)
	})

	m.Row(20, func() {
		m.ColSpace(3)
		m.Col(6, func() {
			m.Text(fmt.Sprintf("Total de Atividades: %d", len(data.Activities)), props.Text{
				Size:  11,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
		m.ColSpace(3)
	})
}

// buildExecutiveSummary cria a página de resumo com estatísticas
func buildExecutiveSummary(m pdf.Maroto, activities []repository.FindDetailedActivityReportRow) {
	m.AddPage()

	// Título da seção
	m.Row(15, func() {
		m.Col(12, func() {
			m.Text("Resumo Executivo", props.Text{
				Size:  18,
				Style: consts.Bold,
				Color: primaryColor,
			})
		})
	})

	m.Row(2, func() {
		m.Col(12, func() {
			m.Line(1, props.Line{Color: primaryColor})
		})
	})

	m.Row(8, func() {})

	// Calcular estatísticas
	stats := calculateStatistics(activities)

	// Cards de estatísticas
	m.Row(25, func() {
		m.Col(4, func() {
			buildStatCard(m, "Total de Atividades", fmt.Sprintf("%d", stats.TotalActivities), primaryColor)
		})
		m.Col(4, func() {
			buildStatCard(m, "Atividades Concluídas", fmt.Sprintf("%d", stats.CompletedActivities), secondaryColor)
		})
		m.Col(4, func() {
			buildStatCard(m, "XP Total Conquistado", fmt.Sprintf("%d XP", stats.TotalXP), warningColor)
		})
	})

	m.Row(5, func() {})

	// Distribuição por nível
	m.Row(12, func() {
		m.Col(12, func() {
			m.Text("Distribuição por Nível", props.Text{
				Size:  12,
				Style: consts.Bold,
			})
		})
	})

	m.Row(8, func() {
		m.Col(3, func() {
			m.Text("Nível", props.Text{
				Size:  10,
				Style: consts.Bold,
				Color: darkGray,
			})
		})
		m.Col(3, func() {
			m.Text("Quantidade", props.Text{
				Size:  10,
				Style: consts.Bold,
				Align: consts.Center,
				Color: darkGray,
			})
		})
		m.Col(3, func() {
			m.Text("XP Total", props.Text{
				Size:  10,
				Style: consts.Bold,
				Align: consts.Center,
				Color: darkGray,
			})
		})
		m.Col(3, func() {
			m.Text("% Conclusão", props.Text{
				Size:  10,
				Style: consts.Bold,
				Align: consts.Center,
				Color: darkGray,
			})
		})
	})

	m.Row(1, func() {
		m.Col(12, func() {
			m.Line(0.5, props.Line{Color: lightGray})
		})
	})

	for level, data := range stats.ByLevel {
		m.Row(8, func() {
			m.Col(3, func() {
				m.Text(string(level), props.Text{
					Size: 10,
				})
			})
			m.Col(3, func() {
				m.Text(fmt.Sprintf("%d", data.Count), props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
			m.Col(3, func() {
				m.Text(fmt.Sprintf("%d", data.TotalXP), props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
			m.Col(3, func() {
				completionRate := 0.0
				if data.Count > 0 {
					completionRate = (float64(data.CompletedCount) / float64(data.Count)) * 100
				}
				m.Text(fmt.Sprintf("%.0f%%", completionRate), props.Text{
					Size:  10,
					Align: consts.Center,
					Color: getColorByCompletion(completionRate),
				})
			})
		})
	}

	m.Row(10, func() {})

	// Distribuição por Pilares
	m.Row(12, func() {
		m.Col(12, func() {
			m.Text("Distribuição por Pilares", props.Text{
				Size:  12,
				Style: consts.Bold,
			})
		})
	})

	pillarCounts := calculatePillarDistribution(activities)
	
	m.Row(8, func() {
		m.Col(6, func() {
			m.Text("Pilar", props.Text{
				Size:  10,
				Style: consts.Bold,
				Color: darkGray,
			})
		})
		m.Col(6, func() {
			m.Text("Atividades", props.Text{
				Size:  10,
				Style: consts.Bold,
				Align: consts.Center,
				Color: darkGray,
			})
		})
	})

	m.Row(1, func() {
		m.Col(12, func() {
			m.Line(0.5, props.Line{Color: lightGray})
		})
	})

	for pillar, count := range pillarCounts {
		m.Row(8, func() {
			m.Col(6, func() {
				m.Text(getPillarName(pillar), props.Text{
					Size: 10,
				})
			})
			m.Col(6, func() {
				m.Text(fmt.Sprintf("%d", count), props.Text{
					Size:  10,
					Align: consts.Center,
				})
			})
		})
	}

	m.Row(10, func() {})

	// Atividades PDI
	pdiCount := 0
	for _, act := range activities {
		if act.IsPdiTarget {
			pdiCount++
		}
	}

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("🎯 Atividades marcadas como PDI: %d de %d (%.0f%%)",
				pdiCount,
				len(activities),
				float64(pdiCount)/float64(len(activities))*100,
			), props.Text{
				Size:  10,
				Style: consts.Bold,
			})
		})
	})
}

// buildStatCard cria um card de estatística
func buildStatCard(m pdf.Maroto, title, value string, bgColor color.Color) {
	m.Text(title, props.Text{
		Size:  8,
		Color: darkGray,
	})
	m.Text(value, props.Text{
		Top:   5,
		Size:  14,
		Style: consts.Bold,
		Color: bgColor,
	})
}

// buildActivitiesSection cria a seção detalhada de atividades
func buildActivitiesSection(m pdf.Maroto, activities []repository.FindDetailedActivityReportRow) {
	m.AddPage()

	// Título da seção
	m.Row(15, func() {
		m.Col(12, func() {
			m.Text("Detalhamento das Atividades", props.Text{
				Size:  18,
				Style: consts.Bold,
				Color: primaryColor,
			})
		})
	})

	m.Row(2, func() {
		m.Col(12, func() {
			m.Line(1, props.Line{Color: primaryColor})
		})
	})

	m.Row(8, func() {})

	// Listar cada atividade
	for idx, act := range activities {
		// Cabeçalho da atividade com fundo colorido
		m.Row(12, func() {
			m.Col(12, func() {
				m.Text(fmt.Sprintf("%d. %s", idx+1, act.Title), props.Text{
					Size:  11,
					Style: consts.Bold,
					Color: primaryColor,
				})
			})
		})

		// Informações principais em tabela
		m.Row(8, func() {
			m.Col(3, func() {
				m.Text("Nível:", props.Text{
					Size:  9,
					Style: consts.Bold,
				})
			})
			m.Col(3, func() {
				m.Text(string(act.Level), props.Text{
					Size: 9,
				})
			})
			m.Col(3, func() {
				m.Text("XP:", props.Text{
					Size:  9,
					Style: consts.Bold,
				})
			})
			m.Col(3, func() {
				m.Text(fmt.Sprintf("%d pontos", act.XpReward), props.Text{
					Size:  9,
					Color: warningColor,
				})
			})
		})

		m.Row(8, func() {
			m.Col(3, func() {
				m.Text("Progresso:", props.Text{
					Size:  9,
					Style: consts.Bold,
				})
			})
			m.Col(3, func() {
				progressColor := getColorByProgress(act.ProgressPercentage)
				m.Text(fmt.Sprintf("%d%%", act.ProgressPercentage), props.Text{
					Size:  9,
					Color: progressColor,
					Style: consts.Bold,
				})
			})
			m.Col(3, func() {
				m.Text("PDI:", props.Text{
					Size:  9,
					Style: consts.Bold,
				})
			})
			m.Col(3, func() {
				pdiText := "Não"
				if act.IsPdiTarget {
					pdiText = "Sim 🎯"
				}
				m.Text(pdiText, props.Text{
					Size: 9,
				})
			})
		})

		// Pilares
		var pillars []string
		if act.Pillars != nil {
			json.Unmarshal([]byte(fmt.Sprintf("%v", act.Pillars)), &pillars)
		}

		if len(pillars) > 0 {
			m.Row(8, func() {
				m.Col(3, func() {
					m.Text("Pilares:", props.Text{
						Size:  9,
						Style: consts.Bold,
					})
				})
				m.Col(9, func() {
					pillarNames := ""
					for i, p := range pillars {
						if i > 0 {
							pillarNames += ", "
						}
						pillarNames += getPillarName(p)
					}
					m.Text(pillarNames, props.Text{
						Size: 9,
					})
				})
			})
		}

		// Evidências
		var evidences []Evidence
		if act.Evidences != nil && len(act.Evidences) > 0 {
			json.Unmarshal(act.Evidences, &evidences)
		}

		if len(evidences) > 0 {
			m.Row(8, func() {
				m.Col(12, func() {
					m.Text(fmt.Sprintf("Evidências (%d):", len(evidences)), props.Text{
						Size:  9,
						Style: consts.Bold,
					})
				})
			})

			for i, ev := range evidences {
				m.Row(7, func() {
					m.Col(1, func() {
						m.Text(fmt.Sprintf("%d.", i+1), props.Text{
							Size:  8,
							Align: consts.Right,
						})
					})
					m.Col(11, func() {
						m.Text(ev.Description, props.Text{
							Size: 8,
						})
					})
				})
				if ev.URL != "" {
					m.Row(6, func() {
						m.ColSpace(1)
						m.Col(11, func() {
							m.Text(ev.URL, props.Text{
								Size:  7,
								Color: primaryColor,
							})
						})
					})
				}
			}
		} else {
			m.Row(7, func() {
				m.Col(12, func() {
					m.Text("Sem evidências cadastradas", props.Text{
						Size:  8,
						Style: consts.Italic,
						Color: darkGray,
					})
				})
			})
		}

		// Linha separadora entre atividades
		m.Row(8, func() {})
		m.Row(1, func() {
			m.Col(12, func() {
				m.Line(0.3, props.Line{Color: lightGray})
			})
		})
		m.Row(8, func() {})
	}

	// Se não há atividades
	if len(activities) == 0 {
		m.Row(20, func() {
			m.Col(12, func() {
				m.Text("Nenhuma atividade registrada", props.Text{
					Size:  12,
					Align: consts.Center,
					Color: darkGray,
					Style: consts.Italic,
				})
			})
		})
	}
}

// Funções auxiliares para estatísticas

type Statistics struct {
	TotalActivities     int
	CompletedActivities int
	InProgressCount     int
	TotalXP             int32
	ByLevel             map[repository.LadderLevel]LevelStats
}

type LevelStats struct {
	Count          int
	CompletedCount int
	TotalXP        int32
}

func calculateStatistics(activities []repository.FindDetailedActivityReportRow) Statistics {
	stats := Statistics{
		ByLevel: make(map[repository.LadderLevel]LevelStats),
	}

	for _, act := range activities {
		stats.TotalActivities++

		if act.ProgressPercentage == 100 {
			stats.CompletedActivities++
			stats.TotalXP += act.XpReward
		} else if act.ProgressPercentage > 0 {
			stats.InProgressCount++
		}

		levelData := stats.ByLevel[act.Level]
		levelData.Count++
		levelData.TotalXP += act.XpReward
		if act.ProgressPercentage == 100 {
			levelData.CompletedCount++
		}
		stats.ByLevel[act.Level] = levelData
	}

	return stats
}

func calculatePillarDistribution(activities []repository.FindDetailedActivityReportRow) map[string]int {
	pillarCounts := make(map[string]int)

	for _, act := range activities {
		var pillars []string
		if act.Pillars != nil {
			json.Unmarshal([]byte(fmt.Sprintf("%v", act.Pillars)), &pillars)
			for _, p := range pillars {
				pillarCounts[p]++
			}
		}
	}

	return pillarCounts
}

func getPillarName(pillar string) string {
	switch pillar {
	case "TECHNICAL":
		return "Técnico"
	case "RESULTS":
		return "Resultados"
	case "INFLUENCE":
		return "Influência"
	default:
		return pillar
	}
}

func getColorByProgress(progress int32) color.Color {
	if progress == 100 {
		return secondaryColor // Verde
	} else if progress >= 50 {
		return warningColor // Laranja
	}
	return darkGray // Cinza
}

func getColorByCompletion(completionRate float64) color.Color {
	if completionRate >= 75 {
		return secondaryColor // Verde
	} else if completionRate >= 50 {
		return warningColor // Laranja
	}
	return darkGray // Cinza
}