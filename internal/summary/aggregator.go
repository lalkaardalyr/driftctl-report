// Package summary provides multi-scan aggregation across historical scan results.
package summary

import (
	"sort"
	"time"

	"github.com/owner/driftctl-report/internal/model"
)

// Entry holds a single scan result paired with its scan time.
type Entry struct {
	ScannedAt time.Time
	Result    model.ScanResult
}

// Aggregate holds rolled-up statistics derived from multiple scan entries.
type Aggregate struct {
	TotalScans      int
	FirstScan       time.Time
	LastScan        time.Time
	AvgDrifted      float64
	AvgUnmanaged    float64
	AvgCoverage     float64
	MaxDrifted      int
	MaxUnmanaged    int
	WorstCoverage   float64
	BestCoverage    float64
}

// Compute derives an Aggregate from a slice of Entries.
// Entries need not be pre-sorted; Compute sorts by ScannedAt internally.
func Compute(entries []Entry) Aggregate {
	if len(entries) == 0 {
		return Aggregate{}
	}

	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ScannedAt.Before(sorted[j].ScannedAt)
	})

	agg := Aggregate{
		TotalScans:    len(sorted),
		FirstScan:     sorted[0].ScannedAt,
		LastScan:      sorted[len(sorted)-1].ScannedAt,
		WorstCoverage: 100.0,
		BestCoverage:  0.0,
	}

	var sumDrifted, sumUnmanaged, sumCoverage float64

	for _, e := range sorted {
		d := len(e.Result.DriftedResources)
		u := len(e.Result.UnmanagedResources)
		cov := e.Result.Summary.Coverage

		sumDrifted += float64(d)
		sumUnmanaged += float64(u)
		sumCoverage += cov

		if d > agg.MaxDrifted {
			agg.MaxDrifted = d
		}
		if u > agg.MaxUnmanaged {
			agg.MaxUnmanaged = u
		}
		if cov < agg.WorstCoverage {
			agg.WorstCoverage = cov
		}
		if cov > agg.BestCoverage {
			agg.BestCoverage = cov
		}
	}

	n := float64(len(sorted))
	agg.AvgDrifted = sumDrifted / n
	agg.AvgUnmanaged = sumUnmanaged / n
	agg.AvgCoverage = sumCoverage / n

	return agg
}
