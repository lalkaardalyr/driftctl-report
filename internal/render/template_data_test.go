package render

import (
	"testing"

	"github.com/owner/driftctl-report/internal/model"
)

func makeTestScanResult() model.ScanResult {
	return model.ScanResult{
		Summary: model.Summary{
			TotalResources: 10,
			DriftedCount:   2,
			MissingCount:   1,
			UnmanagedCount: 3,
			Coverage:       40.0,
		},
		DriftedResources: []model.Resource{
			{Type: "aws_s3_bucket", ID: "my-bucket"},
			{Type: "aws_iam_role", ID: "my-role"},
		},
		MissingResources: []model.Resource{
			{Type: "aws_lambda_function", ID: "my-fn"},
		},
		UnmanagedResources: []model.Resource{},
	}
}

func TestNewTemplateData_Summary(t *testing.T) {
	result := makeTestScanResult()
	data := NewTemplateData(result)

	if data.Summary.TotalResources != 10 {
		t.Errorf("expected TotalResources=10, got %d", data.Summary.TotalResources)
	}
	if data.Summary.DriftedCount != 2 {
		t.Errorf("expected DriftedCount=2, got %d", data.Summary.DriftedCount)
	}
	if data.Summary.MissingCount != 1 {
		t.Errorf("expected MissingCount=1, got %d", data.Summary.MissingCount)
	}
	if data.Summary.UnmanagedCount != 3 {
		t.Errorf("expected UnmanagedCount=3, got %d", data.Summary.UnmanagedCount)
	}
}

func TestNewTemplateData_CoverageClass_Bad(t *testing.T) {
	result := makeTestScanResult() // coverage 40% => Bad
	data := NewTemplateData(result)
	if data.Summary.CoverageClass != "coverage-bad" {
		t.Errorf("expected coverage-bad, got %s", data.Summary.CoverageClass)
	}
}

func TestNewTemplateData_CoverageClass_Good(t *testing.T) {
	result := makeTestScanResult()
	result.Summary.Coverage = 95.0
	data := NewTemplateData(result)
	if data.Summary.CoverageClass != "coverage-ok" {
		t.Errorf("expected coverage-ok, got %s", data.Summary.CoverageClass)
	}
}

func TestNewTemplateData_ResourceViews(t *testing.T) {
	result := makeTestScanResult()
	data := NewTemplateData(result)

	if len(data.DriftedResources) != 2 {
		t.Fatalf("expected 2 drifted resources, got %d", len(data.DriftedResources))
	}
	if data.DriftedResources[0].Type != "aws_s3_bucket" {
		t.Errorf("unexpected type: %s", data.DriftedResources[0].Type)
	}
	if data.DriftedResources[0].ID != "my-bucket" {
		t.Errorf("unexpected id: %s", data.DriftedResources[0].ID)
	}
	if len(data.MissingResources) != 1 {
		t.Fatalf("expected 1 missing resource, got %d", len(data.MissingResources))
	}
	if len(data.UnmanagedResources) != 0 {
		t.Errorf("expected 0 unmanaged resources, got %d", len(data.UnmanagedResources))
	}
}

func TestToResourceViews_Empty(t *testing.T) {
	views := toResourceViews([]model.Resource{})
	if len(views) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(views))
	}
}
