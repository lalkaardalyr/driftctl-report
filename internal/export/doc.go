// Package export provides writers that serialise drift scan results into
// machine-readable formats suitable for downstream processing or archiving.
//
// Currently supported formats:
//
//   - CSV  — via CSVWriter, one row per resource with status and source columns.
//
// Example usage:
//
//	w := export.NewCSVWriter(os.Stdout)
//	if err := w.Write(scanResult); err != nil {
//		log.Fatal(err)
//	}
package export
