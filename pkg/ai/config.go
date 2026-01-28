package ai

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	APIKey  string `env:"OPENAI_API_KEY,required"`
	BaseURL string `env:"OPENAI_BASE_URL" envDefault:"https://api.openai.com/v1"`
	Model   string `env:"OPENAI_MODEL" envDefault:"gpt-4o-mini"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
