// Package model defines the domain types used throughout driftctl-report.
//
// ScanResult is the central struct that carries all information extracted from
// a driftctl JSON output file. It is produced by the parser package and
// consumed by the renderer package to generate HTML reports.
//
// Builder functions (FromAnalysis) convert raw driftctl analyser types into
// the normalised ScanResult representation so that the rest of the application
// remains decoupled from the upstream driftctl library types.
package model
