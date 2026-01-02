package worker

import (
	"codim/internal/utils/env"
	"encoding/json"
	"fmt"
)

type Config struct {
	Driver      string `json:"driver" yaml:"driver"`
	Queue       string `json:"queue" yaml:"queue"`
	Concurrency int    `json:"concurrency" yaml:"concurrency"`
}

// Load loads worker configurations from a JSON environment variable
// Workers are configured via the WORKERS environment variable as a JSON array.
// Example: WORKERS='[{"driver":"node","queue":"codexec.node","concurrency":10},{"driver":"python","queue":"codexec.python","concurrency":10}]'
func Load() ([]Config, error) {
	workersJSON := env.Get("WORKERS", "")
	if workersJSON == "" {
		return nil, fmt.Errorf("WORKERS environment variable is required")
	}

	var workers []Config
	if err := json.Unmarshal([]byte(workersJSON), &workers); err != nil {
		return nil, fmt.Errorf("failed to parse WORKERS JSON: %w", err)
	}

	if len(workers) == 0 {
		return nil, fmt.Errorf("at least one worker configuration is required in WORKERS")
	}

	// Set default concurrency if not specified
	for i := range workers {
		if workers[i].Concurrency == 0 {
			workers[i].Concurrency = 10
		}
	}

	return workers, nil
}
