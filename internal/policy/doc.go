// Package policy implements configurable drift policy evaluation for
// driftctl-report. A policy Rule defines numeric thresholds for drifted
// resources, unmanaged resources, and minimum infrastructure coverage.
//
// Use Evaluate to check a model.ScanResult against a Rule and receive a
// Result that indicates the Severity (ok / warning / critical) and a list
// of human-readable violation messages suitable for reporting or gating
// CI pipelines.
//
// Example:
//
//	rule := policy.DefaultRule()
//	out  := policy.Evaluate(scanResult, rule)
//	if out.Severity == policy.SeverityCritical {
//		os.Exit(1)
//	}
package policy
