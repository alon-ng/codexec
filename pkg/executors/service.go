package executors

import (
	"codim/pkg/executors/drivers"
	"codim/pkg/executors/drivers/models"
	"codim/pkg/utils/logger"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Service struct {
	driver  drivers.Driver
	logger  *logger.Logger
	timeout time.Duration
}

func New(driver drivers.Driver, logger *logger.Logger, timeout time.Duration) *Service {
	return &Service{
		driver:  driver,
		logger:  logger,
		timeout: timeout,
	}
}

func (s *Service) Execute(ctx context.Context, executionRequest models.ExecutionRequest) (models.ExecuteResponse, error) {
	s.logger.Infof("Executing job %s", executionRequest.JobID)

	execCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.driver.Execute(execCtx, executionRequest)
	if err != nil {
		s.logger.Errorf("Failed to execute job %s: %v", executionRequest.JobID, err)
		return models.ExecuteResponse{}, err
	}

	s.logger.Infof("Executed job %s successfully", executionRequest.JobID)
	s.logger.Debugf(
		"Executed job %s with stdout %s, stderr %s, exit code %d, time %f, memory %d, cpu %f",
		executionRequest.JobID,
		res.Stdout,
		res.Stderr,
		res.ExitCode,
		res.Time,
		res.Memory,
		res.CPU,
	)
	return res, nil
}

func (s *Service) ParseExecutionRequest(body []byte) (models.ExecutionRequest, error) {
	var executionRequest models.ExecutionRequest
	if err := json.Unmarshal(body, &executionRequest); err != nil {
		return models.ExecutionRequest{}, fmt.Errorf("failed to unmarshal execution request: %w", err)
	}

	return executionRequest, nil
}
