// Package diff provides utilities for comparing two driftctl scan results
// and producing a structured change summary between runs.
package diff

import (
	"github.com/owner/driftctl-report/internal/model"
)

// ResourceKey uniquely identifies a resource by type and ID.
type ResourceKey struct {
	Type string
	ID   string
}

// Result holds the outcome of comparing two ScanResults.
type Result struct {
	// NewlyDrifted are resources that appear drifted in current but not in baseline.
	NewlyDrifted []model.Resource
	// Resolved are resources that were drifted in baseline but are no longer drifted.
	Resolved []model.Resource
	// NewlyUnmanaged are resources unmanaged in current but not in baseline.
	NewlyUnmanaged []model.Resource
	// StillDrifted are resources drifted in both baseline and current.
	StillDrifted []model.Resource
}

// Compare produces a diff.Result by comparing a baseline ScanResult to a current ScanResult.
func Compare(baseline, current model.ScanResult) Result {
	baselineDrifted := indexResources(baseline.DriftedResources)
	currentDrifted := indexResources(current.DriftedResources)
	baselineUnmanaged := indexResources(baseline.UnmanagedResources)
	currentUnmanaged := indexResources(current.UnmanagedResources)

	var result Result

	for key, res := range currentDrifted {
		if _, existed := baselineDrifted[key]; existed {
			result.StillDrifted = append(result.StillDrifted, res)
		} else {
			result.NewlyDrifted = append(result.NewlyDrifted, res)
		}
	}

	for key, res := range baselineDrifted {
		if _, stillDrifted := currentDrifted[key]; !stillDrifted {
			result.Resolved = append(result.Resolved, res)
		}
	}

	for key, res := range currentUnmanaged {
		if _, existed := baselineUnmanaged[key]; !existed {
			result.NewlyUnmanaged = append(result.NewlyUnmanaged, res)
		}
	}

	return result
}

func indexResources(resources []model.Resource) map[ResourceKey]model.Resource {
	idx := make(map[ResourceKey]model.Resource, len(resources))
	for _, r := range resources {
		key := ResourceKey{Type: r.Type, ID: r.ID}
		idx[key] = r
	}
	return idx
}
