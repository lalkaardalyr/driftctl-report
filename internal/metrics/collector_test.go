package metrics

import (
	"testing"

	"github.com/snyk/driftctl-report/internal/model"
)

func makeCollectorResult(drifted, unmanaged, managed int) model.ScanResult {
	resources := make([]model.Resource, 0, drifted+unmanaged+managed)
	for i := 0; i < drifted; i++ {
		resources = append(resources, model.Resource{ID: fmt.Sprintf("drifted-%d", i), Type: "aws_s3_bucket", Source: "drifted"})
	}
	for i := 0; i < unmanaged; i++ {
		resources = append(resources, model.Resource{ID: fmt.Sprintf("unmanaged-%d", i), Type: "aws_instance", Source: "unmanaged"})
	}
	for i := 0; i < managed; i++ {
		resources = append(resources, model.Resource{ID: fmt.Sprintf("managed-%d", i), Type: "aws_iam_role", Source: "managed"})
	}
	return model.ScanResult{
		DriftedResources:  resources[:drifted],
		UnmanagedResources: resources[drifted : drifted+unmanaged],
		ManagedResources:  resources[drifted+unmanaged:],
	}
}

func TestCollector_RecordIncreasesLen(t *testing.T) {
	c := NewCollector()
	if c.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", c.Len())
	}
	c.Record(makeCollectorResult(2, 1, 5))
	c.Record(makeCollectorResult(0, 0, 10))
	if c.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", c.Len())
	}
}

func TestCollector_SnapshotReturnsCopy(t *testing.T) {
	c := NewCollector()
	c.Record(makeCollectorResult(3, 2, 10))
	snap := c.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 snapshot entry, got %d", len(snap))
	}
	// Mutating snapshot must not affect collector
	snap[0].DriftedCount = 999
	if c.Snapshot()[0].DriftedCount == 999 {
		t.Fatal("snapshot mutation affected internal state")
	}
}

func TestCollector_Reset_ClearsEntries(t *testing.T) {
	c := NewCollector()
	c.Record(makeCollectorResult(1, 1, 5))
	c.Reset()
	if c.Len() != 0 {
		t.Fatalf("expected 0 entries after reset, got %d", c.Len())
	}
	if !c.LastUpdated.IsZero() {
		t.Fatal("expected LastUpdated to be zero after reset")
	}
}

func TestCollector_MetricsValues(t *testing.T) {
	c := NewCollector()
	c.Record(makeCollectorResult(2, 1, 7))
	snap := c.Snapshot()
	entry := snap[0]
	if entry.DriftedCount != 2 {
		t.Errorf("expected DriftedCount=2, got %d", entry.DriftedCount)
	}
	if entry.UnmanagedCount != 1 {
		t.Errorf("expected UnmanagedCount=1, got %d", entry.UnmanagedCount)
	}
	if entry.TotalResources != 10 {
		t.Errorf("expected TotalResources=10, got %d", entry.TotalResources)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero Timestamp")
	}
}
