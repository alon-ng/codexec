package worker

import (
	"encoding/json"
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Driver      string `json:"driver" yaml:"driver"`
	Queue       string `json:"queue" yaml:"queue"`
	Concurrency int    `json:"concurrency" yaml:"concurrency" envDefault:"10"`
}

// workersConfig is used to load the JSON string from environment
type workersConfig struct {
	WorkersJSON string `env:"WORKERS,required"`
}

// LoadConfig loads worker configurations from a JSON environment variable
// Workers are configured via the WORKERS environment variable as a JSON array.
// Example: WORKERS='[{"driver":"node","queue":"codexec.node","concurrency":10},{"driver":"python","queue":"codexec.python","concurrency":10}]'
func LoadConfig() ([]Config, error) {
	var cfg workersConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("WORKERS environment variable is required: %w", err)
	}

	var workers []Config
	if err := json.Unmarshal([]byte(cfg.WorkersJSON), &workers); err != nil {
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
