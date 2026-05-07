// Package summary provides multi-scan aggregation utilities for driftctl-report.
//
// It accepts a slice of [Entry] values — each pairing a scan timestamp with a
// [model.ScanResult] — and produces an [Aggregate] containing rolled-up
// statistics such as average/max drifted resource counts and coverage range.
//
// Two output helpers are provided:
//
//   - [WriteJSON] serialises the Aggregate as indented JSON.
//   - [WriteText] renders a human-readable tabular summary via tabwriter.
//
// Typical usage:
//
//	entries := []summary.Entry{
//		{ScannedAt: t1, Result: result1},
//		{ScannedAt: t2, Result: result2},
//	}
//	agg := summary.Compute(entries)
//	summary.WriteText(os.Stdout, agg)
package summary
