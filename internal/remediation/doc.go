// Package remediation analyses a driftctl ScanResult and produces an
// actionable remediation plan.
//
// Each drifted or unmanaged resource is mapped to an Action that includes:
//   - a human-readable description of the problem
//   - a severity level (high / medium / low) derived from the resource type
//   - a suggested Terraform import command to bring the resource back under
//     IaC management
//
// Plans can be serialised to JSON (WriteJSON) or plain text (WriteText) for
// inclusion in reports, CI pipelines, or notification payloads.
package remediation
