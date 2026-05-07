package snapshot_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/driftctl-report/internal/model"
	"github.com/driftctl-report/internal/snapshot"
)

func makeEntry(drifted []model.Resource) snapshot.Entry {
	return snapshot.Entry{
		ID:        fmt.Sprintf("entry-%d", time.Now().UnixNano()),
		CreatedAt: time.Now().UTC(),
		Result: model.ScanResult{
			Summary: model.Summary{
				TotalResources: 10,
				TotalDrifted:   len(drifted),
			},
			Drifted: drifted,
		},
	}
}

func res(id string) model.Resource {
	return model.Resource{ID: id, Type: "aws_s3_bucket", Source: "driftctl"}
}

func TestCompare_DriftedChange_Positive(t *testing.T) {
	from := makeEntry([]model.Resource{res("a")})
	to := makeEntry([]model.Resource{res("a"), res("b")})
	delta := snapshot.Compare(from, to)
	if delta.DriftedChange != 1 {
		t.Errorf("expected DriftedChange=1, got %d", delta.DriftedChange)
	}
}

func TestCompare_DriftedChange_Negative(t *testing.T) {
	from := makeEntry([]model.Resource{res("a"), res("b")})
	to := makeEntry([]model.Resource{res("a")})
	delta := snapshot.Compare(from, to)
	if delta.DriftedChange != -1 {
		t.Errorf("expected DriftedChange=-1, got %d", delta.DriftedChange)
	}
}

func TestCompare_NewlyDrifted(t *testing.T) {
	from := makeEntry([]model.Resource{res("a")})
	to := makeEntry([]model.Resource{res("a"), res("b")})
	delta := snapshot.Compare(from, to)
	if len(delta.NewlyDrifted) != 1 || delta.NewlyDrifted[0].ID != "b" {
		t.Errorf("expected NewlyDrifted=[b], got %v", delta.NewlyDrifted)
	}
}

func TestCompare_Resolved(t *testing.T) {
	from := makeEntry([]model.Resource{res("a"), res("b")})
	to := makeEntry([]model.Resource{res("a")})
	delta := snapshot.Compare(from, to)
	if len(delta.Resolved) != 1 || delta.Resolved[0].ID != "b" {
		t.Errorf("expected Resolved=[b], got %v", delta.Resolved)
	}
}

func TestDelta_Summary_Improving(t *testing.T) {
	d := snapshot.Delta{DriftedChange: -2}
	if d.Summary() != "drift reduced by 2 resource(s)" {
		t.Errorf("unexpected summary: %s", d.Summary())
	}
}

func TestDelta_Summary_Stable(t *testing.T) {
	d := snapshot.Delta{DriftedChange: 0}
	if d.Summary() != "no change in drift" {
		t.Errorf("unexpected summary: %s", d.Summary())
	}
}
