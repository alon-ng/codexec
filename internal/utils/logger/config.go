package logger

import "codim/internal/utils/env"

type Config struct {
	Level string `yaml:"level" validate:"omitempty,oneof=debug info warn error"`
}

func Load() Config {
	return Config{
		Level: env.Get("LOGGER_LEVEL", "info"),
	}
}
