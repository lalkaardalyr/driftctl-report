package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftctl-report/internal/model"
	"github.com/driftctl-report/internal/snapshot"
)

func makeScanResult(drifted, unmanaged int) model.ScanResult {
	resources := make([]model.Resource, drifted)
	for i := range resources {
		resources[i] = model.Resource{ID: fmt.Sprintf("res-%d", i), Type: "aws_s3_bucket", Source: "driftctl"}
	}
	return model.ScanResult{
		Summary: model.Summary{TotalResources: drifted + unmanaged, TotalDrifted: drifted, TotalUnmanaged: unmanaged},
		Drifted: resources,
	}
}

func TestNewStore_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "snapshots")
	_, err := snapshot.NewStore(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestSave_CreatesFile(t *testing.T) {
	store, _ := snapshot.NewStore(t.TempDir())
	result := model.ScanResult{Summary: model.Summary{TotalResources: 5, TotalDrifted: 1}}
	entry, err := store.Save(result, "ci-run")
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if entry.ID == "" {
		t.Error("expected non-empty ID")
	}
	if entry.Label != "ci-run" {
		t.Errorf("expected label 'ci-run', got %q", entry.Label)
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	store, _ := snapshot.NewStore(t.TempDir())
	result := model.ScanResult{Summary: model.Summary{TotalResources: 10, TotalDrifted: 3}}
	saved, _ := store.Save(result, "test")
	loaded, err := store.Load(saved.ID)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.ID != saved.ID {
		t.Errorf("ID mismatch: got %q, want %q", loaded.ID, saved.ID)
	}
	if loaded.Result.Summary.TotalDrifted != 3 {
		t.Errorf("TotalDrifted mismatch: got %d", loaded.Result.Summary.TotalDrifted)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	store, _ := snapshot.NewStore(t.TempDir())
	_, err := store.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestList_ReturnsAllEntries(t *testing.T) {
	store, _ := snapshot.NewStore(t.TempDir())
	result := model.ScanResult{Summary: model.Summary{TotalResources: 2}}
	for i := 0; i < 3; i++ {
		if _, err := store.Save(result, ""); err != nil {
			t.Fatalf("Save %d failed: %v", i, err)
		}
	}
	entries, err := store.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}
