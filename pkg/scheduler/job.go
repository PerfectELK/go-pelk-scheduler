package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type Job struct {
	name        string
	fn          func(ctx context.Context, arg ...any)
	args        []any
	interval    time.Duration
	lastStarted time.Time
	logger      *slog.Logger
	wg          *sync.WaitGroup
	ran         bool
}

func NewJob(
	name string,
	fn func(ctx context.Context, args ...any),
	interval time.Duration,
	args ...any,
) *Job {
	return &Job{
		name:     name,
		fn:       fn,
		args:     args,
		interval: interval,
	}
}

func (j *Job) TryRun() {
	if j.ran {
		return
	}
	if time.Now().Before(j.lastStarted.Add(j.interval)) {
		return
	}
	j.run()
}

func (j *Job) run() {
	j.ran = true
	j.logger.Info(fmt.Sprintf("job.run: start %s", j.name))
	j.wg.Add(1)
	go func() {
		j.lastStarted = time.Now()
		defer func() {
			j.ran = false
			j.wg.Done()
		}()
		ctx := context.Background()
		j.fn(ctx, j.args)
		j.logger.Info(fmt.Sprintf("job.run: end %s", j.name))
	}()
}

func (j *Job) IsRan() bool {
	return j.ran
}

func (j *Job) GetName() string {
	return j.name
}

func (j *Job) GetLatStarted() time.Time {
	return j.lastStarted
}
