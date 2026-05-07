package snapshot

import (
	"fmt"

	"github.com/driftctl-report/internal/model"
)

// Delta describes the change between two snapshot entries.
type Delta struct {
	From          Entry
	To            Entry
	DriftedChange int // positive = more drift, negative = less drift
	ManagedChange int
	NewlyDrifted  []model.Resource
	Resolved      []model.Resource
}

// Summary returns a human-readable one-line description of the delta.
func (d Delta) Summary() string {
	switch {
	case d.DriftedChange < 0:
		return fmt.Sprintf("drift reduced by %d resource(s)", -d.DriftedChange)
	case d.DriftedChange > 0:
		return fmt.Sprintf("drift increased by %d resource(s)", d.DriftedChange)
	default:
		return "no change in drift"
	}
}

// Compare calculates the Delta between two snapshot entries.
func Compare(from, to Entry) Delta {
	d := Delta{
		From:          from,
		To:            to,
		DriftedChange: to.Result.Summary.TotalDrifted - from.Result.Summary.TotalDrifted,
		ManagedChange: to.Result.Summary.TotalResources - from.Result.Summary.TotalResources,
	}

	oldIndex := indexByID(from.Result.Drifted)
	newIndex := indexByID(to.Result.Drifted)

	for id, res := range newIndex {
		if _, existed := oldIndex[id]; !existed {
			d.NewlyDrifted = append(d.NewlyDrifted, res)
		}
	}
	for id, res := range oldIndex {
		if _, stillDrifted := newIndex[id]; !stillDrifted {
			d.Resolved = append(d.Resolved, res)
		}
	}
	return d
}

func indexByID(resources []model.Resource) map[string]model.Resource {
	m := make(map[string]model.Resource, len(resources))
	for _, r := range resources {
		m[r.ID] = r
	}
	return m
}
