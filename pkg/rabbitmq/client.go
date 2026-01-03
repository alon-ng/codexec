package rabbitmq

import (
	"codim/pkg/utils/logger"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn   *amqp.Connection
	cfg    Config
	logger *logger.Logger
}

// NewClient creates a new RabbitMQ client.
func NewClient(cfg Config, logger *logger.Logger) (*Client, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("rabbitmq URL is required")
	}

	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	logger.Info("Connected to RabbitMQ")

	return &Client{
		conn:   conn,
		cfg:    cfg,
		logger: logger,
	}, nil
}

// Close closes the RabbitMQ connection.
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

// Connection returns the underlying amqp connection.
func (c *Client) Connection() *amqp.Connection {
	return c.conn
}
