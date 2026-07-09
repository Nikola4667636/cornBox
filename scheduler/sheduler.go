package scheduler

import (
	"github.com/robfig/cron/v3"

	"cronBox/domain"
)

type Executor interface {
	Execute(job domain.Job) domain.Result
}

type Logger interface {
	Log(result domain.Result) error
}

type Scheduler struct {
	cron     *cron.Cron
	executor Executor
	logger   Logger
}

func New(executor Executor, logger Logger) *Scheduler {

	return &Scheduler{
		cron: cron.New(
			cron.WithSeconds(),
		),
		executor: executor,
		logger:   logger,
	}
}

// CWE-252
func (s *Scheduler) AddJobs(jobs []domain.Job) error {
	for _, job := range jobs {
		job := job
		// CWE-400
		_, err := s.cron.AddFunc(
			job.Schedule,
			func() {
				result := s.executor.Execute(job)
				_ = s.logger.Log(result)
			},
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
