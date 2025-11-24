package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

type EmailConfig struct {
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	From     string
}

type EmailMessage struct {
	To      string
	Subject string
	Body    string
	IsHTML  bool
}

// LoadEmailConfig - Config'i bir kez yükleyin ve cache'leyin
func LoadEmailConfig() *EmailConfig {
	return &EmailConfig{
		SMTPHost: GetEnv("SMTP_HOST", ""),
		SMTPPort: GetEnv("SMTP_PORT", "587"),
		SMTPUser: GetEnv("SMTP_USER", ""),
		SMTPPass: GetEnv("SMTP_PASSWORD", ""),
		From:     GetEnv("EMAIL_FROM", ""),
	}
}

func (c *EmailConfig) Validate() error {
	if c.SMTPHost == "" || c.SMTPPort == "" || c.SMTPUser == "" || c.SMTPPass == "" || c.From == "" {
		return fmt.Errorf("SMTP configuration is incomplete")
	}
	return nil
}

func SendEmailWithContext(ctx context.Context, to, subject, body string, isHTML bool) error {
	config := LoadEmailConfig()
	
	if err := config.Validate(); err != nil {
		return err
	}

	msg := &EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  isHTML,
	}

	return sendEmailWithRetry(ctx, config, msg, 3)
}

// Backward compatibility
func SendEmail(to string, body string) error {
	return SendEmailWithContext(context.Background(), to, "No Subject", body, false)
}

func sendEmailWithRetry(ctx context.Context, config *EmailConfig, msg *EmailMessage, maxRetries int) error {
	var lastErr error
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
			LogInfo(fmt.Sprintf("Retrying email send (attempt %d/%d)", attempt+1, maxRetries))
		}

		err := sendEmail(ctx, config, msg)
		if err == nil {
			return nil
		}
		
		lastErr = err
		
		// Kalıcı hatalar için retry yapma
		if isPermanentError(err) {
			break
		}
	}
	
	return fmt.Errorf("failed to send email after %d attempts: %w", maxRetries, lastErr)
}

func sendEmail(ctx context.Context, config *EmailConfig, msg *EmailMessage) error {
	// Email validation
	if _, err := mail.ParseAddress(msg.To); err != nil {
		return fmt.Errorf("invalid recipient email: %w", err)
	}
	if _, err := mail.ParseAddress(config.From); err != nil {
		return fmt.Errorf("invalid sender email: %w", err)
	}

	// Context-aware TCP dial with timeout
	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}
	
	conn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(config.SMTPHost, config.SMTPPort))
	if err != nil {
		return fmt.Errorf("tcp dial failed: %w", err)
	}

	c, err := smtp.NewClient(conn, config.SMTPHost)
	if err != nil {
		return fmt.Errorf("smtp client creation failed: %w", err)
	}
	defer c.Close()

	// Check for STARTTLS support
	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName:         config.SMTPHost,
			InsecureSkipVerify: false, // NEVER set to true in production
			MinVersion:         tls.VersionTLS12,
		}

		if err = c.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("starttls failed: %w", err)
		}
	}

	// Authentication - GÜVENLİK: Sadece hata durumunu loglayın, credentials'ı asla!
	auth := smtp.PlainAuth("", config.SMTPUser, config.SMTPPass, config.SMTPHost)
	if err = c.Auth(auth); err != nil {
		LogError("SMTP authentication failed") // Detay vermeyin!
		return fmt.Errorf("smtp auth failed: %w", err)
	}

	// Send mail
	if err = c.Mail(config.From); err != nil {
		return fmt.Errorf("mail from failed: %w", err)
	}

	if err = c.Rcpt(msg.To); err != nil {
		return fmt.Errorf("rcpt to failed: %w", err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("data command failed: %w", err)
	}
	defer w.Close()

	// Construct email with proper headers
	emailBody := buildEmailBody(config.From, msg)
	
	if _, err = w.Write([]byte(emailBody)); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("close data writer failed: %w", err)
	}

	if err = c.Quit(); err != nil {
		return fmt.Errorf("quit failed: %w", err)
	}

	LogInfo("Email sent successfully to: " + maskEmail(msg.To))
	return nil
}

func buildEmailBody(from string, msg *EmailMessage) string {
	contentType := "text/plain; charset=UTF-8"
	if msg.IsHTML {
		contentType = "text/html; charset=UTF-8"
	}

	return fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: %s\r\n\r\n%s",
		from,
		msg.To,
		msg.Subject,
		contentType,
		msg.Body,
	)
}

// GÜVENLİK: Email adresini maskeleyerek loglayın
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***"
	}
	
	username := parts[0]
	if len(username) <= 2 {
		return "**@" + parts[1]
	}
	
	return string(username[0]) + "***" + string(username[len(username)-1]) + "@" + parts[1]
}

func isPermanentError(err error) bool {
	errStr := err.Error()
	permanentErrors := []string{
		"invalid recipient",
		"invalid sender",
		"smtp auth failed",
	}
	
	for _, permErr := range permanentErrors {
		if strings.Contains(errStr, permErr) {
			return true
		}
	}
	return false
}