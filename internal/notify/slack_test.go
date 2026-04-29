package notify_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/owner/driftctl-report/internal/notify"
)

func TestSlackNotifier_Send_Success(t *testing.T) {
	var received string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewSlackNotifier(server.URL)
	msg := notify.Message{
		Level:   notify.LevelWarning,
		Subject: "drift detected",
		Body:    "3 resources drifted",
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received != "application/json" {
		t.Errorf("expected application/json, got %s", received)
	}
}

func TestSlackNotifier_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := notify.NewSlackNotifier(server.URL)
	msg := notify.Message{Level: notify.LevelInfo, Subject: "test", Body: "body"}
	err := n.Send(msg)
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestSlackNotifier_Send_InvalidURL(t *testing.T) {
	n := notify.NewSlackNotifier("http://127.0.0.1:0/invalid")
	msg := notify.Message{Level: notify.LevelInfo, Subject: "test", Body: "body"}
	err := n.Send(msg)
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}

func TestSlackNotifier_Send_PayloadContainsSubject(t *testing.T) {
	var body []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body = make([]byte, r.ContentLength)
		r.Body.Read(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.NewSlackNotifier(server.URL)
	msg := notify.Message{Level: notify.LevelCritical, Subject: "critical-drift", Body: "many resources"}
	_ = n.Send(msg)

	if !contains(string(body), "critical-drift") {
		t.Errorf("expected payload to contain subject, got: %s", string(body))
	}
	_ = fmt.Sprintf("") // suppress unused import
}
