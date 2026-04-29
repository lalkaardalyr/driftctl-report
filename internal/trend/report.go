package trend

import "fmt"

// Report summarises a Series into a human-readable text report.
type Report struct {
	Series Series
}

// NewReport creates a Report from the given Series.
func NewReport(s Series) Report {
	return Report{Series: s}
}

// Summary returns a concise text summary of the trend.
func (r Report) Summary() string {
	if len(r.Series) == 0 {
		return "No scan data available."
	}
	latest, _ := r.Series.Latest()
	direction := r.Series.Trend()
	return fmt.Sprintf(
		"Trend: %s | Scans: %d | Latest — Managed: %d, Drifted: %d, Unmanaged: %d, Coverage: %.1f%%",
		direction,
		len(r.Series),
		latest.TotalManaged,
		latest.Drifted,
		latest.Unmanaged,
		latest.Coverage*100,
	)
}

// DeltaDrifted returns the change in drifted resource count from first to last scan.
// Returns 0 if fewer than two data points exist.
func (r Report) DeltaDrifted() int {
	if len(r.Series) < 2 {
		return 0
	}
	return r.Series[len(r.Series)-1].Drifted - r.Series[0].Drifted
}

// DeltaCoverage returns the change in coverage percentage from first to last scan.
// Returns 0 if fewer than two data points exist.
func (r Report) DeltaCoverage() float64 {
	if len(r.Series) < 2 {
		return 0
	}
	return r.Series[len(r.Series)-1].Coverage - r.Series[0].Coverage
}
