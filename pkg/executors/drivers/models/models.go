package models

import "github.com/google/uuid"

type ExecuteResponse struct {
	JobID    uuid.UUID `json:"job_id"`
	Stdout   string    `json:"stdout"`
	Stderr   string    `json:"stderr"`
	ExitCode int       `json:"exit_code"`
	Time     float64   `json:"time"`
	Memory   int64     `json:"memory"`
	CPU      float64   `json:"cpu"`
}
