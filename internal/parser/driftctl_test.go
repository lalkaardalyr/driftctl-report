package parser_test

import (
	"strings"
	"testing"

	"github.com/example/driftctl-report/internal/parser"
)

const sampleJSON = `{
  "summary": {
    "total_resources": 5,
    "total_drifted": 1,
    "total_unmanaged": 2,
    "total_deleted": 0,
    "total_managed": 3
  },
  "managed": [
    {"id": "bucket-1", "type": "aws_s3_bucket"},
    {"id": "sg-123",   "type": "aws_security_group"}
  ],
  "unmanaged": [
    {"id": "bucket-2", "type": "aws_s3_bucket"}
  ],
  "deleted": [],
  "differences": [
    {
      "res": {"id": "bucket-1", "type": "aws_s3_bucket"},
      "changelog": [
        {"type": "update", "path": ["tags", "env"], "from": "prod", "to": "staging"}
      ]
    }
  ],
  "coverage": 60.0
}`

func TestParse_ValidJSON(t *testing.T) {
	report, err := parser.Parse(strings.NewReader(sampleJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if report.Summary.TotalResources != 5 {
		t.Errorf("expected TotalResources=5, got %d", report.Summary.TotalResources)
	}
	if report.Summary.TotalDrifted != 1 {
		t.Errorf("expected TotalDrifted=1, got %d", report.Summary.TotalDrifted)
	}
	if len(report.Managed) != 2 {
		t.Errorf("expected 2 managed resources, got %d", len(report.Managed))
	}
	if len(report.Unmanaged) != 1 {
		t.Errorf("expected 1 unmanaged resource, got %d", len(report.Unmanaged))
	}
	if len(report.Differences) != 1 {
		t.Fatalf("expected 1 difference, got %d", len(report.Differences))
	}
	if len(report.Differences[0].Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(report.Differences[0].Changes))
	}
	if report.Coverage != 60.0 {
		t.Errorf("expected coverage=60.0, got %f", report.Coverage)
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	_, err := parser.Parse(strings.NewReader(`{not valid json}`))
	if err == nil {
		t.Fatal("expected an error for invalid JSON, got nil")
	}
}

func TestParseFile_MissingFile(t *testing.T) {
	_, err := parser.ParseFile("/nonexistent/path/report.json")
	if err == nil {
		t.Fatal("expected an error for missing file, got nil")
	}
}
