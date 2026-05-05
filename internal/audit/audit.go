// Package audit provides structured audit logging for drift scan events,
// recording who triggered a scan, when, and what the outcome was.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// EventKind classifies the type of audit event.
type EventKind string

const (
	EventScanStarted   EventKind = "scan_started"
	EventScanCompleted EventKind = "scan_completed"
	EventScanFailed    EventKind = "scan_failed"
	EventReportExported EventKind = "report_exported"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp  time.Time         `json:"timestamp"`
	Kind       EventKind         `json:"kind"`
	Actor      string            `json:"actor,omitempty"`
	InputPath  string            `json:"input_path,omitempty"`
	OutputPath string            `json:"output_path,omitempty"`
	Drifted    int               `json:"drifted_count,omitempty"`
	Unmanaged  int               `json:"unmanaged_count,omitempty"`
	Meta       map[string]string `json:"meta,omitempty"`
	Error      string            `json:"error,omitempty"`
}

// Logger writes audit events as newline-delimited JSON to an io.Writer.
type Logger struct {
	w     io.Writer
	actor string
}

// New creates a new audit Logger writing to w.
// actor identifies the user or service triggering events (may be empty).
func New(w io.Writer, actor string) *Logger {
	return &Logger{w: w, actor: actor}
}

// Log writes a single audit event. It stamps the current UTC time.
func (l *Logger) Log(e Event) error {
	e.Timestamp = time.Now().UTC()
	if e.Actor == "" {
		e.Actor = l.actor
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal event: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}

// ScanStarted is a convenience method for logging a scan_started event.
func (l *Logger) ScanStarted(inputPath string) error {
	return l.Log(Event{Kind: EventScanStarted, InputPath: inputPath})
}

// ScanCompleted is a convenience method for logging a scan_completed event.
func (l *Logger) ScanCompleted(inputPath, outputPath string, drifted, unmanaged int) error {
	return l.Log(Event{
		Kind:       EventScanCompleted,
		InputPath:  inputPath,
		OutputPath: outputPath,
		Drifted:    drifted,
		Unmanaged:  unmanaged,
	})
}

// ScanFailed logs a scan_failed event with the provided error message.
func (l *Logger) ScanFailed(inputPath, errMsg string) error {
	return l.Log(Event{Kind: EventScanFailed, InputPath: inputPath, Error: errMsg})
}
