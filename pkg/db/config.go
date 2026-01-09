package db

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ConnectionString string        `env:"DB_CONNECTION_STRING,required"`
	MaxConns         int           `env:"DB_MAX_CONNS" envDefault:"10"`
	MinConns         int           `env:"DB_MIN_CONNS" envDefault:"2"`
	MaxIdleConns     int           `env:"DB_MAX_IDLE_CONNS" envDefault:"10"`
	MaxOpenConns     int           `env:"DB_MAX_OPEN_CONNS" envDefault:"10"`
	ConnMaxLifetime  time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"10s"`
	ConnMaxIdleTime  time.Duration `env:"DB_CONN_MAX_IDLE_TIME" envDefault:"10s"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
