package scheduler

import (
	"context"
	"runtime"
	"time"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

type Scheduler struct {
	logger *zap.Logger
	s      gocron.Scheduler // Use a pointer to access scheduling methods
}

// New initializes the gocron scheduler and wraps it
func New(logger *zap.Logger) (*Scheduler, error) {
	s, err := gocron.NewScheduler(gocron.WithLocation(time.UTC))
	if err != nil {
		return nil, err
	}
	return &Scheduler{
		logger: logger,
		s:      s, // explicitly assign the pointer
	}, nil
}

// RegisterJobs takes all jobs and registers them with gocron
func (sc *Scheduler) RegisterJobs(ctx context.Context, jobs []Job) error {
	for _, job := range jobs {
		task := gocron.NewTask(job.Task(), job.Params()...)

		_, err := sc.s.NewJob(job.Definition(), task)
		if err != nil {
			return err
		}

		sc.logger.Info("Registered job", zap.String("job", job.Name()))
	}

	// Add a memory cleanup job that runs every hour
	_, err := sc.s.NewJob(
		gocron.DurationJob(1*time.Hour),
		gocron.NewTask(sc.cleanupMemory),
	)
	if err != nil {
		return err
	}

	return nil
}

// Start the scheduler
func (sc *Scheduler) Start() {
	sc.logger.Info("Starting scheduler")
	sc.s.Start()
}

// Stop the scheduler gracefully
func (sc *Scheduler) Stop() {
	sc.logger.Info("Stopping scheduler")
	sc.s.Shutdown()
}

// cleanupMemory performs memory cleanup
func (sc *Scheduler) cleanupMemory() {
	sc.logger.Info("Running memory cleanup")

	// Force garbage collection
	runtime.GC()

	// Read memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	sc.logger.Info("Memory stats after cleanup",
		zap.Uint64("alloc", m.Alloc),
		zap.Uint64("total_alloc", m.TotalAlloc),
		zap.Uint64("sys", m.Sys),
		zap.Uint32("num_gc", m.NumGC),
		zap.Int("goroutines", runtime.NumGoroutine()),
	)
}
