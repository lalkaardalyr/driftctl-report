package schedule

import (
	"errors"
	"time"
)

// Config holds the configuration for a single scheduled scan job.
type Config struct {
	// Name is a human-readable identifier for the job.
	Name string `json:"name"`

	// InputPath is the path to the driftctl JSON output file.
	InputPath string `json:"input_path"`

	// OutputPath is the destination for the generated HTML report.
	OutputPath string `json:"output_path"`

	// IntervalSeconds defines how often the job runs.
	IntervalSeconds int `json:"interval_seconds"`
}

// Validate returns an error if the Config is missing required fields or
// contains invalid values.
func (c Config) Validate() error {
	if c.Name == "" {
		return errors.New("schedule config: name is required")
	}
	if c.InputPath == "" {
		return errors.New("schedule config: input_path is required")
	}
	if c.OutputPath == "" {
		return errors.New("schedule config: output_path is required")
	}
	if c.IntervalSeconds <= 0 {
		return errors.New("schedule config: interval_seconds must be positive")
	}
	return nil
}

// Interval converts IntervalSeconds to a time.Duration.
func (c Config) Interval() time.Duration {
	return time.Duration(c.IntervalSeconds) * time.Second
}
