// Package email provides email sending functionality using Resend API
package email

import (
	"fmt"
	"log/slog"

	"github.com/me/level-up-hub/backend/config"
	"github.com/resend/resend-go/v3"
)

// Service handles email operations
type Service struct {
	cfg    *config.Config
	client *resend.Client
}

// NewService creates a new email service
func NewService(cfg *config.Config) *Service {
	var client *resend.Client
	if cfg.ResendAPIKey != "" {
		client = resend.NewClient(cfg.ResendAPIKey)
	}
	return &Service{cfg: cfg, client: client}
}

// SendReportToManager sends a PDF report to the user's manager
func (s *Service) SendReportToManager(managerName, managerEmail, userName, userEmail string, pdfContent []byte) error {
	if s.client == nil {
		slog.Warn("Resend API key not configured, simulating email send",
			slog.String("to", managerEmail),
		)
		return nil
	}

	subject := fmt.Sprintf("Relatório de Desenvolvimento de Carreira - %s", userName)

	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
				<h2 style="color: #1976d2; border-bottom: 2px solid #1976d2; padding-bottom: 10px;">
					Relatório de Desenvolvimento de Carreira
				</h2>
				
				<p>Olá <strong>%s</strong>,</p>
				
				<p>
					O colaborador <strong>%s</strong> (%s) compartilhou seu relatório de desenvolvimento de carreira com você.
				</p>
				
				<p>Este relatório contém:</p>
				
				<ul>
					<li>Resumo executivo com estatísticas de desempenho</li>
					<li>Atividades concluídas e em andamento</li>
					<li>Distribuição por níveis e pilares</li>
					<li>Evidências e documentação</li>
					<li>Progresso em relação aos objetivos do PDI</li>
				</ul>
				
				<p>O relatório completo está disponível no PDF em anexo.</p>
				
				<hr style="border: none; border-top: 1px solid #ddd; margin: 20px 0;">
				
				<p style="font-size: 12px; color: #666;">
					<em>Este é um email automático gerado pelo Level Up Hub.</em>
				</p>
			</div>
		</body>
		</html>
	`, managerName, userName, userEmail)

	filename := fmt.Sprintf("relatorio_%s.pdf", sanitizeFilename(userName))

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", s.cfg.SMTPFromName, s.cfg.SMTPFrom),
		To:      []string{managerEmail},
		Subject: subject,
		Html:    body,
		Attachments: []*resend.Attachment{
			{
				Filename: filename,
				Content:  pdfContent,
			},
		},
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		slog.Error("falha ao enviar email via Resend",
			slog.String("error", err.Error()),
			slog.String("to", managerEmail),
		)
		return fmt.Errorf("falha ao enviar email: %w", err)
	}

	slog.Info("email enviado com sucesso via Resend",
		slog.String("to", managerEmail),
		slog.String("subject", subject),
	)

	return nil
}

// sanitizeFilename removes special characters from filename
func sanitizeFilename(name string) string {
	result := ""
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			result += string(r)
		} else if r == ' ' {
			result += "_"
		}
	}
	return result
}
