package rabbitmq

import (
	"codim/internal/utils/env"
	"fmt"
)

type Config struct {
	URL string `yaml:"url"`
}

func Load() (Config, error) {
	url := env.Get("RABBITMQ_URL", "")
	if url == "" {
		return Config{}, fmt.Errorf("RABBITMQ_URL environment variable is required")
	}

	return Config{
		URL: url,
	}, nil
}
