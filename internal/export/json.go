package export

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/owner/driftctl-report/internal/model"
)

// JSONWriter writes a ScanResult as indented JSON.
type JSONWriter struct {
	w io.Writer
}

// NewJSONWriter creates a new JSONWriter that writes to w.
func NewJSONWriter(w io.Writer) *JSONWriter {
	return &JSONWriter{w: w}
}

// jsonExport is the structure serialised to JSON.
type jsonExport struct {
	Summary   jsonSummary    `json:"summary"`
	Drifted   []jsonResource `json:"drifted_resources"`
	Unmanaged []jsonResource `json:"unmanaged_resources"`
}

type jsonSummary struct {
	Total      int     `json:"total_resources"`
	Managed    int     `json:"managed"`
	Unmanaged  int     `json:"unmanaged"`
	Drifted    int     `json:"drifted"`
	Coverage   float64 `json:"coverage_percent"`
}

type jsonResource struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Source string `json:"source,omitempty"`
}

// Write serialises result to JSON and writes it to the underlying writer.
func (j *JSONWriter) Write(result model.ScanResult) error {
	export := jsonExport{
		Summary: jsonSummary{
			Total:    result.Summary.Total,
			Managed:  result.Summary.Managed,
			Unmanaged: result.Summary.Unmanaged,
			Drifted:  result.Summary.Drifted,
			Coverage: result.Summary.Coverage,
		},
	}

	for _, r := range result.DriftedResources {
		export.Drifted = append(export.Drifted, jsonResource{
			ID:     r.ResourceID,
			Type:   r.ResourceType,
			Source: r.Source,
		})
	}

	for _, r := range result.UnmanagedResources {
		export.Unmanaged = append(export.Unmanaged, jsonResource{
			ID:   r.ResourceID,
			Type: r.ResourceType,
			Source: r.Source,
		})
	}

	enc := json.NewEncoder(j.w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(export); err != nil {
		return fmt.Errorf("json export: encode: %w", err)
	}
	return nil
}
