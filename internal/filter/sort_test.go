package filter_test

import (
	"testing"

	"github.com/snyk/driftctl-report/internal/filter"
	"github.com/snyk/driftctl-report/internal/model"
)

func TestSort_ByID(t *testing.T) {
	input := []model.Resource{
		{ID: "z-resource", Type: "aws_s3_bucket"},
		{ID: "a-resource", Type: "aws_iam_role"},
		{ID: "m-resource", Type: "aws_security_group"},
	}
	sorted := filter.Sort(input, filter.SortByID)
	if sorted[0].ID != "a-resource" || sorted[1].ID != "m-resource" || sorted[2].ID != "z-resource" {
		t.Errorf("unexpected sort order: %v", sorted)
	}
}

func TestSort_ByType(t *testing.T) {
	input := []model.Resource{
		{ID: "b", Type: "aws_s3_bucket"},
		{ID: "a", Type: "aws_iam_role"},
		{ID: "c", Type: "aws_iam_role"},
	}
	sorted := filter.Sort(input, filter.SortByType)
	if sorted[0].Type != "aws_iam_role" {
		t.Errorf("expected aws_iam_role first, got %s", sorted[0].Type)
	}
	// secondary sort by ID within same type
	if sorted[0].ID != "a" || sorted[1].ID != "c" {
		t.Errorf("expected secondary sort by ID: got %s, %s", sorted[0].ID, sorted[1].ID)
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	input := []model.Resource{
		{ID: "z", Type: "aws_s3_bucket"},
		{ID: "a", Type: "aws_iam_role"},
	}
	_ = filter.Sort(input, filter.SortByID)
	if input[0].ID != "z" {
		t.Error("original slice was mutated")
	}
}

func TestUniqueTypes(t *testing.T) {
	input := []model.Resource{
		{ID: "1", Type: "aws_s3_bucket"},
		{ID: "2", Type: "aws_s3_bucket"},
		{ID: "3", Type: "aws_iam_role"},
		{ID: "4", Type: "aws_security_group"},
	}
	types := filter.UniqueTypes(input)
	if len(types) != 3 {
		t.Fatalf("expected 3 unique types, got %d", len(types))
	}
	if types[0] != "aws_iam_role" || types[1] != "aws_s3_bucket" || types[2] != "aws_security_group" {
		t.Errorf("unexpected types order: %v", types)
	}
}
