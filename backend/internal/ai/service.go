package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/me/level-up-hub/backend/config"
	"github.com/me/level-up-hub/backend/internal/repository"
	"google.golang.org/genai"
)

type Service struct {
	client *genai.Client
	repo   *repository.Queries
}

type ClassifyRequest struct {
	Title     string `json:"title" binding:"required"`
	Execution string `json:"execution"`
}

type ClassifyResponse struct {
	Level   string   `json:"level"`
	Pillars []string `json:"pillars"`
}

func NewService(cfg *config.Config, repo *repository.Queries) *Service {
	if cfg.GeminiAPIKey == "" {
		slog.Warn("Gemini API key not configured, AI classification will not work")
		return &Service{repo: repo}
	}

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  cfg.GeminiAPIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		slog.Error("failed to create Gemini client", slog.String("error", err.Error()))
		return &Service{repo: repo}
	}

	return &Service{client: client, repo: repo}
}

func (s *Service) ClassifyTask(ctx context.Context, req ClassifyRequest) (*ClassifyResponse, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Gemini API key not configured")
	}

	// Fetch career ladder from DB
	ladders, err := s.repo.ListEntireLadder(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch career ladder: %w", err)
	}

	// Build system instruction with ladder data
	systemPrompt := buildSystemPrompt(ladders)

	// Build user prompt
	userPrompt := fmt.Sprintf("Titulo da atividade: %s", req.Title)
	if req.Execution != "" {
		userPrompt += fmt.Sprintf("\nExecucao (o que foi feito): %s", req.Execution)
	}

	// Call Gemini
	result, err := s.client.Models.GenerateContent(ctx,
		"gemini-2.0-flash",
		genai.Text(userPrompt),
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{{Text: systemPrompt}},
			},
			ResponseMIMEType: "application/json",
			Temperature:      genai.Ptr[float32](0.2),
		},
	)
	if err != nil {
		slog.Error("Gemini API call failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to classify task: %w", err)
	}

	// Parse response
	text := result.Text()
	var response ClassifyResponse
	if err := json.Unmarshal([]byte(text), &response); err != nil {
		slog.Error("failed to parse Gemini response", slog.String("text", text), slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return &response, nil
}

func buildSystemPrompt(ladders []repository.CareerLadder) string {
	var sb strings.Builder

	sb.WriteString(`Voce e um classificador de atividades de engenharia de software. 
Com base na career ladder fornecida abaixo, classifique a atividade do usuario em:
1. Um nivel da ladder (P1, P2, P3, LT1, LT2, LT3, LT4)
2. Um ou mais pilares (TECHNICAL, RESULTS, INFLUENCE)

Responda APENAS com um JSON no formato:
{"level": "P2", "pillars": ["TECHNICAL", "RESULTS"]}

Career Ladder:
`)

	for _, l := range ladders {
		sb.WriteString(fmt.Sprintf("\n--- %s (XP: %d) ---\n", l.Level, l.XpReward))
		sb.WriteString(fmt.Sprintf("Tecnico: %s\n", l.Technical))
		sb.WriteString(fmt.Sprintf("Resultados: %s\n", l.ExpectedResults))
		sb.WriteString(fmt.Sprintf("Lideranca: %s\n", l.LeadershipScope))
	}

	sb.WriteString(`
Regras:
- Escolha o nivel que melhor se encaixa com a complexidade e escopo da atividade.
- Escolha os pilares que a atividade impacta (pode ser mais de um).
- Retorne APENAS o JSON, sem explicacao.`)

	return sb.String()
}
