package parser

import (
	"encoding/json"
	"fmt"
	"os"

	"cronBox/domain"
)

type config struct {
	Jobs []domain.Job `json:"jobs"`
}

func ParseConfig(filename string) ([]domain.Job, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Read config %s: %w", filename, err)
	}

	var cfg config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("Parse config %s: %w", filename, err)
	}

	return cfg.Jobs, nil
}
