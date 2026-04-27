package export_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/owner/driftctl-report/internal/export"
	"github.com/owner/driftctl-report/internal/model"
)

func makeJSONScanResult() model.ScanResult {
	return model.ScanResult{
		Summary: model.Summary{
			Total:     5,
			Managed:   3,
			Unmanaged: 1,
			Drifted:   1,
			Coverage:  60.0,
		},
		DriftedResources: []model.Resource{
			{ResourceID: "bucket-1", ResourceType: "aws_s3_bucket", Source: "terraform"},
		},
		UnmanagedResources: []model.Resource{
			{ResourceID: "sg-99", ResourceType: "aws_security_group"},
		},
	}
}

func TestJSONWriter_WritesValidJSON(t *testing.T) {
	var buf bytes.Buffer
	w := export.NewJSONWriter(&buf)

	if err := w.Write(makeJSONScanResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}

func TestJSONWriter_SummaryFields(t *testing.T) {
	var buf bytes.Buffer
	w := export.NewJSONWriter(&buf)
	_ = w.Write(makeJSONScanResult())

	var out struct {
		Summary struct {
			Total    int     `json:"total_resources"`
			Coverage float64 `json:"coverage_percent"`
		} `json:"summary"`
	}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Summary.Total != 5 {
		t.Errorf("expected total=5, got %d", out.Summary.Total)
	}
	if out.Summary.Coverage != 60.0 {
		t.Errorf("expected coverage=60.0, got %f", out.Summary.Coverage)
	}
}

func TestJSONWriter_DriftedAndUnmanaged(t *testing.T) {
	var buf bytes.Buffer
	w := export.NewJSONWriter(&buf)
	_ = w.Write(makeJSONScanResult())

	var out struct {
		Drifted   []map[string]interface{} `json:"drifted_resources"`
		Unmanaged []map[string]interface{} `json:"unmanaged_resources"`
	}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.Drifted) != 1 {
		t.Errorf("expected 1 drifted resource, got %d", len(out.Drifted))
	}
	if len(out.Unmanaged) != 1 {
		t.Errorf("expected 1 unmanaged resource, got %d", len(out.Unmanaged))
	}
	if out.Drifted[0]["id"] != "bucket-1" {
		t.Errorf("unexpected drifted id: %v", out.Drifted[0]["id"])
	}
}

func TestJSONWriter_EmptyResult(t *testing.T) {
	var buf bytes.Buffer
	w := export.NewJSONWriter(&buf)
	if err := w.Write(model.ScanResult{}); err != nil {
		t.Fatalf("unexpected error on empty result: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output for empty result")
	}
}
