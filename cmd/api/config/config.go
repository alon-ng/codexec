package config

import (
	"codim/internal/api/v1"
	"codim/internal/db"
	"codim/internal/redis"
	"codim/internal/utils/logger"
	"sync"

	"github.com/caarlos0/env/v11"
)

var (
	once    sync.Once
	config  Config
	loadErr error
)

type Config struct {
	Logger logger.Config
	API    api.Config
	DB     db.Config
	Redis  redis.Config
}

func Load() (Config, error) {
	once.Do(func() {
		loggerCfg, err := logger.LoadConfig()
		if err != nil {
			loadErr = err
			return
		}
		config.Logger = loggerCfg

		apiCfg, err := api.LoadConfig()
		if err != nil {
			loadErr = err
			return
		}
		config.API = apiCfg

		dbCfg, err := db.LoadConfig()
		if err != nil {
			loadErr = err
			return
		}
		config.DB = dbCfg

		redisCfg, err := redis.LoadConfig()
		if err != nil {
			loadErr = err
			return
		}
		config.Redis = redisCfg

		// Parse the remaining fields using caarlos0/env
		if err := env.Parse(&config); err != nil {
			loadErr = err
			return
		}
	})

	return config, loadErr
}
