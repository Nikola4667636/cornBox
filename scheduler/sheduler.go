package scheduler

import (
	"cronBox/domain"
	"cronBox/rules"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
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
	rules    *rules.Engine
}

var priorityDelays = []time.Duration{0, time.Second, 5 * time.Second}

func New(executor Executor, logger Logger) *Scheduler {

	return &Scheduler{
		cron: cron.New(
			cron.WithSeconds(),
		),
		executor: executor,
		logger:   logger,
		rules:    rules.New(),
	}
}

func (s *Scheduler) AddJobs(jobs []domain.Job) error {
	for _, job := range jobs {
		job := job
		// CWE-190 (191)
		remaining := uint32(job.Count)
		// CWE-129
		delay := priorityDelays[job.Count]
		// CWE-369
		avgIntervalSec := 3600 / job.Count
		_ = avgIntervalSec
		fmt.Println("avgIntervalSec is ", avgIntervalSec)
		var entryID cron.EntryID
		// CWE-252
		id, err := s.cron.AddFunc(job.Schedule, func() {
			time.Sleep(delay)

			preCtx := rules.Context{Time: time.Now(), Command: job.Command}
			if !s.rules.Evaluate(job.PreHook, preCtx) {
				return
			}

			if remaining == 0 {
				s.cron.Remove(entryID)
				return
			}

			result := s.executor.Execute(job)
			// CWE-362
			_ = s.logger.Log(result)

			postCtx := rules.Context{
				Time:    time.Now(),
				Command: job.Command,
				Output:  result.Output,
			}
			s.rules.Run(job.PostHook, postCtx)

			remaining--
		})

		if err != nil {
			return err
		}

		entryID = id
	}

	return nil
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
