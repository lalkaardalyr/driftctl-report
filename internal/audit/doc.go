// Package audit provides structured, append-only audit logging for
// driftctl-report operations.
//
// Each audit event is written as a single line of JSON (JSONL format) to any
// io.Writer, making the log easy to ingest into log-aggregation systems such
// as Loki, Splunk, or CloudWatch Logs.
//
// # Basic usage
//
//	logger := audit.New(os.Stdout, "ci")
//	logger.ScanStarted("drift.json")
//	logger.ScanCompleted("drift.json", "report.html", 4, 12)
//
// # Rotating file
//
//	rf, _ := audit.OpenRotating("/var/log/driftctl-report/audit")
//	defer rf.Close()
//	logger := audit.New(rf, "scheduler")
package audit
