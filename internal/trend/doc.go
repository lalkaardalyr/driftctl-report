// Package trend analyses drift metrics across multiple driftctl scan results
// over time, enabling users to observe whether infrastructure drift is
// improving, worsening, or remaining stable between audit runs.
//
// Usage:
//
//	entries := []trend.TimestampedResult{
//		{Timestamp: t1, Result: result1},
//		{Timestamp: t2, Result: result2},
//	}
//	series := trend.FromResults(entries)
//	fmt.Println(series.Trend()) // "improving", "worsening", or "stable"
package trend
