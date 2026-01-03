package main

import (
	"codim/cmd/codexec/config"
	"codim/pkg/executors"
	"codim/pkg/executors/drivers"
	"codim/pkg/rabbitmq"
	"codim/pkg/utils/logger"
	"codim/pkg/worker"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	logger, err := initializeLogger(cfg)
	if err != nil {
		logrus.Fatalf("Failed to initialize logger: %v", err)
	}

	rmqClient, err := initializeRabbitMQ(cfg, logger)
	if err != nil {
		logrus.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rmqClient.Close()

	sigChan := setupSignalHandling()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-sigChan
		logger.Info("Received shutdown signal, shutting down...")
		cancel()
	}()

	for _, wCfg := range cfg.Workers {
		go func() {
			driver, err := drivers.New(wCfg.Driver, cfg.CmdPrefix, logger)
			if err != nil {
				logger.Errorf("Failed to initialize driver %s: %v", wCfg.Driver, err)
				return
			}

			executorService := executors.New(driver, logger, cfg.ExecutionTimeout)

			w := worker.New(rmqClient, executorService, logger, wCfg)
			logger.Infof("Starting worker for queue %s (driver: %s)", wCfg.Queue, wCfg.Driver)

			// Start will handle reconnection internally
			if err := w.Start(ctx); err != nil {
				logger.Errorf("Worker for queue %s stopped with error: %v", wCfg.Queue, err)
			}

			// Ensure worker is stopped
			w.Stop()
			logger.Infof("Worker for queue %s stopped", wCfg.Queue)
		}()
	}

	<-ctx.Done()
	logger.Info("Application stopped")

	os.Exit(0)
}

// initializeLogger creates and initializes the logger from configuration
func initializeLogger(cfg config.Config) (*logger.Logger, error) {
	log, err := logger.New(cfg.Logger)
	if err != nil {
		return nil, err
	}

	log.Info("Logger initialized successfully")
	return log, nil
}

// initializeRabbitMQ creates a new RabbitMQ client from configuration
func initializeRabbitMQ(cfg config.Config, log *logger.Logger) (*rabbitmq.Client, error) {
	rmqClient, err := rabbitmq.NewClient(cfg.RabbitMQ, log)
	if err != nil {
		return nil, err
	}

	return rmqClient, nil
}

// setupSignalHandling configures signal handling for graceful shutdown
func setupSignalHandling() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	return sigChan
}
