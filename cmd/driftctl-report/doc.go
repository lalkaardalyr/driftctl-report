// Package main provides the driftctl-report CLI tool.
//
// Usage:
//
//	driftctl-report --input <driftctl-output.json> [--output <report.html>]
//
// Flags:
//
//	--input   Path to the driftctl JSON scan output file (required).
//	--output  Path for the generated HTML report (default: drift-report.html).
//
// The tool parses the driftctl JSON output, builds a structured summary of
// infrastructure drift, and renders a human-readable HTML report suitable
// for sharing with teams or including in audit artefacts.
package main
