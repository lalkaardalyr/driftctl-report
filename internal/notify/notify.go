// Package notify provides alerting integrations for drift scan results.
// It supports sending notifications when drift is detected above a threshold.
package notify

import (
	"fmt"
	"io"

	"github.com/owner/driftctl-report/internal/model"
)

// Level represents the severity of a drift notification.
type Level string

const (
	LevelInfo    Level = "info"
	LevelWarning Level = "warning"
	LevelCritical Level = "critical"
)

// Message holds the content of a drift notification.
type Message struct {
	Level   Level
	Subject string
	Body    string
}

// Notifier sends a drift notification message.
type Notifier interface {
	Send(msg Message) error
}

// Config controls when notifications are triggered.
type Config struct {
	// WarnThreshold triggers a warning when drifted resource count >= this value.
	WarnThreshold int
	// CriticalThreshold triggers a critical alert when drifted resource count >= this value.
	CriticalThreshold int
}

// DefaultConfig returns a sensible default notification config.
func DefaultConfig() Config {
	return Config{
		WarnThreshold:     1,
		CriticalThreshold: 10,
	}
}

// Evaluate determines the notification level and builds a Message from a ScanResult.
func Evaluate(result model.ScanResult, cfg Config) (Message, bool) {
	drifted := len(result.DriftedResources)
	if drifted == 0 {
		return Message{}, false
	}

	level := LevelInfo
	if drifted >= cfg.CriticalThreshold {
		level = LevelCritical
	} else if drifted >= cfg.WarnThreshold {
		level = LevelWarning
	}

	return Message{
		Level:   level,
		Subject: fmt.Sprintf("[driftctl] %s: %d drifted resource(s) detected", level, drifted),
		Body:    buildBody(result),
	}, true
}

func buildBody(result model.ScanResult) string {
	s := result.Summary
	return fmt.Sprintf(
		"Coverage: %s | Managed: %d | Drifted: %d | Unmanaged: %d | Missing: %d",
		s.CoverageFormatted(),
		s.Managed,
		s.Drifted,
		s.Unmanaged,
		s.Missing,
	)
}

// WriteMessage formats a Message to the given writer.
func WriteMessage(w io.Writer, msg Message) error {
	_, err := fmt.Fprintf(w, "[%s] %s\n%s\n", msg.Level, msg.Subject, msg.Body)
	return err
}
