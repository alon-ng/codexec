package redis

import "github.com/caarlos0/env/v11"

type Config struct {
	Host       string `env:"REDIS_HOST,required"`
	Port       int    `env:"REDIS_PORT" envDefault:"6379"`
	Password   string `env:"REDIS_PASSWORD"`
	DB         int    `env:"REDIS_DB" envDefault:"0"`
	MaxRetries int    `env:"REDIS_MAX_RETRIES" envDefault:"3"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
