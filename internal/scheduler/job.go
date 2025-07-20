package scheduler

import (
	"github.com/go-co-op/gocron/v2"
)

type Job interface {
	Name() string
	Definition() gocron.JobDefinition // e.g. DurationJob(...)
	Task() any                        // The function to run
	Params() []any                    // Params to pass to the task
}
