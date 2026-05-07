package summary

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteJSON serialises the Aggregate to w as indented JSON.
func WriteJSON(w io.Writer, agg Aggregate) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(agg)
}

// WriteText writes a human-readable tabular summary of the Aggregate to w.
func WriteText(w io.Writer, agg Aggregate) error {
	if agg.TotalScans == 0 {
		_, err := fmt.Fprintln(w, "No scan data available.")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	lines := []struct {
		label string
		value string
	}{
		{"Total scans", fmt.Sprintf("%d", agg.TotalScans)},
		{"First scan", agg.FirstScan.Format("2006-01-02 15:04:05 UTC")},
		{"Last scan", agg.LastScan.Format("2006-01-02 15:04:05 UTC")},
		{"Avg drifted resources", fmt.Sprintf("%.1f", agg.AvgDrifted)},
		{"Avg unmanaged resources", fmt.Sprintf("%.1f", agg.AvgUnmanaged)},
		{"Avg coverage", fmt.Sprintf("%.1f%%", agg.AvgCoverage)},
		{"Max drifted (single scan)", fmt.Sprintf("%d", agg.MaxDrifted)},
		{"Max unmanaged (single scan)", fmt.Sprintf("%d", agg.MaxUnmanaged)},
		{"Best coverage", fmt.Sprintf("%.1f%%", agg.BestCoverage)},
		{"Worst coverage", fmt.Sprintf("%.1f%%", agg.WorstCoverage)},
	}

	for _, l := range lines {
		if _, err := fmt.Fprintf(tw, "%s:\t%s\n", l.label, l.value); err != nil {
			return err
		}
	}

	return tw.Flush()
}
