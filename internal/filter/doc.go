// Package filter provides utilities for filtering, sorting, and grouping
// collections of drift resources.
//
// It is intended to be used by rendering and reporting components to
// select and organise resources before presentation.
//
// Example usage:
//
//	drifted := filter.Apply(resources, filter.ByType("aws_s3_bucket"))
//	byType := filter.GroupByType(drifted)
//	sorted := filter.Sort(drifted, filter.SortByID)
package filter
