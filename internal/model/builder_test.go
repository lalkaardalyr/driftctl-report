package model_test

import (
	"testing"

	"github.com/owner/driftctl-report/internal/model"
)

func makeScanResult(managed, unmanaged, deleted, drifted int) model.ScanResult {
	res := model.ScanResult{}
	for i := 0; i < managed; i++ {
		res.ManagedResources = append(res.ManagedResources, model.Resource{ID: "m", Type: "aws_s3_bucket"})
	}
	for i := 0; i < unmanaged; i++ {
		res.UnmanagedResources = append(res.UnmanagedResources, model.Resource{ID: "u", Type: "aws_instance"})
	}
	for i := 0; i < deleted; i++ {
		res.DeletedResources = append(res.DeletedResources, model.Resource{ID: "d", Type: "aws_vpc"})
	}
	for i := 0; i < drifted; i++ {
		res.DriftedResources = append(res.DriftedResources, model.DriftedResource{
			Resource: model.Resource{ID: "dr", Type: "aws_security_group"},
		})
	}
	total := managed + unmanaged + deleted
	var cov float64
	if total > 0 {
		cov = float64(managed) / float64(total) * 100
	}
	res.Summary = model.Summary{
		TotalResources:  total,
		Managed:         managed,
		Unmanaged:       unmanaged,
		Deleted:         deleted,
		Drifted:         drifted,
		CoveragePercent: cov,
	}
	return res
}

func TestSummary_CoverageLabel(t *testing.T) {
	cases := []struct {
		pct   float64
		want  string
	}{
		{100, "Excellent"},
		{90, "Excellent"},
		{75, "Good"},
		{55, "Fair"},
		{30, "Poor"},
	}
	for _, tc := range cases {
		s := model.Summary{CoveragePercent: tc.pct}
		if got := s.CoverageLabel(); got != tc.want {
			t.Errorf("pct=%.0f: got %q, want %q", tc.pct, got, tc.want)
		}
	}
}

func TestScanResult_HasDrift(t *testing.T) {
	clean := makeScanResult(5, 0, 0, 0)
	if clean.HasDrift() {
		t.Error("expected no drift for fully managed result")
	}
	drifty := makeScanResult(5, 2, 1, 1)
	if !drifty.HasDrift() {
		t.Error("expected drift to be detected")
	}
}

func TestSummary_CoverageFormatted(t *testing.T) {
	s := model.Summary{CoveragePercent: 83.333}
	if got := s.CoverageFormatted(); got != "83.3%" {
		t.Errorf("got %q, want \"83.3%%\"", got)
	}
}
