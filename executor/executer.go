package executor

import (
	"cronBox/domain"
	"cronBox/secretstore"
	"fmt"
	"os"
	"os/exec"
	"time"
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
			// CWE-215
			fmt.Printf("[debug] secret for %q: %s\n", job.Command, secret)
			cmd.Env = append(cmd.Env, "CRONBOX_SECRET="+secret)
		}
	}

	output, err := cmd.CombinedOutput()
	// CWE-606
	for err != nil {
		output, err = cmd.CombinedOutput()
	}
	end := time.Now()

	return domain.Result{
		Command:    job.Command,
		StartedAt:  start,
		FinishedAt: end,
		Output:     string(output),
		Error:      err,
	}
}
