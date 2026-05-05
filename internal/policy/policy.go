// Package policy provides drift policy evaluation — defining thresholds
// and rules that determine whether a scan result should be considered
// compliant, at-risk, or in violation.
package policy

import (
	"fmt"

	"github.com/example/driftctl-report/internal/model"
)

// Severity represents the outcome level of a policy evaluation.
type Severity string

const (
	SeverityOK       Severity = "ok"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// Rule defines a single policy threshold.
type Rule struct {
	// MaxDriftedResources is the maximum number of drifted resources allowed.
	// A value of -1 means no limit.
	MaxDriftedResources int
	// MaxUndocumentedResources is the maximum number of unmanaged resources allowed.
	MaxUndocumentedResources int
	// MinCoveragePercent is the minimum acceptable coverage percentage (0–100).
	MinCoveragePercent float64
}

// Result holds the outcome of evaluating a scan against a policy.
type Result struct {
	Severity Severity
	Violations []string
}

// DefaultRule returns a sensible default policy rule.
func DefaultRule() Rule {
	return Rule{
		MaxDriftedResources:      0,
		MaxUndocumentedResources: 5,
		MinCoveragePercent:       80.0,
	}
}

// Evaluate checks a ScanResult against the given Rule and returns a Result.
func Evaluate(r model.ScanResult, rule Rule) Result {
	var violations []string

	if rule.MaxDriftedResources >= 0 && len(r.DriftedResources) > rule.MaxDriftedResources {
		violations = append(violations, fmt.Sprintf(
			"drifted resources (%d) exceeds maximum allowed (%d)",
			len(r.DriftedResources), rule.MaxDriftedResources,
		))
	}

	if rule.MaxUndocumentedResources >= 0 && len(r.UnmanagedResources) > rule.MaxUndocumentedResources {
		violations = append(violations, fmt.Sprintf(
			"unmanaged resources (%d) exceeds maximum allowed (%d)",
			len(r.UnmanagedResources), rule.MaxUndocumentedResources,
		))
	}

	if r.Summary.Coverage < rule.MinCoveragePercent {
		violations = append(violations, fmt.Sprintf(
			"coverage (%.1f%%) is below minimum required (%.1f%%)",
			r.Summary.Coverage, rule.MinCoveragePercent,
		))
	}

	sev := SeverityOK
	if len(violations) > 0 {
		if len(r.DriftedResources) > 0 {
			sev = SeverityCritical
		} else {
			sev = SeverityWarning
		}
	}

	return Result{Severity: sev, Violations: violations}
}
