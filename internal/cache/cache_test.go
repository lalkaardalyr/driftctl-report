package cache_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/driftctl-report/internal/cache"
	"github.com/example/driftctl-report/internal/model"
)

func makeScanResult() model.ScanResult {
	return model.ScanResult{
		Summary: model.Summary{
			TotalResources:    10,
			ManagedResources:  8,
			DriftedResources:  2,
			UnmanagedResources: 1,
			Coverage:          80.0,
		},
	}
}

func writeInputFile(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "input-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	return f.Name()
}

func TestNew_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cache")
	_, err := cache.New(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("expected cache directory to exist: %v", err)
	}
}

func TestGet_MissReturnsNil(t *testing.T) {
	c, _ := cache.New(t.TempDir())
	input := writeInputFile(t)
	result, err := c.Get(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil on cache miss, got %+v", result)
	}
}

func TestPutAndGet_RoundTrip(t *testing.T) {
	c, _ := cache.New(t.TempDir())
	input := writeInputFile(t)
	want := makeScanResult()

	if err := c.Put(input, want); err != nil {
		t.Fatalf("Put: %v", err)
	}
	got, err := c.Get(input)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got == nil {
		t.Fatal("expected cached result, got nil")
	}
	if got.Summary.TotalResources != want.Summary.TotalResources {
		t.Errorf("TotalResources: got %d, want %d", got.Summary.TotalResources, want.Summary.TotalResources)
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c, _ := cache.New(t.TempDir())
	input := writeInputFile(t)
	_ = c.Put(input, makeScanResult())

	if err := c.Invalidate(input); err != nil {
		t.Fatalf("Invalidate: %v", err)
	}
	got, err := c.Get(input)
	if err != nil {
		t.Fatalf("Get after invalidate: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil after invalidation, got %+v", got)
	}
}

func TestPut_WritesValidJSON(t *testing.T) {
	cacheDir := t.TempDir()
	c, _ := cache.New(cacheDir)
	input := writeInputFile(t)
	_ = c.Put(input, makeScanResult())

	entries, _ := os.ReadDir(cacheDir)
	if len(entries) != 1 {
		t.Fatalf("expected 1 cache file, got %d", len(entries))
	}
	data, _ := os.ReadFile(filepath.Join(cacheDir, entries[0].Name()))
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Errorf("cache file is not valid JSON: %v", err)
	}
	if _, ok := m["cached_at"]; !ok {
		t.Error("expected 'cached_at' field in cache entry")
	}
	_ = time.Now() // satisfy import
}
