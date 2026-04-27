package render

import (
	"html/template"
	"io"
	"time"

	"github.com/your-org/driftctl-report/internal/model"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>driftctl Report</title>
  <style>
    body { font-family: sans-serif; margin: 2rem; background: #f9f9f9; color: #333; }
    h1 { color: #2c3e50; }
    .summary { background: #fff; border-radius: 8px; padding: 1rem 2rem; box-shadow: 0 1px 4px rgba(0,0,0,0.1); margin-bottom: 2rem; }
    .label { display: inline-block; padding: 0.2rem 0.6rem; border-radius: 4px; font-weight: bold; }
    .label-ok { background: #d4edda; color: #155724; }
    .label-warn { background: #fff3cd; color: #856404; }
    .label-crit { background: #f8d7da; color: #721c24; }
    table { width: 100%; border-collapse: collapse; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 4px rgba(0,0,0,0.1); }
    th { background: #2c3e50; color: #fff; padding: 0.6rem 1rem; text-align: left; }
    td { padding: 0.5rem 1rem; border-bottom: 1px solid #eee; }
    tr:last-child td { border-bottom: none; }
    .badge { display: inline-block; padding: 0.15rem 0.5rem; border-radius: 3px; font-size: 0.85em; }
    .badge-missing { background: #f8d7da; color: #721c24; }
    .badge-unmanaged { background: #fff3cd; color: #856404; }
    .badge-changed { background: #cce5ff; color: #004085; }
  </style>
</head>
<body>
  <h1>Infrastructure Drift Report</h1>
  <p>Generated: {{.GeneratedAt}}</p>
  <div class="summary">
    <h2>Summary</h2>
    <p>Coverage: <strong>{{.Summary.CoverageFormatted}}</strong>
      <span class="label label-{{.Summary.CoverageLabel | lower}}">{{.Summary.CoverageLabel}}</span>
    </p>
    <p>Total Resources: {{.Summary.TotalResources}} &nbsp;|&nbsp;
       Managed: {{.Summary.ManagedResources}} &nbsp;|&nbsp;
       Unmanaged: {{.Summary.UnmanagedResources}} &nbsp;|&nbsp;
       Missing: {{.Summary.MissingResources}} &nbsp;|&nbsp;
       Changed: {{.Summary.ChangedResources}}</p>
  </div>
  {{if .Resources}}
  <h2>Drifted Resources</h2>
  <table>
    <thead><tr><th>Type</th><th>ID</th><th>Status</th></tr></thead>
    <tbody>
    {{range .Resources}}
      <tr>
        <td>{{.Type}}</td>
        <td>{{.ID}}</td>
        <td><span class="badge badge-{{.Status | lower}}">{{.Status}}</span></td>
      </tr>
    {{end}}
    </tbody>
  </table>
  {{else}}
  <p>No drift detected. &#x2705;</p>
  {{end}}
</body>
</html>`

// ReportData holds all data passed to the HTML template.
type ReportData struct {
	Summary     model.Summary
	Resources   []model.DriftedResource
	GeneratedAt string
}

// Renderer writes an HTML drift report to the provided writer.
type Renderer struct {
	tmpl *template.Template
}

// New creates a new Renderer, returning an error if the template fails to parse.
func New() (*Renderer, error) {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	tmpl, err := template.New("report").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return nil, err
	}
	return &Renderer{tmpl: tmpl}, nil
}

// Render executes the HTML template with the given summary and drifted resources.
func (r *Renderer) Render(w io.Writer, summary model.Summary, resources []model.DriftedResource) error {
	data := ReportData{
		Summary:     summary,
		Resources:   resources,
		GeneratedAt: time.Now().UTC().Format(time.RFC1123),
	}
	return r.tmpl.Execute(w, data)
}
