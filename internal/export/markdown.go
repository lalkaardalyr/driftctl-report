package export

import (
	"fmt"
	"io"
	"text/template"

	"github.com/your-org/driftctl-report/internal/model"
)

const markdownTmpl = `# Drift Report

## Summary

| Metric | Value |
|--------|-------|
| Coverage | {{ .Summary.CoverageFormatted }} |
| Total Resources | {{ .Summary.Total }} |
| Managed | {{ .Summary.Managed }} |
| Unmanaged | {{ .Summary.Unmanaged }} |
| Drifted | {{ .Summary.Drifted }} |
| Missing | {{ .Summary.Missing }} |

## Drifted Resources
{{ if .DriftedResources }}
| ID | Type | Source |
|----|------|--------|
{{ range .DriftedResources }}| {{ .ID }} | {{ .Type }} | {{ .Source }} |
{{ end }}{{ else }}
_No drifted resources detected._
{{ end }}
## Unmanaged Resources
{{ if .UnmanagedResources }}
| ID | Type | Source |
|----|------|--------|
{{ range .UnmanagedResources }}| {{ .ID }} | {{ .Type }} | {{ .Source }} |
{{ end }}{{ else }}
_No unmanaged resources detected._
{{ end }}
`

type markdownRow struct {
	ID     string
	Type   string
	Source string
}

type markdownData struct {
	Summary            model.Summary
	DriftedResources   []markdownRow
	UnmanagedResources []markdownRow
}

// MarkdownWriter writes a Markdown-formatted drift report.
type MarkdownWriter struct {
	tmpl *template.Template
}

// NewMarkdownWriter creates a new MarkdownWriter.
func NewMarkdownWriter() (*MarkdownWriter, error) {
	t, err := template.New("markdown").Parse(markdownTmpl)
	if err != nil {
		return nil, fmt.Errorf("markdown: parse template: %w", err)
	}
	return &MarkdownWriter{tmpl: t}, nil
}

// Write renders the scan result as Markdown to the given writer.
func (m *MarkdownWriter) Write(w io.Writer, result model.ScanResult) error {
	data := markdownData{
		Summary:            result.Summary,
		DriftedResources:   toMarkdownRows(result.DriftedResources),
		UnmanagedResources: toMarkdownRows(result.UnmanagedResources),
	}
	if err := m.tmpl.Execute(w, data); err != nil {
		return fmt.Errorf("markdown: render: %w", err)
	}
	return nil
}

func toMarkdownRows(resources []model.Resource) []markdownRow {
	rows := make([]markdownRow, 0, len(resources))
	for _, r := range resources {
		rows = append(rows, markdownRow{
			ID:     r.ID,
			Type:   r.Type,
			Source: r.Source,
		})
	}
	return rows
}
