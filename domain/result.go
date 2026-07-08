package domain

import "time"

type Result struct {
	Command    string
	StartedAt  time.Time
	FinishedAt time.Time
	Output     string
	Error      error
}
