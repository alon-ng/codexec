package rabbitmq

import (
	"codim/internal/utils/logger"
	"context"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerOption func(*Consumer)

func WithConsumerTag(tag string) ConsumerOption {
	return func(c *Consumer) {
		c.tag = tag
	}
}

func WithAutoAck(autoAck bool) ConsumerOption {
	return func(c *Consumer) {
		c.autoAck = autoAck
	}
}

func WithExclusive(exclusive bool) ConsumerOption {
	return func(c *Consumer) {
		c.exclusive = exclusive
	}
}

func WithNoLocal(noLocal bool) ConsumerOption {
	return func(c *Consumer) {
		c.noLocal = noLocal
	}
}

func WithNoWait(noWait bool) ConsumerOption {
	return func(c *Consumer) {
		c.noWait = noWait
	}
}

func WithArgs(args amqp.Table) ConsumerOption {
	return func(c *Consumer) {
		c.args = args
	}
}

// Consumer handles message consumption from RabbitMQ.
type Consumer struct {
	conn      *amqp.Connection
	logger    *logger.Logger
	tag       string
	autoAck   bool
	exclusive bool
	noLocal   bool
	noWait    bool
	args      amqp.Table
}

// Handler is the function signature for processing messages.
// It receives the message body and returns an error.
// If error is nil, the message is Acknowledged.
// If error is not nil, the message is Negative Acknowledged and requeued.
type Handler func(ctx context.Context, body []byte) error

// NewConsumer creates a new Consumer instance.
func (c *Client) NewConsumer() *Consumer {
	return &Consumer{
		conn:      c.conn,
		logger:    c.logger,
		tag:       "",
		autoAck:   false,
		exclusive: false,
		noLocal:   false,
		noWait:    false,
		args:      nil,
	}
}

// Start begins consuming messages from the specified queue with the given concurrency.
// It runs in the background until the context is cancelled or the connection is lost.
func (c *Consumer) Start(ctx context.Context, queue string, handler Handler, concurrency int, opts ...ConsumerOption) error {
	for _, opt := range opts {
		opt(c)
	}

	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	// Set QoS to ensure we don't overwhelm the consumers.
	// Prefetch count is multiplied by 5 to ensure we don't overwhelm the consumers but still keep the workers busy.
	if err := ch.Qos(concurrency*5, 0, false); err != nil {
		_ = ch.Close()
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := ch.ConsumeWithContext(
		ctx,
		queue,
		c.tag,
		c.autoAck,
		c.exclusive,
		c.noLocal,
		c.noWait,
		c.args,
	)
	if err != nil {
		_ = ch.Close()
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.logger.Infof("Starting consumer for queue %s with concurrency %d", queue, concurrency)

	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-msgs:
					if !ok {
						c.logger.Warn("Channel closed, stopping worker")
						return
					}
					c.processMessage(ctx, msg, handler)
				}
			}
		}()
	}

	// Monitor context cancellation to cleanup
	go func() {
		<-ctx.Done()
		c.logger.Info("Context cancelled, closing consumer channel")
		_ = ch.Close()
	}()

	wg.Wait()
	c.logger.Info("Consumer stopped")

	return nil
}

func (c *Consumer) processMessage(ctx context.Context, msg amqp.Delivery, handler Handler) {
	c.logger.Debugf("Received message: %s with message id %s from queue %s", string(msg.Body), msg.MessageId, msg.RoutingKey)
	defer func() {
		if r := recover(); r != nil {
			c.logger.Errorf("Panic in consumer handler: %v", r)
			_ = msg.Nack(false, true)
		}
	}()

	if err := handler(ctx, msg.Body); err != nil {
		c.logger.Errorf("Failed to process message: %v", err)
		if nackErr := msg.Nack(false, true); nackErr != nil {
			c.logger.Errorf("Failed to nack message: %v", nackErr)
		}
		return
	}

	if ackErr := msg.Ack(false); ackErr != nil {
		c.logger.Errorf("Failed to ack message: %v", ackErr)
	}
}
