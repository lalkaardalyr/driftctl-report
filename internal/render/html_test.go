package render_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/driftctl-report/internal/model"
	"github.com/your-org/driftctl-report/internal/render"
)

func TestNew_ReturnsRenderer(t *testing.T) {
	r, err := render.New()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil renderer")
	}
}

func TestRender_ContainsSummaryInfo(t *testing.T) {
	r, err := render.New()
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	summary := model.Summary{
		TotalResources:     10,
		ManagedResources:   8,
		UnmanagedResources: 1,
		MissingResources:   1,
		ChangedResources:   0,
		Coverage:           80.0,
	}

	var buf bytes.Buffer
	if err := r.Render(&buf, summary, nil); err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	output := buf.String()

	for _, want := range []string{
		"Infrastructure Drift Report",
		"80.00%",
		"No drift detected",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestRender_DriftedResourcesTable(t *testing.T) {
	r, err := render.New()
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	summary := model.Summary{
		TotalResources:   3,
		ManagedResources: 1,
		MissingResources: 1,
		ChangedResources: 1,
		Coverage:         33.33,
	}
	resources := []model.DriftedResource{
		{Type: "aws_s3_bucket", ID: "my-bucket", Status: "missing"},
		{Type: "aws_iam_role", ID: "my-role", Status: "changed"},
	}

	var buf bytes.Buffer
	if err := r.Render(&buf, summary, resources); err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	output := buf.String()

	for _, want := range []string{
		"aws_s3_bucket",
		"my-bucket",
		"missing",
		"aws_iam_role",
		"my-role",
		"changed",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}

	if strings.Contains(output, "No drift detected") {
		t.Error("expected drift table, not 'no drift' message")
	}
}
