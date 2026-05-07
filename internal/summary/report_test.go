package summary_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/owner/driftctl-report/internal/summary"
)

func buildAgg() summary.Aggregate {
	base := time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC)
	entries := []summary.Entry{
		makeEntry(base, 4, 8, 65.0),
		makeEntry(base.Add(24*time.Hour), 2, 3, 80.0),
	}
	return summary.Compute(entries)
}

func TestWriteJSON_ValidOutput(t *testing.T) {
	agg := buildAgg()
	var buf bytes.Buffer
	if err := summary.WriteJSON(&buf, agg); err != nil {
		t.Fatalf("WriteJSON returned error: %v", err)
	}

	var decoded summary.Aggregate
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if decoded.TotalScans != 2 {
		t.Errorf("expected TotalScans 2, got %d", decoded.TotalScans)
	}
}

func TestWriteText_ContainsLabels(t *testing.T) {
	agg := buildAgg()
	var buf bytes.Buffer
	if err := summary.WriteText(&buf, agg); err != nil {
		t.Fatalf("WriteText returned error: %v", err)
	}

	out := buf.String()
	for _, expected := range []string{"Total scans", "Avg drifted", "Best coverage", "Worst coverage"} {
		if !strings.Contains(out, expected) {
			t.Errorf("output missing %q", expected)
		}
	}
}

func TestWriteText_EmptyAggregate(t *testing.T) {
	var buf bytes.Buffer
	if err := summary.WriteText(&buf, summary.Aggregate{}); err != nil {
		t.Fatalf("WriteText returned error: %v", err)
	}
	if !strings.Contains(buf.String(), "No scan data") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteText_CoverageValues(t *testing.T) {
	agg := buildAgg()
	var buf bytes.Buffer
	_ = summary.WriteText(&buf, agg)
	out := buf.String()

	if !strings.Contains(out, "72.5%") {
		t.Errorf("expected average coverage 72.5%% in output, got:\n%s", out)
	}
}
