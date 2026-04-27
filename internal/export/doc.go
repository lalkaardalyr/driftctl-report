// Package export provides writers that serialise a model.ScanResult into
// various output formats.
//
// Supported formats:
//
//	- CSV  via NewCSVWriter
//	- JSON via NewJSONWriter
//
// Each writer accepts an io.Writer so callers can target files, buffers, or
// standard output without coupling to a specific destination.
package export
