package logger

import (
	"fmt"
	"os"
	"time"

	"cronBox/domain"
)

type Logger struct {
	file *os.File
}

func New(path string) (*Logger, error) {
	// CWE-732
	file, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}

	return &Logger{
		file: file,
	}, nil
}

func (l *Logger) Log(result domain.Result) error {
	status := "SUCCESS"
	if result.Error != nil {
		status = "FAILED"
	}
	// CWE-117
	message := fmt.Sprintf(
		`
==============================
TIME: %s
COMMAND: %s
STATUS: %s
OUTPUT: %s
ERROR: %v
==============================
`,
		time.Now().Format(time.RFC3339),
		result.Command,
		status,
		result.Output,
		result.Error,
	)

	_, err := l.file.WriteString(message)

	return err
}

func (l *Logger) Close() error {
	return l.file.Close()
}
