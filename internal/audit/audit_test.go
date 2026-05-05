package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/example/driftctl-report/internal/audit"
)

func TestLog_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf, "ci-bot")

	err := l.Log(audit.Event{
		Kind:      audit.EventScanStarted,
		InputPath: "scan.json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var got audit.Event
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got.Kind != audit.EventScanStarted {
		t.Errorf("expected kind %q, got %q", audit.EventScanStarted, got.Kind)
	}
	if got.Actor != "ci-bot" {
		t.Errorf("expected actor %q, got %q", "ci-bot", got.Actor)
	}
	if got.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestScanStarted_SetsFields(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf, "")
	if err := l.ScanStarted("drift.json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got audit.Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.InputPath != "drift.json" {
		t.Errorf("expected input_path %q, got %q", "drift.json", got.InputPath)
	}
}

func TestScanCompleted_CountsPresent(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf, "svc")
	if err := l.ScanCompleted("in.json", "out.html", 3, 7); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got audit.Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.Drifted != 3 {
		t.Errorf("expected drifted=3, got %d", got.Drifted)
	}
	if got.Unmanaged != 7 {
		t.Errorf("expected unmanaged=7, got %d", got.Unmanaged)
	}
	if got.OutputPath != "out.html" {
		t.Errorf("expected output_path %q, got %q", "out.html", got.OutputPath)
	}
}

func TestScanFailed_ErrorField(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf, "")
	if err := l.ScanFailed("bad.json", "file not found"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got audit.Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.Error != "file not found" {
		t.Errorf("expected error field %q, got %q", "file not found", got.Error)
	}
	if got.Kind != audit.EventScanFailed {
		t.Errorf("expected kind %q, got %q", audit.EventScanFailed, got.Kind)
	}
}
