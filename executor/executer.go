package executor

import (
	"os"
	"os/exec"
	"time"

	"cronBox/domain"
	"cronBox/secretstore"
)

type ShellExecutor struct{}

func New() *ShellExecutor {
	return &ShellExecutor{}
}

func (e *ShellExecutor) Execute(job domain.Job) domain.Result {
	start := time.Now()
	// CWE-78
	cmd := exec.Command("sh", "-c", job.Command)

	cmd.Env = os.Environ()

	if job.Secret != "" {
		secret, err := secretstore.Decrypt(job.Secret)
		if err == nil {
			cmd.Env = append(cmd.Env, "CRONBOX_SECRET="+secret)
		}
	}

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
