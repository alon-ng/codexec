package config

import (
	"codim/internal/rabbitmq"
	"codim/internal/utils/env"
	"codim/internal/utils/logger"
	"codim/internal/worker"
	"sync"
	"time"
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
	CmdPrefix        string
	ExecutionTimeout time.Duration
	ShutdownTimeout  time.Duration
}

func Load() (Config, error) {
	once.Do(func() {
		config.Logger = logger.Load()
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

		// Load CmdPrefix (optional)
		config.CmdPrefix = env.Get("CMD_PREFIX", "")

		// Load timeouts with defaults
		config.ExecutionTimeout = env.GetDuration("EXECUTION_TIMEOUT", 10*time.Second)
		config.ShutdownTimeout = env.GetDuration("SHUTDOWN_TIMEOUT", 30*time.Second)
	})

	return config, loadErr
}
