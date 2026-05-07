// Package metrics provides aggregated drift metrics computed from scan results.
// It is intended for use in dashboards, trend reports, and policy evaluations.
package metrics

import (
	"fmt"
	"sort"

	"github.com/snyk/driftctl-report/internal/model"
)

// ResourceMetrics holds per-resource-type drift counts.
type ResourceMetrics struct {
	Type          string
	DriftedCount  int
	UnmanagedCount int
	Total         int
}

// Report is the top-level metrics report derived from a ScanResult.
type Report struct {
	TotalManaged   int
	TotalDrifted   int
	TotalUnmanaged int
	Coverage       float64
	ByType         []ResourceMetrics
}

// CoverageLabel returns a human-readable label for the coverage percentage.
func (r Report) CoverageLabel() string {
	switch {
	case r.Coverage >= 90:
		return "good"
	case r.Coverage >= 70:
		return "warning"
	default:
		return "critical"
	}
}

// CoverageFormatted returns coverage as a percentage string, e.g. "87.50%".
func (r Report) CoverageFormatted() string {
	return fmt.Sprintf("%.2f%%", r.Coverage)
}

// FromScanResult builds a metrics Report from a model.ScanResult.
func FromScanResult(sr model.ScanResult) Report {
	counts := map[string]*ResourceMetrics{}

	for _, res := range sr.DriftedResources {
		m := getOrCreate(counts, res.Type)
		m.DriftedCount++
		m.Total++
	}

	for _, res := range sr.UnmanagedResources {
		m := getOrCreate(counts, res.Type)
		m.UnmanagedCount++
		m.Total++
	}

	byType := make([]ResourceMetrics, 0, len(counts))
	for _, v := range counts {
		byType = append(byType, *v)
	}
	sort.Slice(byType, func(i, j int) bool {
		return byType[i].Type < byType[j].Type
	})

	return Report{
		TotalManaged:   sr.Summary.TotalManaged,
		TotalDrifted:   len(sr.DriftedResources),
		TotalUnmanaged: len(sr.UnmanagedResources),
		Coverage:       sr.Summary.Coverage,
		ByType:         byType,
	}
}

func getOrCreate(m map[string]*ResourceMetrics, t string) *ResourceMetrics {
	if _, ok := m[t]; !ok {
		m[t] = &ResourceMetrics{Type: t}
	}
	return m[t]
}
