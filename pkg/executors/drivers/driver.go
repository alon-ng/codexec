package drivers

import (
	"codim/pkg/executors/drivers/models"
	"codim/pkg/executors/drivers/node"
	"codim/pkg/executors/drivers/python"
	"codim/pkg/utils/logger"
	"context"
	"fmt"
)

type Driver interface {
	Execute(ctx context.Context, executionRequest models.ExecutionRequest) (models.ExecuteResponse, error)
	SetCmdPrefix(prefix string) error
	CmdPrefix() string
}

func New(driver, cmdPrefix string, logger *logger.Logger) (Driver, error) {
	switch driver {
	case "python":
		return python.New(cmdPrefix, logger), nil
	case "node":
		return node.New(cmdPrefix, logger), nil
	default:
		return nil, fmt.Errorf("driver %s is invalid", driver)
	}
}
