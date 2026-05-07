package remediation_test

import (
	"strings"
	"testing"

	"github.com/snyk/driftctl-report/internal/model"
	"github.com/snyk/driftctl-report/internal/remediation"
)

func makeResult(drifted, unmanaged []model.Resource) model.ScanResult {
	return model.ScanResult{
		DriftedResources:   drifted,
		UnmanagedResources: unmanaged,
	}
}

func TestBuild_EmptyResult_NoPlan(t *testing.T) {
	plan := remediation.Build(makeResult(nil, nil))
	if len(plan.Actions) != 0 {
		t.Fatalf("expected 0 actions, got %d", len(plan.Actions))
	}
}

func TestBuild_DriftedResource_CreatesAction(t *testing.T) {
	result := makeResult([]model.Resource{
		{ResourceID: "sg-abc123", ResourceType: "aws_security_group"},
	}, nil)
	plan := remediation.Build(result)
	if len(plan.Actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(plan.Actions))
	}
	action := plan.Actions[0]
	if action.ResourceID != "sg-abc123" {
		t.Errorf("unexpected resource ID: %s", action.ResourceID)
	}
	if action.Severity != remediation.SeverityHigh {
		t.Errorf("expected high severity for security group, got %s", action.Severity)
	}
}

func TestBuild_UnmanagedResource_MediumSeverity(t *testing.T) {
	result := makeResult(nil, []model.Resource{
		{ResourceID: "i-0123456789", ResourceType: "aws_instance"},
	})
	plan := remediation.Build(result)
	if len(plan.Actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(plan.Actions))
	}
	if plan.Actions[0].Severity != remediation.SeverityMedium {
		t.Errorf("expected medium severity, got %s", plan.Actions[0].Severity)
	}
}

func TestBuild_CommandContainsResourceID(t *testing.T) {
	result := makeResult([]model.Resource{
		{ResourceID: "bucket-xyz", ResourceType: "aws_s3_bucket"},
	}, nil)
	plan := remediation.Build(result)
	if !strings.Contains(plan.Actions[0].Command, "bucket-xyz") {
		t.Errorf("expected command to contain resource ID, got: %s", plan.Actions[0].Command)
	}
}

func TestBuild_IAMResource_HighSeverity(t *testing.T) {
	result := makeResult([]model.Resource{
		{ResourceID: "arn:aws:iam::123:role/MyRole", ResourceType: "aws_iam_role"},
	}, nil)
	plan := remediation.Build(result)
	if plan.Actions[0].Severity != remediation.SeverityHigh {
		t.Errorf("expected high severity for IAM resource, got %s", plan.Actions[0].Severity)
	}
}
