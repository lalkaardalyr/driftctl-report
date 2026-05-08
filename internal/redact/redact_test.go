package redact_test

import (
	"testing"

	"github.com/snyk/driftctl-report/internal/model"
	"github.com/snyk/driftctl-report/internal/redact"
)

func makeResult(drifted []model.Resource, unmanaged []model.Resource) model.ScanResult {
	return model.ScanResult{
		DriftedResources:   drifted,
		UnmanagedResources: unmanaged,
	}
}

func TestApply_RedactsSensitiveAttrs(t *testing.T) {
	r := redact.New(redact.DefaultSensitiveKeys)
	res := makeResult([]model.Resource{
		{ID: "r1", Type: "aws_iam_user", Attrs: map[string]string{
			"name":     "alice",
			"password": "s3cr3t",
		}},
	}, nil)

	out := r.Apply(res)
	attrs := out.DriftedResources[0].Attrs

	if attrs["name"] != "alice" {
		t.Errorf("expected name to be unchanged, got %q", attrs["name"])
	}
	if attrs["password"] != "[REDACTED]" {
		t.Errorf("expected password to be redacted, got %q", attrs["password"])
	}
}

func TestApply_RedactsTokenAndAPIKey(t *testing.T) {
	r := redact.New(redact.DefaultSensitiveKeys)
	res := makeResult([]model.Resource{
		{ID: "r2", Type: "aws_lambda_function", Attrs: map[string]string{
			"api_key":      "key-abc",
			"access_token": "tok-xyz",
			"description":  "safe",
		}},
	}, nil)

	out := r.Apply(res)
	attrs := out.DriftedResources[0].Attrs

	if attrs["api_key"] != "[REDACTED]" {
		t.Errorf("expected api_key redacted, got %q", attrs["api_key"])
	}
	if attrs["access_token"] != "[REDACTED]" {
		t.Errorf("expected access_token redacted, got %q", attrs["access_token"])
	}
	if attrs["description"] != "safe" {
		t.Errorf("expected description unchanged, got %q", attrs["description"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	r := redact.New(redact.DefaultSensitiveKeys)
	orig := makeResult([]model.Resource{
		{ID: "r3", Type: "aws_s3_bucket", Attrs: map[string]string{
			"secret_key": "original-secret",
		}},
	}, nil)

	r.Apply(orig)

	if orig.DriftedResources[0].Attrs["secret_key"] != "original-secret" {
		t.Error("Apply must not mutate the original result")
	}
}

func TestApply_NilResourceLists(t *testing.T) {
	r := redact.New(redact.DefaultSensitiveKeys)
	res := makeResult(nil, nil)
	out := r.Apply(res)

	if out.DriftedResources != nil || out.UnmanagedResources != nil {
		t.Error("nil slices should remain nil after redaction")
	}
}

func TestApply_CustomKeys(t *testing.T) {
	r := redact.New([]string{"internal_id"})
	res := makeResult([]model.Resource{
		{ID: "r4", Type: "custom_resource", Attrs: map[string]string{
			"internal_id": "should-be-redacted",
			"password":    "should-not-be-redacted",
		}},
	}, nil)

	out := r.Apply(res)
	attrs := out.DriftedResources[0].Attrs

	if attrs["internal_id"] != "[REDACTED]" {
		t.Errorf("expected internal_id redacted, got %q", attrs["internal_id"])
	}
	if attrs["password"] != "should-not-be-redacted" {
		t.Errorf("expected password unchanged with custom keys, got %q", attrs["password"])
	}
}
