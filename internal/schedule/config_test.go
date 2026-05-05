package schedule_test

import (
	"testing"
	"time"

	"github.com/org/driftctl-report/internal/schedule"
)

func validConfig() schedule.Config {
	return schedule.Config{
		Name:            "nightly",
		InputPath:       "/data/drift.json",
		OutputPath:      "/reports/drift.html",
		IntervalSeconds: 3600,
	}
}

func TestConfig_Validate_Valid(t *testing.T) {
	if err := validConfig().Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestConfig_Validate_MissingName(t *testing.T) {
	c := validConfig()
	c.Name = ""
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestConfig_Validate_MissingInputPath(t *testing.T) {
	c := validConfig()
	c.InputPath = ""
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for missing input_path")
	}
}

func TestConfig_Validate_MissingOutputPath(t *testing.T) {
	c := validConfig()
	c.OutputPath = ""
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for missing output_path")
	}
}

func TestConfig_Validate_ZeroInterval(t *testing.T) {
	c := validConfig()
	c.IntervalSeconds = 0
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestConfig_Validate_NegativeInterval(t *testing.T) {
	c := validConfig()
	c.IntervalSeconds = -10
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for negative interval")
	}
}

func TestConfig_Interval_ConvertsToDuration(t *testing.T) {
	c := validConfig()
	c.IntervalSeconds = 120
	if got := c.Interval(); got != 120*time.Second {
		t.Errorf("expected 120s, got %s", got)
	}
}
