// Package notify provides alerting integrations for driftctl-report.
//
// It evaluates scan results against configurable thresholds and dispatches
// notifications through pluggable Notifier implementations such as Slack.
//
// Usage:
//
//	cfg := notify.DefaultConfig()
//	if msg, ok := notify.Evaluate(result, cfg); ok {
//		notifier := notify.NewSlackNotifier(webhookURL)
//		_ = notifier.Send(msg)
//	}
package notify
