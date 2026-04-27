package export_test

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"
	"testing"

	"github.com/snyk/driftctl-report/internal/export"
	"github.com/snyk/driftctl-report/internal/model"
)

func makeExportScanResult() model.ScanResult {
	return model.ScanResult{
		Summary: model.Summary{
			DriftedResources: []model.Resource{
				{ResourceID: "bucket-1", ResourceType: "aws_s3_bucket", Source: "aws"},
			},
			UnmanagedResources: []model.Resource{
				{ResourceID: "sg-abc", ResourceType: "aws_security_group", Source: "aws"},
			},
		},
	}
}

func TestCSVWriter_WritesHeader(t *testing.T) {
	var buf bytes.Buffer
	w := export.NewCSVWriter(&buf)

	if err := w.Write(makeExportScanResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := csv.NewReader(strings.NewReader(buf.String()))
	header, err := r.Read()
	if err != nil {
		t.Fatalf("failed to read header: %v", err)
	}

	expected := []string{"resource_id", "resource_type", "status", "source"}
	for i, col := range expected {
		if header[i] != col {
			t.Errorf("header[%d]: got %q, want %q", i, header[i], col)
		}
	}
}

func TestCSVWriter_WritesDriftedAndUnmanaged(t *testing.T) {
	var buf bytes.Buffer
	w := export.NewCSVWriter(&buf)

	if err := w.Write(makeExportScanResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := csv.NewReader(strings.NewReader(buf.String()))
	// skip header
	if _, err := r.Read(); err != nil {
		t.Fatalf("failed to skip header: %v", err)
	}

	var rows [][]string
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("unexpected read error: %v", err)
		}
		rows = append(rows, row)
	}

	if len(rows) != 2 {
		t.Fatalf("expected 2 data rows, got %d", len(rows))
	}
	if rows[0][2] != "drifted" {
		t.Errorf("row 0 status: got %q, want \"drifted\"", rows[0][2])
	}
	if rows[1][2] != "unmanaged" {
		t.Errorf("row 1 status: got %q, want \"unmanaged\"", rows[1][2])
	}
}

func TestCSVWriter_EmptyResult(t *testing.T) {
	var buf bytes.Buffer
	w := export.NewCSVWriter(&buf)

	if err := w.Write(model.ScanResult{}); err != nil {
		t.Fatalf("unexpected error on empty result: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header row, got %d lines", len(lines))
	}
}
