package export_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/snyk/driftctl-report/internal/export"
	"github.com/snyk/driftctl-report/internal/model"
)

func makeXLSXScanResult() model.ScanResult {
	return model.ScanResult{
		DriftedResources: []model.Resource{
			{ResourceID: "sg-abc123", ResourceType: "aws_security_group", Source: "aws"},
		},
		UnmanagedResources: []model.Resource{
			{ResourceID: "i-xyz789", ResourceType: "aws_instance", Source: "aws"},
		},
		Summary: model.Summary{
			TotalResources:   10,
			DriftedCount:     1,
			UnmanagedCount:   1,
			CoveragePercent:  80.0,
		},
	}
}

func TestXLSXWriter_CreatesWriter(t *testing.T) {
	var buf bytes.Buffer
	w := export.NewXLSXWriter(&buf)
	if w == nil {
		t.Fatal("expected non-nil XLSXWriter")
	}
}

func TestXLSXWriter_WritesHeader(t *testing.T) {
	var buf bytes.Buffer
	w := export.NewXLSXWriter(&buf)

	if err := w.Write(makeXLSXScanResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) == 0 {
		t.Fatal("expected at least one line")
	}

	header := lines[0]
	for _, col := range []string{"Resource ID", "Type", "Source", "Status"} {
		if !strings.Contains(header, col) {
			t.Errorf("header missing column %q, got: %s", col, header)
		}
	}
}

func TestXLSXWriter_WritesDriftedAndUnmanaged(t *testing.T) {
	var buf bytes.Buffer
	w := export.NewXLSXWriter(&buf)

	if err := w.Write(makeXLSXScanResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "sg-abc123") {
		t.Error("expected drifted resource ID in output")
	}
	if !strings.Contains(out, "drifted") {
		t.Error("expected status 'drifted' in output")
	}
	if !strings.Contains(out, "i-xyz789") {
		t.Error("expected unmanaged resource ID in output")
	}
	if !strings.Contains(out, "unmanaged") {
		t.Error("expected status 'unmanaged' in output")
	}
}

func TestXLSXWriter_EmptyResult(t *testing.T) {
	var buf bytes.Buffer
	w := export.NewXLSXWriter(&buf)

	empty := model.ScanResult{}
	if err := w.Write(empty); err != nil {
		t.Fatalf("unexpected error on empty result: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header line, got %d lines", len(lines))
	}
}
