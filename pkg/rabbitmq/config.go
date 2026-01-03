package rabbitmq

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	URL string `env:"RABBITMQ_URL,required"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
