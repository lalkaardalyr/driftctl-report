// Package export provides writers that serialise drift scan results into
// various output formats for downstream consumption.
//
// Available writers:
//
//   - CSVWriter    – comma-separated values (.csv)
//   - JSONWriter   – structured JSON (.json)
//   - MarkdownWriter – GitHub-flavoured Markdown (.md)
//   - XLSXWriter   – tab-separated values suitable for spreadsheet import
//
// Each writer accepts an io.Writer so callers can direct output to a file,
// an HTTP response, or an in-memory buffer for testing.
package export
