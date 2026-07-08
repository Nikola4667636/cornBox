package executor

import (
	"os/exec"
	"time"

	"cronBox/domain"
)

type ShellExecutor struct{}

func New() *ShellExecutor {
	return &ShellExecutor{}
}

func (e *ShellExecutor) Execute(job domain.Job) domain.Result {
	start := time.Now()

	cmd := exec.Command(
		"sh",
		"-c",
		job.Command,
	)

	output, err := cmd.CombinedOutput()
	end := time.Now()

	return domain.Result{
		Command:    job.Command,
		StartedAt:  start,
		FinishedAt: end,
		Output:     string(output),
		Error:      err,
	}
}
