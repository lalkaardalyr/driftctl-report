package metrics

import (
	"sync"
	"time"

	"github.com/snyk/driftctl-report/internal/model"
)

// CollectedMetrics holds aggregated metrics across multiple scan results.
type CollectedMetrics struct {
	mu          sync.Mutex
	entries     []ScanMetrics
	LastUpdated time.Time
}

// ScanMetrics represents metrics captured from a single scan result.
type ScanMetrics struct {
	Timestamp       time.Time
	TotalResources  int
	DriftedCount    int
	UnmanagedCount  int
	Coverage        float64
	CoverageLabel   string
}

// NewCollector returns an initialised CollectedMetrics instance.
func NewCollector() *CollectedMetrics {
	return &CollectedMetrics{}
}

// Record appends metrics derived from the provided ScanResult.
func (c *CollectedMetrics) Record(result model.ScanResult) {
	m := FromScanResult(result)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = append(c.entries, ScanMetrics{
		Timestamp:      time.Now().UTC(),
		TotalResources: m.TotalResources,
		DriftedCount:   m.DriftedCount,
		UnmanagedCount: m.UnmanagedCount,
		Coverage:       m.Coverage,
		CoverageLabel:  m.CoverageLabel,
	})
	c.LastUpdated = time.Now().UTC()
}

// Snapshot returns a copy of all collected metric entries.
func (c *CollectedMetrics) Snapshot() []ScanMetrics {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]ScanMetrics, len(c.entries))
	copy(out, c.entries)
	return out
}

// Reset clears all collected entries.
func (c *CollectedMetrics) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = nil
	c.LastUpdated = time.Time{}
}

// Len returns the number of recorded entries.
func (c *CollectedMetrics) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}
