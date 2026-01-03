package worker

import (
	"codim/pkg/executors"
	"codim/pkg/rabbitmq"
	"codim/pkg/utils/logger"
	"context"
	"fmt"
)

type Worker struct {
	concurrency     int
	queue           string
	ctx             context.Context
	cancel          context.CancelFunc
	rmqClient       *rabbitmq.Client
	executorService *executors.Service
	logger          *logger.Logger
}

func New(
	rmqClient *rabbitmq.Client,
	executorService *executors.Service,
	logger *logger.Logger,
	cfg Config,
) *Worker {
	return &Worker{
		concurrency:     cfg.Concurrency,
		queue:           cfg.Queue,
		rmqClient:       rmqClient,
		executorService: executorService,
		logger:          logger,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	w.ctx, w.cancel = context.WithCancel(ctx)

	if err := w.ensureQueueExists(); err != nil {
		return err
	}

	consumer := w.rmqClient.NewConsumer()
	err := consumer.Start(w.ctx, w.queue, w.messageHandler, w.concurrency)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) Stop() error {
	w.cancel()
	return nil
}

func (w *Worker) messageHandler(ctx context.Context, body []byte) error {
	executionRequest, err := w.executorService.ParseExecutionRequest(body)
	if err != nil {
		return err
	}

	_, err = w.executorService.Execute(ctx, executionRequest)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) ensureQueueExists() error {
	ch, err := w.rmqClient.Connection().Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(w.queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	return nil
}
