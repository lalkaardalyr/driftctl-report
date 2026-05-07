package model

import "fmt"

// CoverageLabel returns a human-readable coverage tier.
// Tiers: Excellent (>=90%), Good (>=70%), Fair (>=50%), Poor (<50%).
func (s Summary) CoverageLabel() string {
	switch {
	case s.CoveragePercent >= 90:
		return "Excellent"
	case s.CoveragePercent >= 70:
		return "Good"
	case s.CoveragePercent >= 50:
		return "Fair"
	default:
		return "Poor"
	}
}

// CoverageFormatted returns the coverage percentage as a formatted string.
func (s Summary) CoverageFormatted() string {
	return fmt.Sprintf("%.1f%%", s.CoveragePercent)
}

// HasDrift returns true when any drift was detected in the scan.
func (r ScanResult) HasDrift() bool {
	return len(r.DriftedResources) > 0 ||
		len(r.UnmanagedResources) > 0 ||
		len(r.DeletedResources) > 0
}

// DriftCount returns the total number of drifted, unmanaged, and deleted resources.
func (r ScanResult) DriftCount() int {
	return len(r.DriftedResources) + len(r.UnmanagedResources) + len(r.DeletedResources)
}
