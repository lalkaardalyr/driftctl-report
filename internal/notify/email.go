package notify

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
)

// EmailConfig holds SMTP connection and addressing configuration.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

// EmailNotifier sends drift alert notifications via SMTP.
type EmailNotifier struct {
	cfg  EmailConfig
	send func(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

// NewEmailNotifier creates an EmailNotifier with the given SMTP configuration.
func NewEmailNotifier(cfg EmailConfig) *EmailNotifier {
	return &EmailNotifier{
		cfg:  cfg,
		send: smtp.SendMail,
	}
}

// Send delivers a notification message via SMTP if the alert warrants it.
func (e *EmailNotifier) Send(msg Message) error {
	if msg.Severity == SeverityNone {
		return nil
	}

	body := e.buildEmail(msg)
	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)

	var auth smtp.Auth
	if e.cfg.Username != "" {
		auth = smtp.PlainAuth("", e.cfg.Username, e.cfg.Password, e.cfg.Host)
	}

	if err := e.send(addr, auth, e.cfg.From, e.cfg.To, []byte(body)); err != nil {
		return fmt.Errorf("email notifier: send failed: %w", err)
	}
	return nil
}

func (e *EmailNotifier) buildEmail(msg Message) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "From: %s\r\n", e.cfg.From)
	fmt.Fprintf(&buf, "To: %s\r\n", strings.Join(e.cfg.To, ", "))
	fmt.Fprintf(&buf, "Subject: [driftctl] %s\r\n", msg.Subject)
	fmt.Fprintf(&buf, "Content-Type: text/plain; charset=UTF-8\r\n")
	fmt.Fprintf(&buf, "\r\n")
	fmt.Fprintf(&buf, "%s\r\n", msg.Body)
	return buf.String()
}
