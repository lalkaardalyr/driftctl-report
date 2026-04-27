// Package diff provides utilities for comparing driftctl scan results
// across multiple runs to identify newly introduced or resolved drift.
package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/snyk/driftctl-report/internal/model"
)

// Baseline represents a saved snapshot of a scan result used as a
// reference point for future drift comparisons.
type Baseline struct {
	// CapturedAt is the UTC timestamp when the baseline was recorded.
	CapturedAt time.Time `json:"captured_at"`

	// ScanResult holds the model data captured at baseline time.
	ScanResult model.ScanResult `json:"scan_result"`
}

// SaveBaseline serialises the given ScanResult as a baseline JSON file
// at the specified path. Any existing file at that path is overwritten.
func SaveBaseline(path string, result model.ScanResult) error {
	b := Baseline{
		CapturedAt: time.Now().UTC(),
		ScanResult: result,
	}

	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write failed: %w", err)
	}

	return nil
}

// LoadBaseline reads and deserialises a baseline JSON file from path.
// It returns an error if the file cannot be read or parsed.
func LoadBaseline(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: read failed: %w", err)
	}

	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal failed: %w", err)
	}

	return &b, nil
}

// CompareWithBaseline loads the baseline at path and runs a full diff
// against current, returning a DiffResult that highlights changes since
// the baseline was captured.
//
// If path does not exist the function returns (nil, nil) so callers can
// treat a missing baseline as a first-run scenario without an error.
func CompareWithBaseline(path string, current model.ScanResult) (*DiffResult, error) {
	b, err := LoadBaseline(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	result := Compare(b.ScanResult, current)
	return &result, nil
}
