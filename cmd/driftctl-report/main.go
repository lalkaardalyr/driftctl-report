package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/example/driftctl-report/internal/model"
	"github.com/example/driftctl-report/internal/parser"
	"github.com/example/driftctl-report/internal/render"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	inputFile := flag.String("input", "", "Path to driftctl JSON output file (required)")
	outputFile := flag.String("output", "drift-report.html", "Path to write the HTML report")
	flag.Parse()

	if *inputFile == "" {
		flag.Usage()
		return fmt.Errorf("--input flag is required")
	}

	analysis, err := parser.ParseFile(*inputFile)
	if err != nil {
		return fmt.Errorf("parsing input file: %w", err)
	}

	scanResult := model.FromAnalysis(analysis)

	renderer, err := render.New()
	if err != nil {
		return fmt.Errorf("creating renderer: %w", err)
	}

	out, err := os.Create(*outputFile)
	if err != nil {
		return fmt.Errorf("creating output file %q: %w", *outputFile, err)
	}
	defer out.Close()

	if err := renderer.Render(out, scanResult); err != nil {
		return fmt.Errorf("rendering report: %w", err)
	}

	fmt.Printf("Report written to %s\n", *outputFile)
	return nil
}
