package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/driftctl-report/internal/notify"
)

func TestWebhookNotifier_Send_Success(t *testing.T) {
	var received map[string]string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewWebhookNotifier(server.URL)
	msg := notify.Message{Subject: "Drift detected", Body: "3 resources drifted", Warning: true}

	if err := n.Send(msg); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if received["subject"] != "Drift detected" {
		t.Errorf("subject mismatch: got %q", received["subject"])
	}
	if received["severity"] != "warning" {
		t.Errorf("severity mismatch: got %q", received["severity"])
	}
	if received["timestamp"] == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestWebhookNotifier_Send_CriticalSeverity(t *testing.T) {
	var received map[string]string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received) //nolint:errcheck
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewWebhookNotifier(server.URL)
	msg := notify.Message{Subject: "Critical drift", Body: "10 resources", Critical: true}

	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["severity"] != "critical" {
		t.Errorf("expected critical, got %q", received["severity"])
	}
}

func TestWebhookNotifier_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := notify.NewWebhookNotifier(server.URL)
	msg := notify.Message{Subject: "test", Body: "body"}

	if err := n.Send(msg); err == nil {
		t.Error("expected error for non-2xx status")
	}
}

func TestWebhookNotifier_Send_InvalidURL(t *testing.T) {
	n := notify.NewWebhookNotifier("http://127.0.0.1:0/no-server")
	msg := notify.Message{Subject: "test", Body: "body"}

	if err := n.Send(msg); err == nil {
		t.Error("expected error for unreachable URL")
	}
}
