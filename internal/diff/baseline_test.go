package diff_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/driftctl-report/internal/diff"
	"github.com/example/driftctl-report/internal/model"
)

func makeBaselineResult() model.ScanResult {
	return model.ScanResult{
		Summary: model.Summary{
			TotalResources:    10,
			ManagedResources:  8,
			DriftedResources:  2,
			UnmanagedResources: 1,
			Coverage:          80.0,
		},
		Drifted: []model.Resource{
			{ID: "aws_s3_bucket.logs", Type: "aws_s3_bucket", Source: "aws"},
			{ID: "aws_iam_role.deployer", Type: "aws_iam_role", Source: "aws"},
		},
		Unmanaged: []model.Resource{
			{ID: "aws_ec2_instance.orphan", Type: "aws_instance", Source: "aws"},
		},
	}
}

func TestSaveBaseline_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	result := makeBaselineResult()
	if err := diff.SaveBaseline(result, path); err != nil {
		t.Fatalf("SaveBaseline returned unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected baseline file to exist, but it does not")
	}
}

func TestSaveBaseline_WritesValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	result := makeBaselineResult()
	if err := diff.SaveBaseline(result, path); err != nil {
		t.Fatalf("SaveBaseline returned unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read baseline file: %v", err)
	}

	var parsed model.ScanResult
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("baseline file is not valid JSON: %v", err)
	}

	if len(parsed.Drifted) != 2 {
		t.Errorf("expected 2 drifted resources, got %d", len(parsed.Drifted))
	}
}

func TestLoadBaseline_ReturnsResult(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	original := makeBaselineResult()
	if err := diff.SaveBaseline(original, path); err != nil {
		t.Fatalf("SaveBaseline returned unexpected error: %v", err)
	}

	loaded, err := diff.LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline returned unexpected error: %v", err)
	}

	if len(loaded.Drifted) != len(original.Drifted) {
		t.Errorf("expected %d drifted resources, got %d", len(original.Drifted), len(loaded.Drifted))
	}

	if len(loaded.Unmanaged) != len(original.Unmanaged) {
		t.Errorf("expected %d unmanaged resources, got %d", len(original.Unmanaged), len(loaded.Unmanaged))
	}
}

func TestLoadBaseline_MissingFile(t *testing.T) {
	_, err := diff.LoadBaseline("/nonexistent/path/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestCompareWithBaseline_DetectsNewAndResolved(t *testing.T) {
	baseline := makeBaselineResult()

	// Current scan: one original drift resolved, one new drift added
	current := model.ScanResult{
		Drifted: []model.Resource{
			{ID: "aws_s3_bucket.logs", Type: "aws_s3_bucket", Source: "aws"},
			{ID: "aws_lambda_function.processor", Type: "aws_lambda_function", Source: "aws"},
		},
		Unmanaged: []model.Resource{},
	}

	report := diff.CompareWithBaseline(baseline, current)

	if len(report.NewlyDrifted) != 1 {
		t.Errorf("expected 1 newly drifted resource, got %d", len(report.NewlyDrifted))
	}
	if report.NewlyDrifted[0].ID != "aws_lambda_function.processor" {
		t.Errorf("unexpected newly drifted resource: %s", report.NewlyDrifted[0].ID)
	}

	if len(report.Resolved) != 1 {
		t.Errorf("expected 1 resolved resource, got %d", len(report.Resolved))
	}
	if report.Resolved[0].ID != "aws_iam_role.deployer" {
		t.Errorf("unexpected resolved resource: %s", report.Resolved[0].ID)
	}
}
