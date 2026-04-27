// Package diff compares two driftctl ScanResult values to identify
// infrastructure changes between audit runs.
//
// Use diff.Compare to obtain a Result that categorises resources as:
//   - NewlyDrifted  – resources that have drifted since the baseline scan.
//   - Resolved      – resources that were drifted but are now in sync.
//   - StillDrifted  – resources that remain drifted across both scans.
//   - NewlyUnmanaged – resources that appeared unmanaged since the baseline.
//
// Example:
//
//	result := diff.Compare(baselineScan, currentScan)
//	fmt.Println("Newly drifted:", len(result.NewlyDrifted))
package diff
