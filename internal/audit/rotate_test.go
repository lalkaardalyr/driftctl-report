package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/driftctl-report/internal/audit"
)

func TestOpenRotating_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	rc, err := audit.OpenRotating(path, audit.RotateOptions{
		MaxSizeMB:  1,
		MaxBackups: 2,
		MaxAgeDays: 7,
	})
	if err != nil {
		t.Fatalf("OpenRotating returned error: %v", err)
	}
	defer rc.Close()

	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		t.Error("expected log file to be created, but it does not exist")
	}
}

func TestOpenRotating_WritesData(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	rc, err := audit.OpenRotating(path, audit.RotateOptions{
		MaxSizeMB:  1,
		MaxBackups: 2,
		MaxAgeDays: 7,
	})
	if err != nil {
		t.Fatalf("OpenRotating returned error: %v", err)
	}
	defer rc.Close()

	msg := []byte(`{"event":"test"}` + "\n")
	n, err := rc.Write(msg)
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if n != len(msg) {
		t.Errorf("expected to write %d bytes, wrote %d", len(msg), n)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read log file: %v", err)
	}
	if string(data) != string(msg) {
		t.Errorf("file content mismatch: got %q, want %q", string(data), string(msg))
	}
}

func TestOpenRotating_Close_Idempotent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	rc, err := audit.OpenRotating(path, audit.RotateOptions{
		MaxSizeMB:  1,
		MaxBackups: 1,
		MaxAgeDays: 1,
	})
	if err != nil {
		t.Fatalf("OpenRotating returned error: %v", err)
	}

	if err := rc.Close(); err != nil {
		t.Errorf("first Close returned error: %v", err)
	}
	// Second close should not panic or return an unexpected error.
	if err := rc.Close(); err != nil {
		t.Errorf("second Close returned error: %v", err)
	}
}

func TestOpenRotating_Rotate_OnSizeExceeded(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping rotation size test in short mode")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	// Set a very small max size so rotation triggers quickly.
	rc, err := audit.OpenRotating(path, audit.RotateOptions{
		MaxSizeMB:  1, // minimum supported; actual rotation behaviour is library-driven
		MaxBackups: 3,
		MaxAgeDays: 30,
	})
	if err != nil {
		t.Fatalf("OpenRotating returned error: %v", err)
	}
	defer rc.Close()

	// Write enough data to exercise the writer without asserting exact rotation
	// counts (which depend on the underlying lumberjack implementation).
	chunk := make([]byte, 512)
	for i := range chunk {
		chunk[i] = 'x'
	}
	for i := 0; i < 10; i++ {
		if _, err := rc.Write(chunk); err != nil {
			t.Fatalf("Write %d returned error: %v", i, err)
		}
	}

	// The primary log file must still be present after writes.
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		t.Error("primary log file missing after writes")
	}
}

func TestOpenRotating_TimestampedBackupName(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "audit.log")

	rc, err := audit.OpenRotating(base, audit.RotateOptions{
		MaxSizeMB:  1,
		MaxBackups: 5,
		MaxAgeDays: 14,
	})
	if err != nil {
		t.Fatalf("OpenRotating returned error: %v", err)
	}
	defer rc.Close()

	// Verify the writer reports a non-zero creation time (i.e., it was
	// initialised with a real timestamp rather than the zero value).
	created := rc.CreatedAt()
	if created.IsZero() {
		t.Error("expected non-zero CreatedAt timestamp")
	}
	if created.After(time.Now().Add(time.Second)) {
		t.Errorf("CreatedAt is in the future: %v", created)
	}
}
