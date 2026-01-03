package config

import (
	"codim/internal/rabbitmq"
	"codim/internal/utils/logger"
	"codim/internal/worker"
	"sync"
	"time"

	"github.com/caarlos0/env/v11"
)

var (
	once    sync.Once
	config  Config
	loadErr error
)

type Config struct {
	Logger           logger.Config
	RabbitMQ         rabbitmq.Config
	Workers          []worker.Config
	CmdPrefix        string        `env:"CMD_PREFIX"`
	ExecutionTimeout time.Duration `env:"EXECUTION_TIMEOUT" envDefault:"10s"`
	ShutdownTimeout  time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"30s"`
}

func Load() (Config, error) {
	once.Do(func() {
		loggerCfg, err := logger.Load()
		if err != nil {
			loadErr = err
			return
		}
		config.Logger = loggerCfg

		rmqConfig, err := rabbitmq.Load()
		if err != nil {
			loadErr = err
			return
		}
		config.RabbitMQ = rmqConfig

		workers, err := worker.Load()
		if err != nil {
			loadErr = err
			return
		}
		config.Workers = workers

		// Parse the remaining fields using caarlos0/env
		if err := env.Parse(&config); err != nil {
			loadErr = err
			return
		}
	})

	return config, loadErr
}
