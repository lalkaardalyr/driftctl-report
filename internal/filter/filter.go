// Package filter provides utilities for filtering and grouping drift resources.
package filter

import "github.com/snyk/driftctl-report/internal/model"

// ResourceFilter defines a predicate for selecting resources.
type ResourceFilter func(r model.Resource) bool

// ByType returns a ResourceFilter that matches resources of the given type.
func ByType(resourceType string) ResourceFilter {
	return func(r model.Resource) bool {
		return r.Type == resourceType
	}
}

// BySource returns a ResourceFilter that matches resources from the given source.
func BySource(source string) ResourceFilter {
	return func(r model.Resource) bool {
		return r.Source == source
	}
}

// Apply returns only the resources that satisfy all provided filters.
func Apply(resources []model.Resource, filters ...ResourceFilter) []model.Resource {
	result := make([]model.Resource, 0, len(resources))
	for _, r := range resources {
		if matchAll(r, filters...) {
			result = append(result, r)
		}
	}
	return result
}

// GroupByType groups resources by their Type field.
func GroupByType(resources []model.Resource) map[string][]model.Resource {
	groups := make(map[string][]model.Resource)
	for _, r := range resources {
		groups[r.Type] = append(groups[r.Type], r)
	}
	return groups
}

// GroupBySource groups resources by their Source field.
func GroupBySource(resources []model.Resource) map[string][]model.Resource {
	groups := make(map[string][]model.Resource)
	for _, r := range resources {
		groups[r.Source] = append(groups[r.Source], r)
	}
	return groups
}

func matchAll(r model.Resource, filters ...ResourceFilter) bool {
	for _, f := range filters {
		if !f(r) {
			return false
		}
	}
	return true
}
