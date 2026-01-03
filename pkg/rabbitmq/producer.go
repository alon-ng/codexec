package rabbitmq

import (
	"codim/pkg/utils/logger"
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Producer handles message publishing to RabbitMQ.
type Producer struct {
	conn   *amqp.Connection
	logger *logger.Logger
}

// NewProducer creates a new Producer instance.
func (c *Client) NewProducer() *Producer {
	return &Producer{
		conn:   c.conn,
		logger: c.logger,
	}
}

// PublishOption allows customizing the publishing behavior.
type PublishOption func(*amqp.Publishing)

// WithContentType sets the content type of the message.
func WithContentType(contentType string) PublishOption {
	return func(p *amqp.Publishing) {
		p.ContentType = contentType
	}
}

// WithHeaders sets the headers of the message.
func WithHeaders(headers amqp.Table) PublishOption {
	return func(p *amqp.Publishing) {
		p.Headers = headers
	}
}

// Publish sends a message to the specified exchange with the given routing key.
// It creates a temporary channel for thread safety.
func (p *Producer) Publish(ctx context.Context, exchange, routingKey string, body []byte, opts ...PublishOption) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}

	for _, opt := range opts {
		opt(&msg)
	}

	if err := ch.PublishWithContext(ctx, exchange, routingKey, false, false, msg); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// PublishObject publishes an object by marshaling it to JSON.
// It only accepts struct types or pointers to structs, rejecting primitives.
func (p *Producer) PublishObject(ctx context.Context, exchange, routingKey string, obj any, opts ...PublishOption) error {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("PublishObject only accepts struct types or pointers to structs, got %T", obj)
	}

	body, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal object: %w", err)
	}

	return p.Publish(ctx, exchange, routingKey, body, opts...)
}
