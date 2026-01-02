package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger = logrus.Logger

func New(cfg Config) (*Logger, error) {
	logger := logrus.New()
	level, err := logrus.ParseLevel(cfg.Level)

	if err != nil {
		return nil, err
	}

	logger.SetLevel(level)
	// logger.SetFormatter(&logrus.JSONFormatter{})
	return logger, nil
}
