package export

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/snyk/driftctl-report/internal/model"
)

// XLSXWriter writes drift scan results in a tab-separated format suitable
// for import into spreadsheet applications (Excel, LibreOffice Calc, etc.).
// It reuses the CSV writer logic but with tab delimiters.
type XLSXWriter struct {
	w *csv.Writer
}

// NewXLSXWriter creates a new XLSXWriter that writes tab-separated values to w.
func NewXLSXWriter(w io.Writer) *XLSXWriter {
	cw := csv.NewWriter(w)
	cw.Comma = '\t'
	return &XLSXWriter{w: cw}
}

// Write outputs the scan result as a tab-separated spreadsheet with a header
// row followed by one row per drifted or unmanaged resource.
func (x *XLSXWriter) Write(result model.ScanResult) error {
	header := []string{"Resource ID", "Type", "Source", "Status"}
	if err := x.w.Write(header); err != nil {
		return fmt.Errorf("xlsx: write header: %w", err)
	}

	for _, r := range result.DriftedResources {
		row := []string{r.ResourceID, r.ResourceType, r.Source, "drifted"}
		if err := x.w.Write(row); err != nil {
			return fmt.Errorf("xlsx: write drifted row: %w", err)
		}
	}

	for _, r := range result.UnmanagedResources {
		row := []string{r.ResourceID, r.ResourceType, r.Source, "unmanaged"}
		if err := x.w.Write(row); err != nil {
			return fmt.Errorf("xlsx: write unmanaged row: %w", err)
		}
	}

	x.w.Flush()
	return x.w.Error()
}
