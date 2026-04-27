package render

import "github.com/owner/driftctl-report/internal/model"

// TemplateData holds all data passed to the HTML template during rendering.
type TemplateData struct {
	Summary            SummaryView
	DriftedResources   []ResourceView
	MissingResources   []ResourceView
	UnmanagedResources []ResourceView
}

// SummaryView is a presentation-layer representation of model.Summary.
type SummaryView struct {
	CoverageFormatted string
	CoverageLabel     string
	CoverageClass     string
	TotalResources    int
	DriftedCount      int
	MissingCount      int
	UnmanagedCount    int
}

// ResourceView is a flat, template-friendly representation of a resource.
type ResourceView struct {
	Type string
	ID   string
}

// NewTemplateData converts a model.ScanResult into TemplateData ready for rendering.
func NewTemplateData(result model.ScanResult) TemplateData {
	var coverageClass string
	switch result.Summary.CoverageLabel() {
	case "Good":
		coverageClass = "coverage-ok"
	case "Fair":
		coverageClass = "coverage-warn"
	default:
		coverageClass = "coverage-bad"
	}

	return TemplateData{
		Summary: SummaryView{
			CoverageFormatted: result.Summary.CoverageFormatted(),
			CoverageLabel:     result.Summary.CoverageLabel(),
			CoverageClass:     coverageClass,
			TotalResources:    result.Summary.TotalResources,
			DriftedCount:      result.Summary.DriftedCount,
			MissingCount:      result.Summary.MissingCount,
			UnmanagedCount:    result.Summary.UnmanagedCount,
		},
		DriftedResources:   toResourceViews(result.DriftedResources),
		MissingResources:   toResourceViews(result.MissingResources),
		UnmanagedResources: toResourceViews(result.UnmanagedResources),
	}
}

func toResourceViews(resources []model.Resource) []ResourceView {
	views := make([]ResourceView, 0, len(resources))
	for _, r := range resources {
		views = append(views, ResourceView{Type: r.Type, ID: r.ID})
	}
	return views
}
