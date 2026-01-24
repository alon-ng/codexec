package websocket

import (
	"codim/pkg/api/v1/errors"
	"codim/pkg/api/v1/modules/progress"
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
	"github.com/jackc/pgx/v5/pgxpool"
)

type JobClient struct {
	JobID        uuid.UUID
	ExerciseUuid uuid.UUID
	Client       *Client
}

type Hub struct {
	clients     map[*Client]bool
	jobClients  map[uuid.UUID]*JobClient
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
	progressSvc *progress.Service
}

func NewHub(rmqClient *rabbitmq.Client, logger *logger.Logger, q *db.Queries, p *pgxpool.Pool) *Hub {
	producer := rmqClient.NewProducer()
	consumer := rmqClient.NewConsumer()
	progressSvc := progress.NewService(q, p)
	return &Hub{
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		jobClients:  make(map[uuid.UUID]*JobClient),
		registerJob: make(chan *JobClient),
		rmqClient:   rmqClient,
		producer:    producer,
		consumer:    consumer,
		logger:      logger,
		q:           q,
		progressSvc: progressSvc,
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
			h.registerJobClient(jobClient)
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

func (h *Hub) registerJobClient(jobClient *JobClient) {
	h.jobMutex.Lock()
	h.jobClients[jobClient.JobID] = jobClient
	h.jobMutex.Unlock()
}

func (h *Hub) unregisterJobClient(jobID uuid.UUID) {
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

	h.jobMutex.Lock()
	jobClient, ok := h.jobClients[res.JobID]
	if ok {
		delete(h.jobClients, res.JobID)
	}
	h.jobMutex.Unlock()

	response := progress.UserExerciseSubmissionResponse{
		ExecuteResponse: res,
		Passed:          res.Passed(),
	}
	if response.Passed {
		nextLessonUuid, nextExerciseUuid, err := h.progressSvc.CompleteUserExercise(ctx, jobClient.Client.userID, jobClient.ExerciseUuid)
		if err != nil {
			errors.HandleError(nil, h.logger, err, http.StatusInternalServerError)
			return err.OriginalError
		}

		response.NextLessonUuid = nextLessonUuid
		response.NextExerciseUuid = nextExerciseUuid
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		errors.HandleError(nil, h.logger, errors.NewAPIError(err, "Internal server error"), http.StatusInternalServerError)
		return err
	}

	if ok {
		select {
		case jobClient.Client.send <- responseBytes:
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
