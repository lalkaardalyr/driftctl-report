package export

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/driftctl-report/internal/model"
)

func makeMarkdownScanResult() model.ScanResult {
	return model.ScanResult{
		Summary: model.Summary{
			Total:     3,
			Managed:   1,
			Unmanaged: 1,
			Drifted:   1,
			Missing:   0,
			Coverage:  33.3,
		},
		DriftedResources: []model.Resource{
			{ID: "bucket-1", Type: "aws_s3_bucket", Source: "terraform.tfstate"},
		},
		UnmanagedResources: []model.Resource{
			{ID: "sg-abc", Type: "aws_security_group", Source: ""},
		},
	}
}

func TestMarkdownWriter_CreatesWriter(t *testing.T) {
	w, err := NewMarkdownWriter()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil MarkdownWriter")
	}
}

func TestMarkdownWriter_ContainsSummary(t *testing.T) {
	w, _ := NewMarkdownWriter()
	var buf bytes.Buffer
	result := makeMarkdownScanResult()

	if err := w.Write(&buf, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "# Drift Report") {
		t.Error("expected markdown heading")
	}
	if !strings.Contains(out, "33.3") {
		t.Error("expected coverage value in output")
	}
	if !strings.Contains(out, "3") {
		t.Error("expected total resource count")
	}
}

func TestMarkdownWriter_DriftedResourcesTable(t *testing.T) {
	w, _ := NewMarkdownWriter()
	var buf bytes.Buffer
	result := makeMarkdownScanResult()

	_ = w.Write(&buf, result)
	out := buf.String()

	if !strings.Contains(out, "bucket-1") {
		t.Error("expected drifted resource ID in output")
	}
	if !strings.Contains(out, "aws_s3_bucket") {
		t.Error("expected drifted resource type in output")
	}
}

func TestMarkdownWriter_EmptyResult(t *testing.T) {
	w, _ := NewMarkdownWriter()
	var buf bytes.Buffer
	empty := model.ScanResult{}

	if err := w.Write(&buf, empty); err != nil {
		t.Fatalf("unexpected error on empty result: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "_No drifted resources detected._") {
		t.Error("expected no-drift message")
	}
	if !strings.Contains(out, "_No unmanaged resources detected._") {
		t.Error("expected no-unmanaged message")
	}
}
