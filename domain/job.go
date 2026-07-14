package domain

type Job struct {
	Schedule string `json:"schedule"`
	Command  string `json:"command"`
	Secret   string `json:"secret"`
	Count    int    `json:"count"`
	PreHook  string `json:"preHook"`
	PostHook string `json:"postHook"`
}
