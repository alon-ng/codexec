package websocket

import (
	"codim/pkg/executors/drivers/models"
	"codim/pkg/rabbitmq"
	"codim/pkg/utils/logger"
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

type JobClient struct {
	JobID  string
	Client *Client
}

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Map jobID to Client
	jobClients map[string]*Client
	jobMutex   sync.Mutex

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Register job to client mapping
	registerJob chan *JobClient

	rmqClient *rabbitmq.Client
	producer  *rabbitmq.Producer
	consumer  *rabbitmq.Consumer
	logger    *logger.Logger
}

func NewHub(rmqClient *rabbitmq.Client, logger *logger.Logger) *Hub {
	producer := rmqClient.NewProducer()
	consumer := rmqClient.NewConsumer()
	return &Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		jobClients:  make(map[string]*Client),
		registerJob: make(chan *JobClient),
		rmqClient:   rmqClient,
		producer:    producer,
		consumer:    consumer,
		logger:      logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case jobClient := <-h.registerJob:
			h.jobMutex.Lock()
			h.jobClients[jobClient.JobID] = jobClient.Client
			h.jobMutex.Unlock()
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) ListenToRabbitMQ(ctx context.Context, exchangeName string) error {
	ch, err := h.rmqClient.Connection().Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	// 1. Declare exchange (Fanout)
	if err := ch.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// 2. Declare Temporary Queue (Exclusive)
	q, err := ch.QueueDeclare(
		"",    // name (empty means random)
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// 3. Bind Queue to Exchange
	if err := ch.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		exchangeName, // exchange
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	// 4. Start consuming
	handler := func(ctx context.Context, body []byte) error {
		var res models.ExecuteResponse
		if err := json.Unmarshal(body, &res); err != nil {
			h.logger.Errorf("failed to unmarshal execute response: %v", err)
			return nil // Don't requeue malformed messages
		}

		jobID := res.JobID.String()

		h.jobMutex.Lock()
		client, ok := h.jobClients[jobID]
		// Clean up job mapping
		if ok {
			delete(h.jobClients, jobID)
		}
		h.jobMutex.Unlock()

		if ok {
			// Send to client
			// We need to marshal it back to JSON or send raw bytes?
			// `res` is struct. `Client.send` is `chan []byte`.
			// `writePump` writes it as TextMessage.

			// Let's send the original body if we want to forward exact response,
			// or marshal `res` if we modified it or want to ensure clean JSON.
			// Using `body` is efficient.

			select {
			case client.send <- body:
			default:
				// Client buffer full or closed
			}
		}

		return nil
	}

	return h.consumer.Start(ctx, q.Name, handler, 1)
}
