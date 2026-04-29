// Package trend provides utilities for tracking drift metrics over time
// by comparing multiple scan results chronologically.
package trend

import (
	"sort"
	"time"

	"github.com/owner/driftctl-report/internal/model"
)

// Point represents a single scan result at a specific point in time.
type Point struct {
	Timestamp    time.Time
	TotalManaged int
	Drifted      int
	Unmanaged    int
	Coverage     float64
}

// Series is an ordered collection of trend points.
type Series []Point

// Direction describes whether drift is improving, worsening, or stable.
type Direction string

const (
	Improving Direction = "improving"
	Worsening Direction = "worsening"
	Stable    Direction = "stable"
)

// FromResults builds a Series from a slice of ScanResults paired with timestamps.
func FromResults(entries []TimestampedResult) Series {
	s := make(Series, 0, len(entries))
	for _, e := range entries {
		s = append(s, Point{
			Timestamp:    e.Timestamp,
			TotalManaged: e.Result.Summary.TotalManaged,
			Drifted:      e.Result.Summary.TotalDrifted,
			Unmanaged:    e.Result.Summary.TotalUnmanaged,
			Coverage:     e.Result.Summary.Coverage,
		})
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].Timestamp.Before(s[j].Timestamp)
	})
	return s
}

// TimestampedResult pairs a ScanResult with a timestamp.
type TimestampedResult struct {
	Timestamp time.Time
	Result    model.ScanResult
}

// Trend returns the overall drift direction across the series.
func (s Series) Trend() Direction {
	if len(s) < 2 {
		return Stable
	}
	first := s[0]
	last := s[len(s)-1]
	switch {
	case last.Drifted < first.Drifted:
		return Improving
	case last.Drifted > first.Drifted:
		return Worsening
	default:
		return Stable
	}
}

// Latest returns the most recent Point in the series, and false if empty.
func (s Series) Latest() (Point, bool) {
	if len(s) == 0 {
		return Point{}, false
	}
	return s[len(s)-1], true
}
