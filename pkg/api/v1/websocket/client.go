package websocket

import (
	"codim/pkg/api/auth"
	"codim/pkg/api/v1/cache"
	"codim/pkg/db"
	"codim/pkg/executors"
	"codim/pkg/fs"
	"codim/pkg/utils/logger"
	"context"
	"encoding/json"
	"net/http"
	"time"

	v "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, logger *logger.Logger, authProvider *auth.Provider, userCache *cache.UserCache, q *db.Queries, w http.ResponseWriter, r *http.Request) {
	// Check authentication cookie before upgrading
	cookie, err := r.Cookie(auth.AuthCookieName)
	if err != nil || cookie == nil || cookie.Value == "" {
		logger.Warn("WebSocket connection rejected: no auth cookie")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Verify token
	userUUID, renewalRequired, err := authProvider.VerifyToken(cookie.Value)
	if err != nil {
		logger.Warnf("WebSocket connection rejected: invalid token: %v", err)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Get user from cache (verify user exists)
	_, err = userCache.GetUser(r.Context(), userUUID)
	if err != nil {
		logger.Errorf("WebSocket connection rejected: failed to get user: %v", err)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Renew token if needed
	if renewalRequired {
		newToken, err := authProvider.GenerateToken(userUUID)
		if err != nil {
			logger.Errorf("Failed to generate new token for WebSocket: %v", err)
			// Continue anyway, token is still valid
		} else {
			// Set new token cookie
			http.SetCookie(w, &http.Cookie{
				Name:     auth.AuthCookieName,
				Value:    newToken,
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
		}
	}

	// Enable CORS for WebSocket (consider restricting this in production)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorf("error upgrading websocket: %v", err)
		return
	}

	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		logger: logger,
		userID: userUUID,
		q:      q,
	}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	logger *logger.Logger
	userID uuid.UUID
	q      *db.Queries
}

type SubmissionMessage struct {
	ExerciseUuid uuid.UUID `json:"exercise_uuid" validate:"required"`
	Submission   fs.Entry  `json:"submission" validate:"required"`
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	validate := v.New()
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Errorf("error: %v", err)
			}
			break
		}

		// Parse submission
		var submission SubmissionMessage
		if err := json.Unmarshal(message, &submission); err != nil {
			c.logger.Errorf("error parsing submission: %v", err)
			continue
		}

		if err := validate.Struct(submission); err != nil {
			c.logger.Errorf("error validating submission: %v", err)
			continue
		}

		// Create execution request
		jobID := uuid.New()

		row, err := c.q.GetExerciseSubjectAndType(context.Background(), submission.ExerciseUuid)
		if err != nil {
			c.logger.Errorf("error getting exercise subject and type: %v", err)
			continue
		}

		if row.Type == db.ExerciseTypeCode {
			runCodeSubmission(c, jobID, submission, row.Subject)
		}
	}
}

func runCodeSubmission(c *Client, jobID uuid.UUID, submission SubmissionMessage, subject string) {
	// Map language to queue/driver
	// For now we assume queue names based on language
	queueName := "codexec." + subject

	req := executors.ExecutionRequest{
		JobID:      jobID,
		Source:     submission.Submission,
		EntryPoint: "main." + getExtension(subject),
	}

	c.hub.registerJob <- &JobClient{
		JobID:  jobID.String(),
		Client: c,
	}

	// Publish to default exchange (empty string) with queue name as routing key
	// This routes directly to the queue without needing to declare an exchange
	err := c.hub.producer.PublishObject(context.Background(), "", queueName, req)
	if err != nil {
		c.logger.Errorf("error publishing to rabbitmq: %v", err)
	}
}

func getExtension(lang string) string {
	switch lang {
	case "python":
		return "py"
	case "node", "javascript":
		return "js"
	default:
		return "txt"
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
