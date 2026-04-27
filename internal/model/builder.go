package model

import "github.com/snyk/driftctl/pkg/analyser"

// FromAnalysis converts a driftctl analyser.Analysis into a ScanResult.
func FromAnalysis(a *analyser.Analysis) ScanResult {
	result := ScanResult{}

	for _, r := range a.Managed() {
		result.ManagedResources = append(result.ManagedResources, Resource{
			ID:   r.ResourceId(),
			Type: r.ResourceType(),
		})
	}

	for _, r := range a.Unmanaged() {
		result.UnmanagedResources = append(result.UnmanagedResources, Resource{
			ID:   r.ResourceId(),
			Type: r.ResourceType(),
		})
	}

	for _, r := range a.Deleted() {
		result.DeletedResources = append(result.DeletedResources, Resource{
			ID:   r.ResourceId(),
			Type: r.ResourceType(),
		})
	}

	for _, dr := range a.Differences() {
		drifted := DriftedResource{
			Resource: Resource{
				ID:   dr.Res.ResourceId(),
				Type: dr.Res.ResourceType(),
			},
		}
		for _, ch := range dr.Changelog {
			drifted.Differences = append(drifted.Differences, Difference{
				FieldPath: ch.Path,
				Previous:  ch.From,
				Current:   ch.To,
			})
		}
		result.DriftedResources = append(result.DriftedResources, drifted)
	}

	total := len(result.ManagedResources) + len(result.UnmanagedResources) + len(result.DeletedResources)
	var coverage float64
	if total > 0 {
		coverage = float64(len(result.ManagedResources)) / float64(total) * 100
	}

	result.Summary = Summary{
		TotalResources:  total,
		Managed:         len(result.ManagedResources),
		Unmanaged:       len(result.UnmanagedResources),
		Deleted:         len(result.DeletedResources),
		Drifted:         len(result.DriftedResources),
		CoveragePercent: coverage,
	}

	return result
}
