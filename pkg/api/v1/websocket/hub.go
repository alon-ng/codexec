package websocket

import (
	"codim/pkg/api/v1/errors"
	"codim/pkg/db"
	"codim/pkg/executors/drivers/models"
	"codim/pkg/rabbitmq"
	"codim/pkg/utils/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type JobClient struct {
	JobID  string
	Client *Client
}

type Hub struct {
	clients     map[*Client]bool
	jobClients  map[string]*Client
	jobMutex    sync.Mutex
	register    chan *Client
	unregister  chan *Client
	registerJob chan *JobClient
	rmqClient   *rabbitmq.Client
	producer    *rabbitmq.Producer
	consumer    *rabbitmq.Consumer
	logger      *logger.Logger
	q           *db.Queries
	upgrader    websocket.Upgrader
}

func NewHub(rmqClient *rabbitmq.Client, logger *logger.Logger, q *db.Queries) *Hub {
	producer := rmqClient.NewProducer()
	consumer := rmqClient.NewConsumer()
	return &Hub{
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		jobClients:  make(map[string]*Client),
		registerJob: make(chan *JobClient),
		rmqClient:   rmqClient,
		producer:    producer,
		consumer:    consumer,
		logger:      logger,
		q:           q,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case jobClient := <-h.registerJob:
			h.registerJobClient(jobClient.JobID, jobClient.Client)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.clients[client] = true
}

func (h *Hub) unregisterClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	}
}

func (h *Hub) registerJobClient(jobID string, client *Client) {
	h.jobMutex.Lock()
	h.jobClients[jobID] = client
	h.jobMutex.Unlock()
}

func (h *Hub) unregisterJobClient(jobID string) {
	h.jobMutex.Lock()
	delete(h.jobClients, jobID)
	h.jobMutex.Unlock()
}

func (h *Hub) ListenToRabbitMQ(ctx context.Context, exchangeName string) error {
	ch, err := h.rmqClient.Connection().Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare Temporary Queue (Exclusive)
	q, err := ch.QueueDeclare(
		"", // name (empty means random)
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := ch.QueueBind(
		q.Name,
		"",
		exchangeName,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return h.consumer.Start(ctx, q.Name, h.messageHandler, 1)
}

func (h *Hub) messageHandler(ctx context.Context, body []byte) error {
	var res models.ExecuteResponse
	if err := json.Unmarshal(body, &res); err != nil {
		h.logger.Errorf("failed to unmarshal execute response: %v", err)
		return err
	}

	jobID := res.JobID.String()

	h.jobMutex.Lock()
	client, ok := h.jobClients[jobID]
	if ok {
		delete(h.jobClients, jobID)
	}
	h.jobMutex.Unlock()

	if ok {
		select {
		case client.send <- body:
		default:
			// Client buffer full or closed
		}
	}

	return nil
}

func (h *Hub) ServeWs(c *gin.Context) {
	uuidStr, exists := c.Get("user_uuid")
	if !exists {
		errors.HandleError(c, h.logger, errors.NewAPIError(nil, "User UUID not found in context"), http.StatusUnauthorized)
		c.Abort()
		return
	}

	uuidStrValue, ok := uuidStr.(string)
	if !ok || uuidStrValue == "" {
		errors.HandleError(c, h.logger, errors.NewAPIError(nil, "Invalid user UUID in context"), http.StatusUnauthorized)
		c.Abort()
		return
	}

	userID, err := uuid.Parse(uuidStrValue)
	if err != nil {
		errors.HandleError(c, h.logger, errors.NewAPIError(err, "Failed to parse user UUID"), http.StatusBadRequest)
		c.Abort()
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Errorf("error upgrading websocket: %v", err)
		c.Abort()
		return
	}

	client := &Client{
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
		logger: h.logger,
		userID: userID,
		q:      h.q,
	}
	h.register <- client

	go client.writePump()
	go client.readPump()
}
