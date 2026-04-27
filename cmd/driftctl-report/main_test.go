package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun_MissingInputFlag(t *testing.T) {
	// Simulate missing --input by calling run directly with no flag set.
	// We rely on the flag package default being empty string.
	os.Args = []string{"driftctl-report"}
	err := run()
	if err == nil {
		t.Fatal("expected error when --input is missing, got nil")
	}
}

func TestRun_InvalidInputFile(t *testing.T) {
	os.Args = []string{"driftctl-report", "--input", "/nonexistent/path.json"}
	err := run()
	if err == nil {
		t.Fatal("expected error for missing input file, got nil")
	}
}

func TestRun_ValidInput(t *testing.T) {
	const sampleJSON = `{
		"summary": {"total_resources": 3, "total_managed": 2, "total_unmanaged": 1, "total_missing": 0, "total_changed": 1, "coverage": 66},
		"managed": [],
		"unmanaged": [],
		"missing": [],
		"differences": []
	}`

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "drift.json")
	outputPath := filepath.Join(tmpDir, "report.html")

	if err := os.WriteFile(inputPath, []byte(sampleJSON), 0o644); err != nil {
		t.Fatalf("failed to write temp input file: %v", err)
	}

	os.Args = []string{"driftctl-report", "--input", inputPath, "--output", outputPath}
	if err := run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output file is empty")
	}
}
