package models

import (
	"codim/pkg/executors/checkers"
	"codim/pkg/fs"

	"github.com/google/uuid"
)

type ExecutionRequest struct {
	JobID       uuid.UUID             `json:"job_id"`
	Source      fs.Entry              `json:"src"`
	EntryPoint  string                `json:"entry_point"`
	IOChecker   *checkers.IOChecker   `json:"io_data_checker,omitempty"`
	CodeChecker *checkers.CodeChecker `json:"code_checker,omitempty"`
}

type ExecuteResponse struct {
	JobID          uuid.UUID                `json:"job_id"`
	Stdout         string                   `json:"stdout"`
	Stderr         string                   `json:"stderr"`
	ExitCode       int                      `json:"exit_code"`
	Time           float64                  `json:"time"`
	Memory         int64                    `json:"memory"`
	CPU            float64                  `json:"cpu"`
	CheckerResults []checkers.CheckerResult `json:"checker_results"`
}

func (e *ExecuteResponse) Passed() bool {
	if e.ExitCode != 0 {
		return false
	}

	for _, checkerResult := range e.CheckerResults {
		if !checkerResult.Success {
			return false
		}
	}

	return true
}
