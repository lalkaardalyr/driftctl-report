// Package schedule provides a lightweight job scheduler for automating
// periodic drift scan report generation.
//
// A Scheduler manages one or more named Jobs, each associated with a
// time.Duration interval and a Task function. Jobs are started concurrently
// and execute immediately upon registration, then repeat on the configured
// interval until the context is cancelled.
//
// Example usage:
//
//	s := schedule.New(log.Default())
//	s.Register(&schedule.Job{
//		Name:     "nightly-drift",
//		Interval: 24 * time.Hour,
//		Task:     runDriftReport,
//	})
//	s.Run(ctx)
package schedule
