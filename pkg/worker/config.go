package worker

import (
	"encoding/json"
	"fmt"

	"github.com/caarlos0/env/v11"
	v "github.com/go-playground/validator/v10"
)

type Config struct {
	Driver       string `json:"driver" validate:"required"`
	Queue        string `json:"queue" validate:"required"`
	ResultsQueue string `json:"results_queue" validate:"required"`
	Concurrency  int    `json:"concurrency"  envDefault:"10"`
}

// workersConfig is used to load the JSON string from environment
type workersConfig struct {
	WorkersJSON string `env:"WORKERS,required"`
}

// LoadConfig loads worker configurations from a JSON environment variable
// Workers are configured via the WORKERS environment variable as a JSON array.
// Example: WORKERS='[{"driver":"node","queue":"codexec.node","concurrency":10},{"driver":"python","queue":"codexec.python","concurrency":10}]'
func LoadConfig() ([]Config, error) {
	validate := v.New()
	var cfg workersConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("WORKERS environment variable is required: %w", err)
	}

	var workers []Config
	if err := json.Unmarshal([]byte(cfg.WorkersJSON), &workers); err != nil {
		return nil, fmt.Errorf("failed to parse WORKERS JSON: %w", err)
	}

	for _, worker := range workers {
		if err := validate.Struct(worker); err != nil {
			return nil, fmt.Errorf("invalid worker configuration: %w", err)
		}
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
