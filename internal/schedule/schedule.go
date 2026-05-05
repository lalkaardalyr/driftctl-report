// Package schedule provides cron-style scheduling for periodic drift scans
// and report generation.
package schedule

import (
	"context"
	"log"
	"time"
)

// Job represents a scheduled task that runs a drift scan pipeline.
type Job struct {
	Name     string
	Interval time.Duration
	Task     func(ctx context.Context) error
}

// Scheduler manages and runs a collection of periodic jobs.
type Scheduler struct {
	jobs   []*Job
	logger *log.Logger
}

// New creates a Scheduler with the provided logger.
func New(logger *log.Logger) *Scheduler {
	return &Scheduler{logger: logger}
}

// Register adds a Job to the scheduler.
func (s *Scheduler) Register(job *Job) {
	s.jobs = append(s.jobs, job)
}

// Run starts all registered jobs and blocks until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	for _, job := range s.jobs {
		go s.runJob(ctx, job)
	}
	<-ctx.Done()
	s.logger.Println("scheduler: context cancelled, stopping all jobs")
}

func (s *Scheduler) runJob(ctx context.Context, job *Job) {
	s.logger.Printf("scheduler: starting job %q with interval %s", job.Name, job.Interval)
	ticker := time.NewTicker(job.Interval)
	defer ticker.Stop()

	// Run immediately on start.
	s.execute(ctx, job)

	for {
		select {
		case <-ticker.C:
			s.execute(ctx, job)
		case <-ctx.Done():
			s.logger.Printf("scheduler: stopping job %q", job.Name)
			return
		}
	}
}

func (s *Scheduler) execute(ctx context.Context, job *Job) {
	s.logger.Printf("scheduler: executing job %q", job.Name)
	if err := job.Task(ctx); err != nil {
		s.logger.Printf("scheduler: job %q failed: %v", job.Name, err)
	}
}
