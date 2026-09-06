package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/me/level-up-hub/backend/config"
	"github.com/me/level-up-hub/backend/internal/repository"
)

type Service struct {
	apiKey string
	repo   *repository.Queries
}

type ClassifyRequest struct {
	Title         string `json:"title" binding:"required"`
	Execution     string `json:"execution"`
	ImpactSummary string `json:"impact_summary"`
}

type ClassifyResponse struct {
	Level   string   `json:"level"`
	Pillars []string `json:"pillars"`
}

type openRouterRequest struct {
	Model    string          `json:"model"`
	Messages []openRouterMsg `json:"messages"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

type openRouterMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewService(cfg *config.Config, repo *repository.Queries) *Service {
	// Check both OPENROUTER_API_KEY and OPENAI_API_KEY
	apiKey := cfg.OpenRouterAPIKey
	if apiKey == "" {
		apiKey = cfg.OpenAIAPIKey
	}

	if apiKey == "" {
		slog.Warn("OpenRouter API key not configured, AI classification will not work")
		return &Service{repo: repo}
	}

	return &Service{apiKey: apiKey, repo: repo}
}

func (s *Service) ClassifyTask(ctx context.Context, req ClassifyRequest) (*ClassifyResponse, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("OpenRouter API key not configured")
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
	if req.ImpactSummary != "" {
		userPrompt += fmt.Sprintf("\nResumo do Impacto: %s", req.ImpactSummary)
	}

	// Build OpenRouter request
	reqBody := openRouterRequest{
		Model: "xiaomi/mimo-v2.5",
		Messages: []openRouterMsg{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		ResponseFormat: &responseFormat{Type: "json_object"},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Call OpenRouter API
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://leveluphub.com")
	httpReq.Header.Set("X-Title", "Level Up Hub")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		slog.Error("OpenRouter API call failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to call AI API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("OpenRouter API error", slog.Int("status", resp.StatusCode), slog.String("body", string(body)))
		return nil, fmt.Errorf("AI API returned status %d", resp.StatusCode)
	}

	// Parse response
	var openResp openRouterResponse
	if err := json.Unmarshal(body, &openResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(openResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	content := openResp.Choices[0].Message.Content

	// Parse the JSON response
	var response ClassifyResponse
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		slog.Error("failed to parse AI response", slog.String("text", content), slog.String("error", err.Error()))
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