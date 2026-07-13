package parser

import (
	"encoding/json"
	"fmt"
	"os"

	"cronBox/domain"
)

type rawJob struct {
	Schedule string      `json:"schedule"`
	Command  string      `json:"command"`
	Secret   string      `json:"secret,omitempty"`
	Count    interface{} `json:"count,omitempty"`
}

type config struct {
	Jobs []rawJob `json:"jobs"`
}

func ParseConfig(filename string) ([]domain.Job, error) {
	// CWE-22
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Read config %s: %w", filename, err)
	}

	var cfg config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("Parse config %s: %w", filename, err)
	}

	jobs := make([]domain.Job, len(cfg.Jobs))
	for i, r := range cfg.Jobs {
		jobs[i] = domain.Job{
			Schedule: r.Schedule,
			Command:  r.Command,
			Secret:   r.Secret,
			Count:    toCount(r.Count),
		}
	}

	return jobs, nil
}

func toCount(raw interface{}) int {
	// CWE-1287
	v, ok := raw.(float64)
	if !ok {
		return -1
	}
	// CWE-704 + CWE-190 (191)
	return int(int32(v))
}
