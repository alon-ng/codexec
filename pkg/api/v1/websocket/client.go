package websocket

import (
	"codim/pkg/api/v1/modules/progress"
	"codim/pkg/db"
	"codim/pkg/executors/checkers"
	"codim/pkg/executors/drivers/models"
	"codim/pkg/fs"
	"codim/pkg/utils/logger"
	"context"
	"encoding/json"
	"fmt"
	"time"

	v "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	logger *logger.Logger
	userID uuid.UUID
	q      *db.Queries
}

type SubmissionMessage struct {
	ExerciseUuid uuid.UUID   `json:"exercise_uuid" validate:"required"`
	Submission   interface{} `json:"submission" validate:"required"`
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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
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

		row, err := c.q.GetExerciseForSubmission(context.Background(), submission.ExerciseUuid)
		if err != nil {
			c.logger.Errorf("error getting exercise subject and type: %v", err)
			continue
		}

		switch row.Type {
		case db.ExerciseTypeCode:
			runCodeSubmission(c, submission, row)
		case db.ExerciseTypeQuiz:
			runQuizSubmission(c, submission, row)
		}
	}
}

func runCodeSubmission(c *Client, submission SubmissionMessage, exercise db.GetExerciseForSubmissionRow) {
	// Unmarshal the submission into the quiz submission struct
	submissionBytes, err := json.Marshal(submission.Submission)
	if err != nil {
		c.logger.Errorf("error marshalling submission: %v", err)
		return
	}

	var codeSubmission progress.UserExerciseSubmissionCode
	if err := json.Unmarshal(submissionBytes, &codeSubmission); err != nil {
		c.logger.Errorf("error unmarshalling code submission: %v", err)
		return
	}

	jobID := uuid.New()
	queueName := "codexec." + exercise.Subject

	var codeChecker *checkers.CodeChecker
	var ioChecker *checkers.IOChecker
	if exercise.CodeChecker != nil {
		if err := json.Unmarshal(*exercise.CodeChecker, &codeChecker); err != nil {
			c.logger.Errorf("error unmarshalling code checker: %v", err)
			return
		}
	}
	if exercise.IoChecker != nil {
		if err := json.Unmarshal(*exercise.IoChecker, &ioChecker); err != nil {
			c.logger.Errorf("error unmarshalling io checker: %v", err)
			return
		}
	}

	req := models.ExecutionRequest{
		JobID:       jobID,
		Source:      fs.Entry(codeSubmission),
		EntryPoint:  "main." + getExtension(exercise.Subject),
		CodeChecker: codeChecker,
		IOChecker:   ioChecker,
	}

	c.hub.registerJob <- &JobClient{
		JobID:        jobID,
		ExerciseUuid: submission.ExerciseUuid,
		Client:       c,
	}

	err = c.hub.producer.PublishObject(context.Background(), "", queueName, req)
	if err != nil {
		c.logger.Errorf("error publishing to rabbitmq: %v", err)
	}
}

func runQuizSubmission(c *Client, submission SubmissionMessage, exercise db.GetExerciseForSubmissionRow) {
	// Unmarshal the submission into the quiz submission struct
	submissionBytes, err := json.Marshal(submission.Submission)
	if err != nil {
		c.logger.Errorf("error marshalling submission: %v", err)
		return
	}

	var quizSubmission progress.UserExerciseSubmissionQuiz
	if err := json.Unmarshal(submissionBytes, &quizSubmission); err != nil {
		c.logger.Errorf("error unmarshalling quiz submission: %v", err)
		return
	}

	var quizChecker map[string]string
	if exercise.QuizChecker != nil {
		if err := json.Unmarshal(*exercise.QuizChecker, &quizChecker); err != nil {
			c.logger.Errorf("error unmarshalling quiz checker: %v", err)
			return
		}
	}

	jobID := uuid.New()
	checkerResults := make([]checkers.CheckerResult, 0)

	for question, answer := range quizSubmission.Answers {
		correctAnswer, ok := quizChecker[question]
		if !ok {
			c.logger.Warnf("question %s not found in quiz data", question)
			checkerResults = append(checkerResults, checkers.CheckerResult{
				Type:    checkers.CheckerType(question),
				Success: false,
				Message: "",
			})
			continue
		}

		markAsCorrect := answer == correctAnswer
		checkerResults = append(checkerResults, checkers.CheckerResult{
			Type:    checkers.CheckerType(question),
			Success: markAsCorrect,
			Message: "",
		})
	}

	res := models.ExecuteResponse{
		JobID:          jobID,
		Stdout:         "",
		Stderr:         "",
		ExitCode:       0,
		Time:           0,
		Memory:         0,
		CPU:            0,
		CheckerResults: checkerResults,
	}

	// Register job client so response can be routed back
	c.hub.registerJob <- &JobClient{
		JobID:        jobID,
		ExerciseUuid: submission.ExerciseUuid,
		Client:       c,
	}

	// Publish response to rabbitmq
	err = c.hub.producer.PublishObject(context.Background(), "codexec.results", "", res)
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
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.writeMessage(message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.ping(); err != nil {
				return
			}
		}
	}
}

func (c *Client) ping() error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(websocket.PingMessage, nil)
}

func (c *Client) writeMessage(message []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))

	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return fmt.Errorf("error getting next writer: %v", err)
	}
	w.Write(message)

	// Add queued chat messages to the current websocket message.
	n := len(c.send)
	for i := 0; i < n; i++ {
		w.Write(newline)
		w.Write(<-c.send)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("error closing writer: %v", err)
	}

	return nil
}
