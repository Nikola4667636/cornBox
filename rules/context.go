package rules

import "time"

type Context struct {
	Time    time.Time
	Command string
	Output  string
}
