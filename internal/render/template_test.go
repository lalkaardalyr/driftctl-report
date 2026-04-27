package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/owner/driftctl-report/internal/model"
)

func TestParsedTemplate_NotNil(t *testing.T) {
	if ParsedTemplate == nil {
		t.Fatal("ParsedTemplate should not be nil")
	}
}

func TestParsedTemplate_RendersSummary(t *testing.T) {
	result := model.ScanResult{
		Summary: model.Summary{
			TotalResources: 5,
			DriftedCount:   1,
			MissingCount:   0,
			UnmanagedCount: 2,
			Coverage:       80.0,
		},
		DriftedResources:   []model.Resource{{Type: "aws_s3_bucket", ID: "test-bucket"}},
		MissingResources:   []model.Resource{},
		UnmanagedResources: []model.Resource{},
	}

	data := NewTemplateData(result)
	var buf bytes.Buffer
	if err := ParsedTemplate.Execute(&buf, data); err != nil {
		t.Fatalf("template execution failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Driftctl Infrastructure Report") {
		t.Error("expected report title in output")
	}
	if !strings.Contains(output, "80.00%") {
		t.Error("expected coverage percentage in output")
	}
	if !strings.Contains(output, "aws_s3_bucket") {
		t.Error("expected drifted resource type in output")
	}
	if !strings.Contains(output, "test-bucket") {
		t.Error("expected drifted resource ID in output")
	}
}

func TestParsedTemplate_NoDriftMessage(t *testing.T) {
	result := model.ScanResult{
		Summary: model.Summary{
			TotalResources: 3,
			Coverage:       100.0,
		},
		DriftedResources:   []model.Resource{},
		MissingResources:   []model.Resource{},
		UnmanagedResources: []model.Resource{},
	}

	data := NewTemplateData(result)
	var buf bytes.Buffer
	if err := ParsedTemplate.Execute(&buf, data); err != nil {
		t.Fatalf("template execution failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No drifted resources found") {
		t.Error("expected no-drift message for drifted section")
	}
	if !strings.Contains(output, "No missing resources found") {
		t.Error("expected no-drift message for missing section")
	}
	if !strings.Contains(output, "No unmanaged resources found") {
		t.Error("expected no-drift message for unmanaged section")
	}
}
