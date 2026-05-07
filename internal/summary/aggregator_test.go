package summary_test

import (
	"testing"
	"time"

	"github.com/owner/driftctl-report/internal/model"
	"github.com/owner/driftctl-report/internal/summary"
)

func makeEntry(t time.Time, drifted, unmanaged int, coverage float64) summary.Entry {
	return summary.Entry{
		ScannedAt: t,
		Result: model.ScanResult{
			DriftedResources:   make([]model.Resource, drifted),
			UnmanagedResources: make([]model.Resource, unmanaged),
			Summary:            model.Summary{Coverage: coverage},
		},
	}
}

func TestCompute_Empty(t *testing.T) {
	agg := summary.Compute(nil)
	if agg.TotalScans != 0 {
		t.Errorf("expected 0 scans, got %d", agg.TotalScans)
	}
}

func TestCompute_SingleEntry(t *testing.T) {
	now := time.Now()
	entry := makeEntry(now, 3, 5, 72.5)
	agg := summary.Compute([]summary.Entry{entry})

	if agg.TotalScans != 1 {
		t.Errorf("expected 1 scan, got %d", agg.TotalScans)
	}
	if agg.AvgDrifted != 3.0 {
		t.Errorf("expected AvgDrifted 3.0, got %f", agg.AvgDrifted)
	}
	if agg.AvgCoverage != 72.5 {
		t.Errorf("expected AvgCoverage 72.5, got %f", agg.AvgCoverage)
	}
}

func TestCompute_MultipleEntries_Averages(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := []summary.Entry{
		makeEntry(base, 2, 4, 80.0),
		makeEntry(base.Add(24*time.Hour), 4, 6, 60.0),
		makeEntry(base.Add(48*time.Hour), 6, 2, 70.0),
	}

	agg := summary.Compute(entries)

	if agg.TotalScans != 3 {
		t.Errorf("expected 3 scans, got %d", agg.TotalScans)
	}
	if agg.AvgDrifted != 4.0 {
		t.Errorf("expected AvgDrifted 4.0, got %f", agg.AvgDrifted)
	}
	if agg.MaxDrifted != 6 {
		t.Errorf("expected MaxDrifted 6, got %d", agg.MaxDrifted)
	}
	if agg.WorstCoverage != 60.0 {
		t.Errorf("expected WorstCoverage 60.0, got %f", agg.WorstCoverage)
	}
	if agg.BestCoverage != 80.0 {
		t.Errorf("expected BestCoverage 80.0, got %f", agg.BestCoverage)
	}
}

func TestCompute_SortsChronologically(t *testing.T) {
	base := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	entries := []summary.Entry{
		makeEntry(base.Add(48*time.Hour), 1, 1, 90.0),
		makeEntry(base, 2, 2, 50.0),
		makeEntry(base.Add(24*time.Hour), 3, 3, 70.0),
	}

	agg := summary.Compute(entries)

	if !agg.FirstScan.Equal(base) {
		t.Errorf("expected FirstScan %v, got %v", base, agg.FirstScan)
	}
	if !agg.LastScan.Equal(base.Add(48 * time.Hour)) {
		t.Errorf("expected LastScan %v, got %v", base.Add(48*time.Hour), agg.LastScan)
	}
}
