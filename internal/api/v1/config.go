package api

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port                int           `env:"API_PORT" envDefault:"8080"`
	PasswordSalt        string        `env:"PASSWORD_SALT,required"`
	JwtSecret           string        `env:"JWT_SECRET,required"`
	JwtTTL              time.Duration `env:"JWT_TTL" envDefault:"6h"`
	JwtRenewalThreshold time.Duration `env:"JWT_RENEWAL_THRESHOLD" envDefault:"15m"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
