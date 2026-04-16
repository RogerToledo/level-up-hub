// Package email provides email sending functionality using SMTP
package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/smtp"

	"github.com/me/level-up-hub/backend/config"
)

// Service handles email operations
type Service struct {
	cfg *config.Config
}

// NewService creates a new email service
func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

// EmailData contains the data needed to send an email
type EmailData struct {
	To         string
	Subject    string
	Body       string
	Attachment *Attachment
}

// Attachment represents an email attachment
type Attachment struct {
	Filename string
	Content  []byte
	MimeType string
}

// SendEmail sends an email with optional attachment
func (s *Service) SendEmail(data EmailData) error {
	if s.cfg.SMTPUser == "" || s.cfg.SMTPPassword == "" {
		slog.Warn("SMTP credentials not configured, simulating email send",
			slog.String("to", data.To),
			slog.String("subject", data.Subject),
		)
		// Simula envio bem-sucedido para testes
		slog.Info("email simulated successfully (SMTP not configured)",
			slog.String("to", data.To),
			slog.String("subject", data.Subject),
		)
		return nil
	}

	// Setup authentication
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	// Build message
	var message bytes.Buffer
	boundary := "boundary-level-up-hub"

	// Headers
	message.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.cfg.SMTPFromName, s.cfg.SMTPFrom))
	message.WriteString(fmt.Sprintf("To: %s\r\n", data.To))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", data.Subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", boundary))
	message.WriteString("\r\n")

	// Body
	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	message.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	message.WriteString("Content-Transfer-Encoding: 7bit\r\n")
	message.WriteString("\r\n")
	message.WriteString(data.Body)
	message.WriteString("\r\n")

	// Attachment
	if data.Attachment != nil {
		message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		message.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", data.Attachment.MimeType, data.Attachment.Filename))
		message.WriteString("Content-Transfer-Encoding: base64\r\n")
		message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", data.Attachment.Filename))
		message.WriteString("\r\n")

		// Encode attachment in base64
		encoded := encodeBase64(data.Attachment.Content)
		message.WriteString(encoded)
		message.WriteString("\r\n")
	}

	message.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	// Send email via SMTP
	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)

	// Connect to SMTP server (without TLS first for port 587)
	client, err := smtp.Dial(addr)
	if err != nil {
		slog.Error("failed to connect to SMTP server",
			slog.String("error", err.Error()),
			slog.String("host", s.cfg.SMTPHost),
			slog.Int("port", s.cfg.SMTPPort),
		)
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	// Start TLS (STARTTLS for port 587)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.cfg.SMTPHost,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		slog.Error("failed to start TLS",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// Authenticate
	if err = client.Auth(auth); err != nil {
		slog.Error("SMTP authentication failed",
			slog.String("error", err.Error()),
			slog.String("user", s.cfg.SMTPUser),
		)
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	// Set sender and recipient
	if err = client.Mail(s.cfg.SMTPFrom); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err = client.Rcpt(data.To); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send message
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = w.Write(message.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	client.Quit()

	slog.Info("email sent successfully",
		slog.String("to", data.To),
		slog.String("subject", data.Subject),
	)

	return nil
}

// encodeBase64 encodes bytes to base64 string with line breaks
func encodeBase64(data []byte) string {
	const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var buf bytes.Buffer

	encoder := make([]byte, 4)
	for i := 0; i < len(data); i += 3 {
		var n int
		if i+3 <= len(data) {
			n = 3
		} else {
			n = len(data) - i
		}

		var v uint32
		for j := 0; j < n; j++ {
			v |= uint32(data[i+j]) << (16 - j*8)
		}

		for j := 0; j < 4; j++ {
			if n == 1 && j > 1 {
				encoder[j] = '='
			} else if n == 2 && j > 2 {
				encoder[j] = '='
			} else {
				encoder[j] = base64Table[(v>>(18-j*6))&0x3F]
			}
		}

		buf.Write(encoder)

		// Add line break every 76 characters
		if (i/3+1)%19 == 0 {
			buf.WriteString("\r\n")
		}
	}

	return buf.String()
}

// SendReportToManager sends a PDF report to the user's manager
func (s *Service) SendReportToManager(managerName, managerEmail, userName, userEmail string, pdfContent []byte) error {
	subject := fmt.Sprintf("Career Development Report - %s", userName)

	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
				<h2 style="color: #1976d2; border-bottom: 2px solid #1976d2; padding-bottom: 10px;">
					Career Development Report
				</h2>
				
				<p>Hello <strong>%s</strong>,</p>
				
				<p>
					The employee <strong>%s</strong> (%s) has shared their career development report with you.
				</p>
				
				<p>
					This comprehensive report includes:
				</p>
				
				<ul>
					<li>Executive summary with performance statistics</li>
					<li>Completed and ongoing activities</li>
					<li>Distribution by levels and pillars</li>
					<li>Evidence and documentation</li>
					<li>Progress towards PDI objectives</li>
				</ul>
				
				<p>
					The complete report is available as a PDF attachment.
				</p>
				
				<hr style="border: none; border-top: 1px solid #ddd; margin: 20px 0;">
				
				<p style="font-size: 12px; color: #666;">
					<em>This is an automated email generated by Level Up Hub system.</em><br>
					<em>Generation date: %s</em>
				</p>
			</div>
		</body>
		</html>
	`, managerName, userName, userEmail, getCurrentDate())

	return s.SendEmail(EmailData{
		To:      managerEmail,
		Subject: subject,
		Body:    body,
		Attachment: &Attachment{
			Filename: fmt.Sprintf("report_%s.pdf", sanitizeFilename(userName)),
			Content:  pdfContent,
			MimeType: "application/pdf",
		},
	})
}

// sanitizeFilename removes special characters from filename
func sanitizeFilename(name string) string {
	// Simple sanitization - replace spaces with underscores and remove special chars
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

// getCurrentDate returns the current date formatted for display
func getCurrentDate() string {
	// Simplified date formatting
	return "Hoje"
}
