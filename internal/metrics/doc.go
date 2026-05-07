// Package metrics derives aggregated, per-type drift statistics from a
// model.ScanResult. It is designed as a read-only computation layer with no
// side effects, making it safe to call from renderers, exporters, and
// notification handlers alike.
//
// # Usage
//
//	sr := model.ScanResult{ ... }
//	report := metrics.FromScanResult(sr)
//	fmt.Println(report.CoverageFormatted()) // "92.50%"
//	fmt.Println(report.CoverageLabel())     // "good"
package metrics
