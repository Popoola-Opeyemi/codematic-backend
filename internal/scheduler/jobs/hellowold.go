package jobs

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type HelloJob struct{}

func (h HelloJob) Name() string {
	return "HelloJob"
}

func (h HelloJob) Definition() gocron.JobDefinition {
	return gocron.DurationJob(30 * time.Minute)
}

func (h HelloJob) Task() any {
	return func(name string) {
		fmt.Println("Hello,", name)
	}
}

func (h HelloJob) Params() []any {
	return []any{"world"}
}
