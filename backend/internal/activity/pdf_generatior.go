package activity

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/me/level-up-hub/backend/internal/repository"
)

// GenerateDossierPDF cria o PDF com base nos dados do banco
func GenerateDossierPDF(activities []repository.FindDetailedActivityReportRow) (*bytes.Buffer, error) {
	// Cria um documento A4 em formato retrato (Portrait)
	m := pdf.NewMaroto(consts.Portrait, consts.A4)

	// --- CABEÇALHO ---
	m.Row(20, func() {
		m.Col(12, func() {
			m.Text("Dossiê de Evolução de Carreira", props.Text{
				Size:  18,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Relatório consolidado de atividades, pilares e evidências.", props.Text{
				Size:  10,
				Align: consts.Center,
			})
		})
	})

	// --- LISTA DE ATIVIDADES ---
	for _, act := range activities {
		// Linha de espaçamento
		m.Row(10, func() {}) 

		// Título e Nível
		m.Row(10, func() {
			m.Col(10, func() {
				m.Text(fmt.Sprintf("Atividade: %s", act.Title), props.Text{
					Style: consts.Bold,
					Size:  12,
				})
			})
			m.Col(2, func() {
				m.Text(fmt.Sprintf("[%s]", act.Level), props.Text{
					Style: consts.Bold,
					Align: consts.Right,
				})
			})
		})

		// Extraindo as Evidências do JSON do banco
		var evidences []struct {
			URL         string `json:"url"`
			Description string `json:"description"`
		}
		if act.Evidences != nil {
			json.Unmarshal(act.Evidences, &evidences)
		}

		// Listando as Evidências
		if len(evidences) > 0 {
			m.Row(8, func() {
				m.Col(12, func() {
					m.Text("Evidências anexadas:", props.Text{Size: 10, Style: consts.Italic})
				})
			})

			for _, ev := range evidences {
				m.Row(6, func() {
					m.Col(1, func() {
						m.Text("-", props.Text{Align: consts.Right})
					})
					m.Col(11, func() {
						m.Text(fmt.Sprintf("%s: %s", ev.Description, ev.URL), props.Text{
							Size: 9,
						})
					})
				})
			}
		}
	}

	buffer, err := m.Output()
	if err != nil {
		return nil, err
	}
	return &buffer, nil
}