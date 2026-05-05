package policy_test

import (
	"testing"

	"github.com/example/driftctl-report/internal/model"
	"github.com/example/driftctl-report/internal/policy"
)

func makeResult(drifted, unmanaged int, coverage float64) model.ScanResult {
	dr := make([]model.Resource, drifted)
	for i := range dr {
		dr[i] = model.Resource{ID: fmt.Sprintf("res-%d", i), Type: "aws_s3_bucket"}
	}
	ur := make([]model.Resource, unmanaged)
	for i := range ur {
		ur[i] = model.Resource{ID: fmt.Sprintf("unmanaged-%d", i), Type: "aws_instance"}
	}
	return model.ScanResult{
		DriftedResources:   dr,
		UnmanagedResources: ur,
		Summary: model.Summary{
			Coverage: coverage,
		},
	}
}

func TestEvaluate_NoDrift_OK(t *testing.T) {
	r := makeResult(0, 0, 95.0)
	rule := policy.DefaultRule()
	out := policy.Evaluate(r, rule)
	if out.Severity != policy.SeverityOK {
		t.Errorf("expected OK, got %s", out.Severity)
	}
	if len(out.Violations) != 0 {
		t.Errorf("expected no violations, got %v", out.Violations)
	}
}

func TestEvaluate_DriftedResources_Critical(t *testing.T) {
	r := makeResult(3, 0, 90.0)
	rule := policy.DefaultRule() // MaxDriftedResources = 0
	out := policy.Evaluate(r, rule)
	if out.Severity != policy.SeverityCritical {
		t.Errorf("expected Critical, got %s", out.Severity)
	}
	if len(out.Violations) == 0 {
		t.Error("expected at least one violation")
	}
}

func TestEvaluate_LowCoverage_Warning(t *testing.T) {
	r := makeResult(0, 2, 60.0)
	rule := policy.DefaultRule()
	out := policy.Evaluate(r, rule)
	if out.Severity != policy.SeverityWarning {
		t.Errorf("expected Warning, got %s", out.Severity)
	}
}

func TestEvaluate_TooManyUnmanaged_Warning(t *testing.T) {
	r := makeResult(0, 10, 85.0)
	rule := policy.DefaultRule() // MaxUndocumentedResources = 5
	out := policy.Evaluate(r, rule)
	if out.Severity == policy.SeverityOK {
		t.Error("expected non-OK severity for too many unmanaged resources")
	}
}

func TestEvaluate_UnlimitedDrift_NoViolation(t *testing.T) {
	r := makeResult(100, 0, 95.0)
	rule := policy.Rule{
		MaxDriftedResources:      -1,
		MaxUndocumentedResources: -1,
		MinCoveragePercent:       0,
	}
	out := policy.Evaluate(r, rule)
	if out.Severity != policy.SeverityOK {
		t.Errorf("expected OK with unlimited rule, got %s", out.Severity)
	}
}

func TestDefaultRule_Values(t *testing.T) {
	rule := policy.DefaultRule()
	if rule.MaxDriftedResources != 0 {
		t.Errorf("expected MaxDriftedResources=0, got %d", rule.MaxDriftedResources)
	}
	if rule.MinCoveragePercent != 80.0 {
		t.Errorf("expected MinCoveragePercent=80.0, got %.1f", rule.MinCoveragePercent)
	}
}
