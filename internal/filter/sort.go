package filter

import (
	"sort"

	"github.com/snyk/driftctl-report/internal/model"
)

// SortBy defines the field to sort resources by.
type SortBy int

const (
	// SortByID sorts resources alphabetically by ID.
	SortByID SortBy = iota
	// SortByType sorts resources alphabetically by Type, with ID as a tiebreaker.
	SortByType
)

// Sort returns a new slice of resources sorted by the given field.
// The original slice is not modified.
func Sort(resources []model.Resource, by SortBy) []model.Resource {
	copy_ := make([]model.Resource, len(resources))
	copy(copy_, resources)

	switch by {
	case SortByType:
		sort.Slice(copy_, func(i, j int) bool {
			if copy_[i].Type == copy_[j].Type {
				return copy_[i].ID < copy_[j].ID
			}
			return copy_[i].Type < copy_[j].Type
		})
	default: // SortByID
		sort.Slice(copy_, func(i, j int) bool {
			return copy_[i].ID < copy_[j].ID
		})
	}

	return copy_
}

// UniqueTypes returns a sorted, deduplicated list of resource types.
func UniqueTypes(resources []model.Resource) []string {
	seen := make(map[string]struct{})
	for _, r := range resources {
		seen[r.Type] = struct{}{}
	}
	types := make([]string, 0, len(seen))
	for t := range seen {
		types = append(types, t)
	}
	sort.Strings(types)
	return types
}

// FilterByType returns a new slice containing only resources of the given type.
// The original slice is not modified.
func FilterByType(resources []model.Resource, resourceType string) []model.Resource {
	result := make([]model.Resource, 0)
	for _, r := range resources {
		if r.Type == resourceType {
			result = append(result, r)
		}
	}
	return result
}
