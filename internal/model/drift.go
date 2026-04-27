package model

// Summary holds aggregated counts from a driftctl scan result.
type Summary struct {
	TotalResources  int
	Managed         int
	Unmanaged       int
	Deleted         int
	Drifted         int
	CoveragePercent float64
}

// Resource represents a single cloud resource from the scan.
type Resource struct {
	ID   string
	Type string
}

// DriftedResource represents a resource that has drifted from its IaC definition.
type DriftedResource struct {
	Resource
	Differences []Difference
}

// Difference describes a single field-level change on a drifted resource.
type Difference struct {
	FieldPath string
	Previous  interface{}
	Current   interface{}
}

// ScanResult is the normalised representation of a driftctl JSON output.
type ScanResult struct {
	Summary          Summary
	ManagedResources []Resource
	UnmanagedResources []Resource
	DeletedResources []Resource
	DriftedResources []DriftedResource
}
