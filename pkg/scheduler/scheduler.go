package scheduler

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type Scheduler struct {
	jobs   []*Job
	wg     sync.WaitGroup
	logger *slog.Logger
	jdi    JobDataStorage
}

func New(logger *slog.Logger) *Scheduler {
	return &Scheduler{
		logger: logger,
	}
}

func (s *Scheduler) SetJobDataStorage(jdi JobDataStorage) {
	s.jdi = jdi
}

func (s *Scheduler) jobDataImport() {
	if s.jdi == nil {
		s.logger.Info("scheduler.jobDataImport: JobDataStorage is nil, skipped")
		return
	}
	jobData := s.jdi.importJobsData()
	if len(jobData) == 0 {
		return
	}
	jobMap := make(map[string]*Job)
	for _, j := range s.jobs {
		jobMap[j.name] = j
	}
	for _, jd := range jobData {
		j, ok := jobMap[jd.name]
		if !ok {
			continue
		}
		j.lastStarted = jd.lastStarted
	}
}

func (s *Scheduler) AddJob(j *Job) {
	j.logger = s.logger
	j.wg = &s.wg
	s.jobs = append(s.jobs, j)
}

// Run all jobs until shutdown
func (s *Scheduler) Run(ctx context.Context) {
	s.jobDataImport()
	tick := time.NewTicker(time.Second)
forLoop:
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("scheduler.Run: shutdown")
			break forLoop
		case <-tick.C:
			for _, j := range s.jobs {
				j.TryRun()
			}
		}
	}

	s.wg.Wait()
}

type JobData struct {
	name        string
	lastStarted time.Time
}

type JobDataStorage interface {
	importJobsData() []*JobData
	exportJobsData(jobs []*Job) bool
}
