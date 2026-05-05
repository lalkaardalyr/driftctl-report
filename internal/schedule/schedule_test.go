package schedule_test

import (
	"context"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/org/driftctl-report/internal/schedule"
)

func newLogger() *log.Logger {
	return log.New(os.Stdout, "", 0)
}

func TestScheduler_RegisterAndRun_ExecutesTask(t *testing.T) {
	var count int32

	job := &schedule.Job{
		Name:     "test-job",
		Interval: 50 * time.Millisecond,
		Task: func(ctx context.Context) error {
			atomic.AddInt32(&count, 1)
			return nil
		},
	}

	s := schedule.New(newLogger())
	s.Register(job)

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()

	go s.Run(ctx)
	<-ctx.Done()

	// Expect at least 2 executions (immediate + at least one tick).
	if got := atomic.LoadInt32(&count); got < 2 {
		t.Errorf("expected at least 2 executions, got %d", got)
	}
}

func TestScheduler_StopsOnContextCancel(t *testing.T) {
	var count int32

	job := &schedule.Job{
		Name:     "cancel-job",
		Interval: 10 * time.Millisecond,
		Task: func(ctx context.Context) error {
			atomic.AddInt32(&count, 1)
			return nil
		},
	}

	s := schedule.New(newLogger())
	s.Register(job)

	ctx, cancel := context.WithCancel(context.Background())
	go s.Run(ctx)

	time.Sleep(35 * time.Millisecond)
	cancel()
	time.Sleep(20 * time.Millisecond)

	snapshot := atomic.LoadInt32(&count)
	time.Sleep(30 * time.Millisecond)

	if after := atomic.LoadInt32(&count); after != snapshot {
		t.Errorf("job ran after context cancel: before=%d after=%d", snapshot, after)
	}
}

func TestScheduler_MultipleJobs(t *testing.T) {
	var a, b int32

	s := schedule.New(newLogger())
	s.Register(&schedule.Job{
		Name:     "job-a",
		Interval: 50 * time.Millisecond,
		Task: func(_ context.Context) error { atomic.AddInt32(&a, 1); return nil },
	})
	s.Register(&schedule.Job{
		Name:     "job-b",
		Interval: 50 * time.Millisecond,
		Task: func(_ context.Context) error { atomic.AddInt32(&b, 1); return nil },
	})

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	go s.Run(ctx)
	<-ctx.Done()

	if atomic.LoadInt32(&a) < 1 || atomic.LoadInt32(&b) < 1 {
		t.Error("expected both jobs to have executed at least once")
	}
}
