// Package export provides functionality to export drift scan results
// into various machine-readable formats such as CSV.
package export

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/snyk/driftctl-report/internal/model"
)

// CSVWriter writes a ScanResult to a CSV format.
type CSVWriter struct {
	w *csv.Writer
}

// NewCSVWriter returns a new CSVWriter that writes to w.
func NewCSVWriter(w io.Writer) *CSVWriter {
	return &CSVWriter{w: csv.NewWriter(w)}
}

// Write serialises the drifted and unmanaged resources from result into CSV
// rows and flushes them to the underlying writer.
func (c *CSVWriter) Write(result model.ScanResult) error {
	header := []string{"resource_id", "resource_type", "status", "source"}
	if err := c.w.Write(header); err != nil {
		return fmt.Errorf("export csv: write header: %w", err)
	}

	for _, r := range result.Summary.DriftedResources {
		row := []string{r.ResourceID, r.ResourceType, "drifted", r.Source}
		if err := c.w.Write(row); err != nil {
			return fmt.Errorf("export csv: write row: %w", err)
		}
	}

	for _, r := range result.Summary.UnmanagedResources {
		row := []string{r.ResourceID, r.ResourceType, "unmanaged", r.Source}
		if err := c.w.Write(row); err != nil {
			return fmt.Errorf("export csv: write row: %w", err)
		}
	}

	c.w.Flush()
	return c.w.Error()
}
