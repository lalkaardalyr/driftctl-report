package filter_test

import (
	"testing"

	"github.com/snyk/driftctl-report/internal/filter"
	"github.com/snyk/driftctl-report/internal/model"
)

func makeResources() []model.Resource {
	return []model.Resource{
		{ID: "bucket-1", Type: "aws_s3_bucket", Source: "terraform"},
		{ID: "bucket-2", Type: "aws_s3_bucket", Source: "remote"},
		{ID: "sg-1", Type: "aws_security_group", Source: "terraform"},
		{ID: "sg-2", Type: "aws_security_group", Source: "remote"},
		{ID: "role-1", Type: "aws_iam_role", Source: "terraform"},
	}
}

func TestApply_ByType(t *testing.T) {
	res := filter.Apply(makeResources(), filter.ByType("aws_s3_bucket"))
	if len(res) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(res))
	}
}

func TestApply_BySource(t *testing.T) {
	res := filter.Apply(makeResources(), filter.BySource("terraform"))
	if len(res) != 3 {
		t.Fatalf("expected 3 resources, got %d", len(res))
	}
}

func TestApply_MultipleFilters(t *testing.T) {
	res := filter.Apply(makeResources(), filter.ByType("aws_s3_bucket"), filter.BySource("terraform"))
	if len(res) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(res))
	}
	if res[0].ID != "bucket-1" {
		t.Errorf("expected bucket-1, got %s", res[0].ID)
	}
}

func TestApply_NoFilters(t *testing.T) {
	res := filter.Apply(makeResources())
	if len(res) != 5 {
		t.Fatalf("expected 5 resources, got %d", len(res))
	}
}

func TestGroupByType(t *testing.T) {
	groups := filter.GroupByType(makeResources())
	if len(groups["aws_s3_bucket"]) != 2 {
		t.Errorf("expected 2 s3 buckets, got %d", len(groups["aws_s3_bucket"]))
	}
	if len(groups["aws_security_group"]) != 2 {
		t.Errorf("expected 2 security groups, got %d", len(groups["aws_security_group"]))
	}
	if len(groups["aws_iam_role"]) != 1 {
		t.Errorf("expected 1 iam role, got %d", len(groups["aws_iam_role"]))
	}
}

func TestGroupBySource(t *testing.T) {
	groups := filter.GroupBySource(makeResources())
	if len(groups["terraform"]) != 3 {
		t.Errorf("expected 3 terraform resources, got %d", len(groups["terraform"]))
	}
	if len(groups["remote"]) != 2 {
		t.Errorf("expected 2 remote resources, got %d", len(groups["remote"]))
	}
}
