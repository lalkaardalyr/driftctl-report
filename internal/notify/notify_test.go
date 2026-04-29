package notify_test

import (
	"bytes"
	"testing"

	"github.com/owner/driftctl-report/internal/model"
	"github.com/owner/driftctl-report/internal/notify"
)

func makeResult(drifted, managed, unmanaged, missing int) model.ScanResult {
	resources := make([]model.Resource, drifted)
	for i := range resources {
		resources[i] = model.Resource{ID: fmt.Sprintf("res-%d", i), Type: "aws_instance"}
	}
	return model.ScanResult{
		DriftedResources:   resources,
		UnmanagedResources: make([]model.Resource, unmanaged),
		Summary: model.Summary{
			Managed:   managed,
			Drifted:   drifted,
			Unmanaged: unmanaged,
			Missing:   missing,
			Total:     managed + drifted + unmanaged + missing,
		},
	}
}

func TestEvaluate_NoDrift_ReturnsFalse(t *testing.T) {
	result := makeResult(0, 5, 0, 0)
	_, ok := notify.Evaluate(result, notify.DefaultConfig())
	if ok {
		t.Fatal("expected no notification when no drift")
	}
}

func TestEvaluate_Warning(t *testing.T) {
	result := makeResult(3, 10, 0, 0)
	msg, ok := notify.Evaluate(result, notify.DefaultConfig())
	if !ok {
		t.Fatal("expected notification")
	}
	if msg.Level != notify.LevelWarning {
		t.Errorf("expected warning, got %s", msg.Level)
	}
}

func TestEvaluate_Critical(t *testing.T) {
	result := makeResult(12, 20, 0, 0)
	msg, ok := notify.Evaluate(result, notify.Config{WarnThreshold: 1, CriticalThreshold: 10})
	if !ok {
		t.Fatal("expected notification")
	}
	if msg.Level != notify.LevelCritical {
		t.Errorf("expected critical, got %s", msg.Level)
	}
}

func TestEvaluate_SubjectContainsDriftCount(t *testing.T) {
	result := makeResult(5, 10, 2, 1)
	msg, _ := notify.Evaluate(result, notify.DefaultConfig())
	if !contains(msg.Subject, "5") {
		t.Errorf("expected subject to contain drift count, got: %s", msg.Subject)
	}
}

func TestWriteMessage(t *testing.T) {
	msg := notify.Message{
		Level:   notify.LevelWarning,
		Subject: "test subject",
		Body:    "test body",
	}
	var buf bytes.Buffer
	if err := notify.WriteMessage(&buf, msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(buf.String(), "warning") {
		t.Errorf("expected output to contain level, got: %s", buf.String())
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
