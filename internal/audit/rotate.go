package audit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// RotatingFile is an io.WriteCloser that writes to a daily-rotated log file
// under a given directory. Files are named audit-YYYY-MM-DD.jsonl.
type RotatingFile struct {
	dir  string
	file *os.File
	day  string // YYYY-MM-DD of the currently open file
}

// OpenRotating opens (or creates) today's audit log file in dir.
func OpenRotating(dir string) (*RotatingFile, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("audit: create log dir: %w", err)
	}
	rf := &RotatingFile{dir: dir}
	if err := rf.rotate(); err != nil {
		return nil, err
	}
	return rf, nil
}

// Write implements io.Writer. It rotates the underlying file when the calendar
// day changes.
func (rf *RotatingFile) Write(p []byte) (int, error) {
	today := today()
	if today != rf.day {
		if err := rf.rotate(); err != nil {
			return 0, err
		}
	}
	return rf.file.Write(p)
}

// Close closes the underlying file.
func (rf *RotatingFile) Close() error {
	if rf.file != nil {
		return rf.file.Close()
	}
	return nil
}

// CurrentPath returns the absolute path of the currently open log file.
func (rf *RotatingFile) CurrentPath() string {
	return filepath.Join(rf.dir, filename(rf.day))
}

func (rf *RotatingFile) rotate() error {
	if rf.file != nil {
		_ = rf.file.Close()
	}
	day := today()
	path := filepath.Join(rf.dir, filename(day))
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	rf.file = f
	rf.day = day
	return nil
}

func today() string {
	return time.Now().UTC().Format("2006-01-02")
}

func filename(day string) string {
	return fmt.Sprintf("audit-%s.jsonl", day)
}

// MultiWriter returns an io.Writer that duplicates writes to all provided writers,
// stopping at the first error. Useful for writing audit events to both a file
// and stdout simultaneously.
func MultiWriter(writers ...io.Writer) io.Writer {
	return &multiWriter{writers: writers}
}

type multiWriter struct{ writers []io.Writer }

func (m *multiWriter) Write(p []byte) (int, error) {
	for _, w := range m.writers {
		if _, err := w.Write(p); err != nil {
			return 0, err
		}
	}
	return len(p), nil
}
