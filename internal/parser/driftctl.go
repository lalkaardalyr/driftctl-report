package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Resource represents a single cloud resource in the driftctl report.
type Resource struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// DriftctlReport is the top-level structure of a driftctl JSON output file.
type DriftctlReport struct {
	Summary struct {
		TotalResources  int `json:"total_resources"`
		TotalDrifted    int `json:"total_drifted"`
		TotalUnmanaged  int `json:"total_unmanaged"`
		TotalDeleted    int `json:"total_deleted"`
		TotalManaged    int `json:"total_managed"`
	} `json:"summary"`
	Managed    []Resource          `json:"managed"`
	Unmanaged  []Resource          `json:"unmanaged"`
	Deleted    []Resource          `json:"deleted"`
	Differences []Difference       `json:"differences"`
	Coverage   float64             `json:"coverage"`
}

// Difference represents a resource that exists in both state and cloud but has drifted.
type Difference struct {
	Res     Resource  `json:"res"`
	Changes []Change  `json:"changelog"`
}

// Change describes a single field-level change within a drifted resource.
type Change struct {
	Type string      `json:"type"`
	Path []string    `json:"path"`
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
}

// ParseFile reads and parses a driftctl JSON report from the given file path.
func ParseFile(path string) (*DriftctlReport, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parser: opening file %q: %w", path, err)
	}
	defer f.Close()
	return Parse(f)
}

// Parse reads and parses a driftctl JSON report from an io.Reader.
func Parse(r io.Reader) (*DriftctlReport, error) {
	var report DriftctlReport
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&report); err != nil {
		return nil, fmt.Errorf("parser: decoding JSON: %w", err)
	}
	return &report, nil
}
