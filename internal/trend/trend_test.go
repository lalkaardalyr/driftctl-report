package trend_test

import (
	"testing"
	"time"

	"github.com/owner/driftctl-report/internal/model"
	"github.com/owner/driftctl-report/internal/trend"
)

func makeEntry(ts time.Time, drifted, managed, unmanaged int, coverage float64) trend.TimestampedResult {
	return trend.TimestampedResult{
		Timestamp: ts,
		Result: model.ScanResult{
			Summary: model.Summary{
				TotalManaged:   managed,
				TotalDrifted:   drifted,
				TotalUnmanaged: unmanaged,
				Coverage:       coverage,
			},
		},
	}
}

func TestFromResults_SortsChronologically(t *testing.T) {
	now := time.Now()
	entries := []trend.TimestampedResult{
		makeEntry(now.Add(2*time.Hour), 3, 10, 2, 0.8),
		makeEntry(now, 5, 10, 2, 0.6),
		makeEntry(now.Add(time.Hour), 4, 10, 2, 0.7),
	}
	s := trend.FromResults(entries)
	if len(s) != 3 {
		t.Fatalf("expected 3 points, got %d", len(s))
	}
	if !s[0].Timestamp.Equal(now) {
		t.Errorf("expected first point at %v, got %v", now, s[0].Timestamp)
	}
}

func TestSeries_Trend_Improving(t *testing.T) {
	now := time.Now()
	s := trend.FromResults([]trend.TimestampedResult{
		makeEntry(now, 10, 20, 5, 0.5),
		makeEntry(now.Add(time.Hour), 5, 20, 5, 0.75),
	})
	if got := s.Trend(); got != trend.Improving {
		t.Errorf("expected Improving, got %s", got)
	}
}

func TestSeries_Trend_Worsening(t *testing.T) {
	now := time.Now()
	s := trend.FromResults([]trend.TimestampedResult{
		makeEntry(now, 2, 20, 5, 0.9),
		makeEntry(now.Add(time.Hour), 8, 20, 5, 0.6),
	})
	if got := s.Trend(); got != trend.Worsening {
		t.Errorf("expected Worsening, got %s", got)
	}
}

func TestSeries_Trend_Stable(t *testing.T) {
	now := time.Now()
	s := trend.FromResults([]trend.TimestampedResult{
		makeEntry(now, 4, 20, 5, 0.8),
		makeEntry(now.Add(time.Hour), 4, 20, 5, 0.8),
	})
	if got := s.Trend(); got != trend.Stable {
		t.Errorf("expected Stable, got %s", got)
	}
}

func TestSeries_Latest_Empty(t *testing.T) {
	var s trend.Series
	_, ok := s.Latest()
	if ok {
		t.Error("expected false for empty series")
	}
}

func TestSeries_Latest_ReturnsMostRecent(t *testing.T) {
	now := time.Now()
	s := trend.FromResults([]trend.TimestampedResult{
		makeEntry(now, 5, 10, 2, 0.5),
		makeEntry(now.Add(time.Hour), 2, 10, 2, 0.8),
	})
	p, ok := s.Latest()
	if !ok {
		t.Fatal("expected ok=true")
	}
	if p.Drifted != 2 {
		t.Errorf("expected latest drifted=2, got %d", p.Drifted)
	}
}
