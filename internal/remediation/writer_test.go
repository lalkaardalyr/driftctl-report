package remediation_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/snyk/driftctl-report/internal/remediation"
)

func samplePlan() remediation.Plan {
	return remediation.Plan{
		Actions: []remediation.Action{
			{
				ResourceID:   "sg-001",
				ResourceType: "aws_security_group",
				Severity:     remediation.SeverityHigh,
				Description:  "Resource \"sg-001\" has drifted.",
				Command:      "terraform import aws_security_group.this sg-001",
			},
		},
	}
}

func TestWriteJSON_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	if err := remediation.WriteJSON(&buf, samplePlan()); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var decoded struct {
		Actions []remediation.Action `json:"Actions"`
	}
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(decoded.Actions) != 1 {
		t.Errorf("expected 1 action, got %d", len(decoded.Actions))
	}
}

func TestWriteText_ContainsCommand(t *testing.T) {
	var buf bytes.Buffer
	if err := remediation.WriteText(&buf, samplePlan()); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	if !strings.Contains(buf.String(), "terraform import") {
		t.Errorf("expected terraform import command in output")
	}
}

func TestWriteText_EmptyPlan_NoDriftMessage(t *testing.T) {
	var buf bytes.Buffer
	if err := remediation.WriteText(&buf, remediation.Plan{}); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	if !strings.Contains(buf.String(), "No remediation actions required") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWriteText_SeverityLabel(t *testing.T) {
	var buf bytes.Buffer
	_ = remediation.WriteText(&buf, samplePlan())
	if !strings.Contains(buf.String(), "HIGH") {
		t.Errorf("expected HIGH severity label in output")
	}
}
