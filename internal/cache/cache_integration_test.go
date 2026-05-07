package cache_test

import (
	"os"
	"testing"

	"github.com/example/driftctl-report/internal/cache"
	"github.com/example/driftctl-report/internal/model"
)

// TestCache_StaleOnFileChange verifies that modifying the input file produces
// a cache miss even when a prior entry exists.
func TestCache_StaleOnFileChange(t *testing.T) {
	c, err := cache.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	// Write initial input file and populate cache.
	input := writeInputFile(t)
	want := makeScanResult()
	if err := c.Put(input, want); err != nil {
		t.Fatalf("Put: %v", err)
	}

	// Modify the file (change mtime via truncate + write).
	f, err := os.OpenFile(input, os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(`{"updated": true}`)
	_ = f.Close()

	// Cache key changes because mtime changed — expect a miss.
	got, err := c.Get(input)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != nil {
		t.Errorf("expected cache miss after file change, got %+v", got)
	}
}

// TestInvalidate_IdempotentOnMissingEntry ensures calling Invalidate on an
// entry that does not exist returns no error.
func TestInvalidate_IdempotentOnMissingEntry(t *testing.T) {
	c, _ := cache.New(t.TempDir())
	input := writeInputFile(t)
	if err := c.Invalidate(input); err != nil {
		t.Errorf("expected no error on missing entry, got: %v", err)
	}
}

// TestPut_OverwritesExistingEntry confirms that calling Put twice with the
// same input path and unchanged mtime overwrites the previous value.
func TestPut_OverwritesExistingEntry(t *testing.T) {
	c, _ := cache.New(t.TempDir())
	input := writeInputFile(t)

	first := makeScanResult()
	first.Summary.TotalResources = 5
	_ = c.Put(input, first)

	second := makeScanResult()
	second.Summary.TotalResources = 99
	_ = c.Put(input, second)

	got, err := c.Get(input)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("expected result, got nil")
	}
	if got.Summary.TotalResources != 99 {
		t.Errorf("expected TotalResources=99, got %d", got.Summary.TotalResources)
	}
}

// TestNew_InvalidDir checks that New returns an error when the directory
// cannot be created (e.g. parent is a file).
func TestNew_InvalidDir(t *testing.T) {
	// Create a regular file, then try to use it as a cache directory.
	f, err := os.CreateTemp(t.TempDir(), "notadir")
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	_, err = cache.New(f.Name() + "/subdir")
	if err == nil {
		t.Error("expected error when parent path is a file")
	}
}

// Compile-time check: model.ScanResult is used across test files.
var _ model.ScanResult
