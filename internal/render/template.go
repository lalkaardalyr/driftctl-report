package render

import "html/template"

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Driftctl Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 2rem; background: #f5f5f5; color: #333; }
        h1 { color: #2c3e50; }
        .summary { background: #fff; border-radius: 8px; padding: 1.5rem; margin-bottom: 2rem; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .summary-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 1rem; margin-top: 1rem; }
        .metric { text-align: center; padding: 1rem; border-radius: 6px; background: #ecf0f1; }
        .metric .value { font-size: 2rem; font-weight: bold; }
        .metric .label { font-size: 0.85rem; color: #666; }
        .coverage-ok { color: #27ae60; }
        .coverage-warn { color: #e67e22; }
        .coverage-bad { color: #e74c3c; }
        table { width: 100%; border-collapse: collapse; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        th { background: #2c3e50; color: #fff; padding: 0.75rem 1rem; text-align: left; }
        td { padding: 0.75rem 1rem; border-bottom: 1px solid #ecf0f1; }
        tr:last-child td { border-bottom: none; }
        tr:hover td { background: #f9f9f9; }
        .badge { display: inline-block; padding: 0.2rem 0.6rem; border-radius: 4px; font-size: 0.8rem; font-weight: bold; }
        .badge-drifted { background: #fde8e8; color: #c0392b; }
        .badge-missing { background: #fef9e7; color: #d68910; }
        .badge-unmanaged { background: #eaf4fb; color: #1a5276; }
        .section-title { margin-top: 2rem; margin-bottom: 1rem; color: #2c3e50; }
        .no-drift { color: #27ae60; font-style: italic; }
    </style>
</head>
<body>
    <h1>&#128202; Driftctl Infrastructure Report</h1>
    <div class="summary">
        <h2>Summary</h2>
        <div class="summary-grid">
            <div class="metric">
                <div class="value {{.Summary.CoverageClass}}">{{.Summary.CoverageFormatted}}</div>
                <div class="label">Coverage ({{.Summary.CoverageLabel}})</div>
            </div>
            <div class="metric">
                <div class="value">{{.Summary.TotalResources}}</div>
                <div class="label">Total Resources</div>
            </div>
            <div class="metric">
                <div class="value" style="color:#c0392b">{{.Summary.DriftedCount}}</div>
                <div class="label">Drifted</div>
            </div>
            <div class="metric">
                <div class="value" style="color:#d68910">{{.Summary.MissingCount}}</div>
                <div class="label">Missing</div>
            </div>
            <div class="metric">
                <div class="value" style="color:#1a5276">{{.Summary.UnmanagedCount}}</div>
                <div class="label">Unmanaged</div>
            </div>
        </div>
    </div>

    <h2 class="section-title">Drifted Resources</h2>
    {{if .DriftedResources}}
    <table>
        <thead><tr><th>Type</th><th>ID</th><th>Status</th></tr></thead>
        <tbody>
        {{range .DriftedResources}}
        <tr><td>{{.Type}}</td><td>{{.ID}}</td><td><span class="badge badge-drifted">drifted</span></td></tr>
        {{end}}
        </tbody>
    </table>
    {{else}}<p class="no-drift">&#10003; No drifted resources found.</p>{{end}}

    <h2 class="section-title">Missing Resources</h2>
    {{if .MissingResources}}
    <table>
        <thead><tr><th>Type</th><th>ID</th><th>Status</th></tr></thead>
        <tbody>
        {{range .MissingResources}}
        <tr><td>{{.Type}}</td><td>{{.ID}}</td><td><span class="badge badge-missing">missing</span></td></tr>
        {{end}}
        </tbody>
    </table>
    {{else}}<p class="no-drift">&#10003; No missing resources found.</p>{{end}}

    <h2 class="section-title">Unmanaged Resources</h2>
    {{if .UnmanagedResources}}
    <table>
        <thead><tr><th>Type</th><th>ID</th><th>Status</th></tr></thead>
        <tbody>
        {{range .UnmanagedResources}}
        <tr><td>{{.Type}}</td><td>{{.ID}}</td><td><span class="badge badge-unmanaged">unmanaged</span></td></tr>
        {{end}}
        </tbody>
    </table>
    {{else}}<p class="no-drift">&#10003; No unmanaged resources found.</p>{{end}}
</body>
</html>`

// ParsedTemplate is the compiled HTML template used by the renderer.
var ParsedTemplate = template.Must(template.New("report").Parse(htmlTemplate))
