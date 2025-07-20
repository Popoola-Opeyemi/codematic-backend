
## ⏰ Background Job Scheduler

robust background job scheduler built with [gocron](https://github.com/go-co-op/gocron) that handles periodic tasks, data synchronization, and system maintenance.

### Scheduler Features

- **UTC Timezone**: All jobs run in UTC for consistency
- **Memory Management**: Automatic memory cleanup every hour
- **Graceful Shutdown**: Proper cleanup on application termination
- **Job Monitoring**: Comprehensive logging for all job executions
- **Flexible Scheduling**: Support for cron expressions and duration-based jobs

### Current Jobs

#### 1. **HelloJob** (Development/Testing)
- **Schedule**: Every 10 seconds
- **Purpose**: Simple test job for development
- **Location**: `internal/scheduler/jobs/hellowold.go`

#### 2. **JupiterJob** (Production)
- **Schedule**: Every Monday at 2:00 AM UTC
- **Purpose**: Fetches and syncs Jupiter token data
- **Features**:
  - Fetches all tokens from Jupiter API
  - Filters out existing tokens to avoid duplicates
  - Batch inserts new tokens into database
  - Memory cleanup after execution
  - 5-minute timeout with context cancellation
- **Location**: `internal/scheduler/jobs/job.jupiter.go`

#### 3. **Memory Cleanup Job** (System)
- **Schedule**: Every hour
- **Purpose**: System maintenance and memory optimization
- **Features**:
  - Forces garbage collection
  - Logs memory statistics
  - Monitors goroutine count
  - Tracks memory allocation metrics

### Creating New Jobs

To add a new background job:

1. **Create a new job file** in `internal/scheduler/jobs/`:

```go
package jobs

import (
    "time"
    "github.com/go-co-op/gocron/v2"
)

type MyNewJob struct {
    // Add any dependencies here
}

func (j MyNewJob) Name() string {
    return "MyNewJob"
}

func (j MyNewJob) Definition() gocron.JobDefinition {
    // Run every day at 3 AM UTC
    return gocron.CronJob("0 3 * * *", true)
    
    // Or run every 30 minutes
    // return gocron.DurationJob(30 * time.Minute)
}

func (j MyNewJob) Task() any {
    return func() {
        // Your job logic here
        // Remember to handle errors and add logging
    }
}

func (j MyNewJob) Params() []any {
    return []any{} // Add parameters if needed
}
```

2. **Register the job** in `cmd/main.go`:

```go
// Register background jobs
jobList := []scheduler.Job{
    jobs.HelloJob{},
    jobs.NewJupiterJob(portfolioService),
    jobs.MyNewJob{}, // Add your new job here
}
```

### Job Best Practices

1. **Error Handling**: Always handle errors gracefully and log them
2. **Timeouts**: Use context with timeouts for long-running operations
3. **Memory Management**: Clean up large objects and force GC when needed
4. **Logging**: Use structured logging with relevant context
5. **Dependencies**: Inject dependencies through the job struct
6. **Idempotency**: Make jobs safe to run multiple times

### Monitoring Jobs

Jobs are automatically logged with:
- Job name and execution time
- Success/failure status
- Memory usage statistics (for cleanup job)
- Error details when failures occur

### Scheduler Configuration

The scheduler is configured in `cmd/main.go` and includes:
- UTC timezone for consistency
- Graceful shutdown handling
- Automatic memory cleanup
- Comprehensive logging