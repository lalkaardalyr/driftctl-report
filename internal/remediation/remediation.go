// Package remediation provides suggestions for resolving infrastructure drift.
package remediation

import (
	"fmt"
	"strings"

	"github.com/snyk/driftctl-report/internal/model"
)

// Severity indicates how urgent a remediation action is.
type Severity string

const (
	SeverityHigh   Severity = "high"
	SeverityMedium Severity = "medium"
	SeverityLow    Severity = "low"
)

// Action represents a single remediation suggestion for a drifted resource.
type Action struct {
	ResourceID   string
	ResourceType string
	Severity     Severity
	Description  string
	Command      string
}

// Plan holds all remediation actions derived from a scan result.
type Plan struct {
	Actions []Action
}

// Build constructs a remediation Plan from a ScanResult.
func Build(result model.ScanResult) Plan {
	var actions []Action
	for _, r := range result.DriftedResources {
		actions = append(actions, Action{
			ResourceID:   r.ResourceID,
			ResourceType: r.ResourceType,
			Severity:     severityFor(r.ResourceType),
			Description:  fmt.Sprintf("Resource %q has drifted from its desired state.", r.ResourceID),
			Command:      importCommand(r.ResourceType, r.ResourceID),
		})
	}
	for _, r := range result.UnmanagedResources {
		actions = append(actions, Action{
			ResourceID:   r.ResourceID,
			ResourceType: r.ResourceType,
			Severity:     SeverityMedium,
			Description:  fmt.Sprintf("Resource %q is unmanaged by IaC.", r.ResourceID),
			Command:      importCommand(r.ResourceType, r.ResourceID),
		})
	}
	return Plan{Actions: actions}
}

func severityFor(resourceType string) Severity {
	switch {
	case strings.HasPrefix(resourceType, "aws_iam"),
		strings.HasPrefix(resourceType, "aws_security_group"):
		return SeverityHigh
	case strings.HasPrefix(resourceType, "aws_s3"),
		strings.HasPrefix(resourceType, "aws_kms"):
		return SeverityHigh
	default:
		return SeverityMedium
	}
}

func importCommand(resourceType, resourceID string) string {
	return fmt.Sprintf("terraform import %s %s", resourceType+".this", resourceID)
}
