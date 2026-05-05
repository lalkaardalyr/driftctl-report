package notify

import (
	"errors"
	"net/smtp"
	"strings"
	"testing"
)

func TestEmailNotifier_Send_SkipsOnNoDrift(t *testing.T) {
	var called bool
	n := NewEmailNotifier(EmailConfig{Host: "localhost", Port: 25, From: "a@b.com", To: []string{"x@y.com"}})
	n.send = func(_ string, _ smtp.Auth, _ string, _ []string, _ []byte) error {
		called = true
		return nil
	}

	msg := Message{Severity: SeverityNone, Subject: "no drift", Body: "all good"}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected SMTP send to be skipped for SeverityNone")
	}
}

func TestEmailNotifier_Send_Success(t *testing.T) {
	var capturedMsg []byte
	n := NewEmailNotifier(EmailConfig{Host: "smtp.example.com", Port: 587, From: "bot@example.com", To: []string{"ops@example.com"}})
	n.send = func(_ string, _ smtp.Auth, _ string, _ []string, msg []byte) error {
		capturedMsg = msg
		return nil
	}

	msg := Message{Severity: SeverityWarning, Subject: "3 drifted resources", Body: "Check your infra."}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(capturedMsg), "3 drifted resources") {
		t.Error("expected subject in email body")
	}
	if !strings.Contains(string(capturedMsg), "Check your infra.") {
		t.Error("expected body text in email")
	}
}

func TestEmailNotifier_Send_PropagatesError(t *testing.T) {
	n := NewEmailNotifier(EmailConfig{Host: "smtp.example.com", Port: 587, From: "bot@example.com", To: []string{"ops@example.com"}})
	n.send = func(_ string, _ smtp.Auth, _ string, _ []string, _ []byte) error {
		return errors.New("connection refused")
	}

	msg := Message{Severity: SeverityWarning, Subject: "drift", Body: "some drift"}
	err := n.Send(msg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestEmailNotifier_Send_UsesCorrectAddress(t *testing.T) {
	var capturedAddr string
	n := NewEmailNotifier(EmailConfig{Host: "mail.corp.io", Port: 465, From: "ci@corp.io", To: []string{"team@corp.io"}})
	n.send = func(addr string, _ smtp.Auth, _ string, _ []string, _ []byte) error {
		capturedAddr = addr
		return nil
	}

	_ = n.Send(Message{Severity: SeverityCritical, Subject: "critical drift", Body: "urgent"})
	if capturedAddr != "mail.corp.io:465" {
		t.Errorf("expected addr mail.corp.io:465, got %s", capturedAddr)
	}
}
