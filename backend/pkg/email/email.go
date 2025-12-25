package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"strings"

	"github.com/bricksocoolxd/bengi-investment-system/pkg/config"
)

// Config holds email configuration
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

var emailConfig *Config

// Initialize sets up the email service
func Initialize() {
	emailConfig = &Config{
		Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		Port:     getEnvInt("SMTP_PORT", 587),
		Username: getEnv("SMTP_USERNAME", ""),
		Password: getEnv("SMTP_PASSWORD", ""),
		From:     getEnv("SMTP_FROM", "noreply@bengi.io"),
	}

	if emailConfig.Username == "" {
		log.Println("[Email] SMTP not configured, email sending disabled")
	} else {
		log.Println("✅ Email service initialized")
	}
}

// IsConfigured returns true if email is configured
func IsConfigured() bool {
	return emailConfig != nil && emailConfig.Username != ""
}

// SendEmail sends a plain text email
func SendEmail(to, subject, body string) error {
	return SendEmailMultiple([]string{to}, subject, body)
}

// SendEmailMultiple sends email to multiple recipients
func SendEmailMultiple(to []string, subject, body string) error {
	if !IsConfigured() {
		log.Printf("[Email] Email not configured, skipping: %s", subject)
		return nil
	}

	auth := smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.Host)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n",
		emailConfig.From,
		strings.Join(to, ","),
		subject,
		body,
	)

	addr := fmt.Sprintf("%s:%d", emailConfig.Host, emailConfig.Port)
	err := smtp.SendMail(addr, auth, emailConfig.From, to, []byte(msg))
	if err != nil {
		log.Printf("[Email] Failed to send email: %v", err)
		return err
	}

	log.Printf("[Email] Sent email to %v: %s", to, subject)
	return nil
}

// SendHTMLEmail sends an HTML email
func SendHTMLEmail(to, subject, htmlBody string) error {
	if !IsConfigured() {
		return nil
	}

	auth := smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.Host)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n",
		emailConfig.From,
		to,
		subject,
		htmlBody,
	)

	addr := fmt.Sprintf("%s:%d", emailConfig.Host, emailConfig.Port)
	return smtp.SendMail(addr, auth, emailConfig.From, []string{to}, []byte(msg))
}

// SendTemplate sends an email using a template
func SendTemplate(to, subject, templateName string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("templates/email/%s.html", templateName))
	if err != nil {
		// Try inline template
		return SendEmail(to, subject, fmt.Sprintf("%v", data))
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	return SendHTMLEmail(to, subject, buf.String())
}

// Notification types
func SendWelcomeEmail(to, name string) error {
	subject := "Welcome to Bengi Investment System"
	body := fmt.Sprintf("Hello %s,\n\nWelcome to Bengi Investment System!\n\nStart trading now.\n\nBest regards,\nBengi Team", name)
	return SendEmail(to, subject, body)
}

func SendOrderConfirmation(to, orderID, symbol, side string, quantity, price float64) error {
	subject := fmt.Sprintf("Order Confirmation - %s %s", side, symbol)
	body := fmt.Sprintf("Your %s order for %s has been placed.\n\nOrder ID: %s\nQuantity: %.4f\nPrice: $%.2f\n\nThank you for trading with us!",
		side, symbol, orderID, quantity, price)
	return SendEmail(to, subject, body)
}

func SendTradeExecution(to, tradeID, symbol, side string, quantity, price, total float64) error {
	subject := fmt.Sprintf("Trade Executed - %s %s", side, symbol)
	body := fmt.Sprintf("Your %s order for %s has been executed.\n\nTrade ID: %s\nQuantity: %.4f\nPrice: $%.2f\nTotal: $%.2f",
		side, symbol, tradeID, quantity, price, total)
	return SendEmail(to, subject, body)
}

func SendPasswordChanged(to string) error {
	subject := "Password Changed"
	body := "Your password has been successfully changed.\n\nIf you did not make this change, please contact support immediately."
	return SendEmail(to, subject, body)
}

// Helper functions
func getEnv(key, defaultVal string) string {
	if val := config.AppConfig; val != nil {
		// Use config values if available
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	return defaultVal
}
