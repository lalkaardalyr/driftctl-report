package diff_test

import (
	"testing"

	"github.com/owner/driftctl-report/internal/diff"
	"github.com/owner/driftctl-report/internal/model"
)

func makeResource(rtype, id string) model.Resource {
	return model.Resource{Type: rtype, ID: id}
}

func makeResult(drifted, unmanaged []model.Resource) model.ScanResult {
	return model.ScanResult{
		DriftedResources:   drifted,
		UnmanagedResources: unmanaged,
	}
}

func TestCompare_NewlyDrifted(t *testing.T) {
	baseline := makeResult(nil, nil)
	current := makeResult([]model.Resource{makeResource("aws_s3_bucket", "my-bucket")}, nil)

	r := diff.Compare(baseline, current)

	if len(r.NewlyDrifted) != 1 {
		t.Fatalf("expected 1 newly drifted, got %d", len(r.NewlyDrifted))
	}
	if r.NewlyDrifted[0].ID != "my-bucket" {
		t.Errorf("unexpected resource id: %s", r.NewlyDrifted[0].ID)
	}
}

func TestCompare_Resolved(t *testing.T) {
	baseline := makeResult([]model.Resource{makeResource("aws_s3_bucket", "old-bucket")}, nil)
	current := makeResult(nil, nil)

	r := diff.Compare(baseline, current)

	if len(r.Resolved) != 1 {
		t.Fatalf("expected 1 resolved, got %d", len(r.Resolved))
	}
	if r.Resolved[0].ID != "old-bucket" {
		t.Errorf("unexpected resource id: %s", r.Resolved[0].ID)
	}
}

func TestCompare_StillDrifted(t *testing.T) {
	res := makeResource("aws_instance", "i-123")
	baseline := makeResult([]model.Resource{res}, nil)
	current := makeResult([]model.Resource{res}, nil)

	r := diff.Compare(baseline, current)

	if len(r.StillDrifted) != 1 {
		t.Fatalf("expected 1 still drifted, got %d", len(r.StillDrifted))
	}
	if len(r.NewlyDrifted) != 0 {
		t.Errorf("expected 0 newly drifted, got %d", len(r.NewlyDrifted))
	}
}

func TestCompare_NewlyUnmanaged(t *testing.T) {
	baseline := makeResult(nil, nil)
	current := makeResult(nil, []model.Resource{makeResource("aws_iam_role", "role-x")})

	r := diff.Compare(baseline, current)

	if len(r.NewlyUnmanaged) != 1 {
		t.Fatalf("expected 1 newly unmanaged, got %d", len(r.NewlyUnmanaged))
	}
}

func TestCompare_EmptyBothSides(t *testing.T) {
	r := diff.Compare(makeResult(nil, nil), makeResult(nil, nil))

	if len(r.NewlyDrifted)+len(r.Resolved)+len(r.StillDrifted)+len(r.NewlyUnmanaged) != 0 {
		t.Error("expected all diff slices to be empty")
	}
}
