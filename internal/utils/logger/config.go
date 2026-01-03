package logger

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Level string `env:"LOGGER_LEVEL" envDefault:"info"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
