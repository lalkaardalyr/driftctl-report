package metrics_test

import (
	"testing"

	"github.com/snyk/driftctl-report/internal/metrics"
	"github.com/snyk/driftctl-report/internal/model"
)

func makeScanResult() model.ScanResult {
	return model.ScanResult{
		Summary: model.Summary{
			TotalManaged: 10,
			Coverage:     85.0,
		},
		DriftedResources: []model.Resource{
			{ID: "res-1", Type: "aws_s3_bucket"},
			{ID: "res-2", Type: "aws_s3_bucket"},
			{ID: "res-3", Type: "aws_iam_role"},
		},
		UnmanagedResources: []model.Resource{
			{ID: "res-4", Type: "aws_lambda_function"},
		},
	}
}

func TestFromScanResult_Totals(t *testing.T) {
	sr := makeScanResult()
	r := metrics.FromScanResult(sr)

	if r.TotalManaged != 10 {
		t.Errorf("expected TotalManaged=10, got %d", r.TotalManaged)
	}
	if r.TotalDrifted != 3 {
		t.Errorf("expected TotalDrifted=3, got %d", r.TotalDrifted)
	}
	if r.TotalUnmanaged != 1 {
		t.Errorf("expected TotalUnmanaged=1, got %d", r.TotalUnmanaged)
	}
	if r.Coverage != 85.0 {
		t.Errorf("expected Coverage=85.0, got %f", r.Coverage)
	}
}

func TestFromScanResult_ByType(t *testing.T) {
	sr := makeScanResult()
	r := metrics.FromScanResult(sr)

	if len(r.ByType) != 3 {
		t.Fatalf("expected 3 resource types, got %d", len(r.ByType))
	}
	// sorted alphabetically: aws_iam_role, aws_lambda_function, aws_s3_bucket
	if r.ByType[0].Type != "aws_iam_role" {
		t.Errorf("expected first type=aws_iam_role, got %s", r.ByType[0].Type)
	}
	if r.ByType[2].DriftedCount != 2 {
		t.Errorf("expected aws_s3_bucket DriftedCount=2, got %d", r.ByType[2].DriftedCount)
	}
	if r.ByType[1].UnmanagedCount != 1 {
		t.Errorf("expected aws_lambda_function UnmanagedCount=1, got %d", r.ByType[1].UnmanagedCount)
	}
}

func TestCoverageLabel(t *testing.T) {
	cases := []struct {
		coverage float64
		want     string
	}{
		{95.0, "good"},
		{75.0, "warning"},
		{50.0, "critical"},
	}
	for _, tc := range cases {
		r := metrics.Report{Coverage: tc.coverage}
		if got := r.CoverageLabel(); got != tc.want {
			t.Errorf("coverage=%.1f: expected label=%q, got %q", tc.coverage, tc.want, got)
		}
	}
}

func TestCoverageFormatted(t *testing.T) {
	r := metrics.Report{Coverage: 87.5}
	if got := r.CoverageFormatted(); got != "87.50%" {
		t.Errorf("expected \"87.50%%\", got %q", got)
	}
}
