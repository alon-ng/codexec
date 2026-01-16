package config

import (
	"codim/pkg/api/v1"
	"codim/pkg/db"
	"sync"

	"github.com/caarlos0/env/v11"
)

var (
	once    sync.Once
	config  Config
	loadErr error
)

type Config struct {
	API api.Config
	DB  db.Config
}

func Load() (Config, error) {
	once.Do(func() {
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

		// Parse the remaining fields using caarlos0/env
		if err := env.Parse(&config); err != nil {
			loadErr = err
			return
		}
	})

	return config, loadErr
}
