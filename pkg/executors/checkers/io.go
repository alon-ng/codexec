package checkers

import (
	"context"
	"fmt"
)

type IOChecker struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
}

func (c *IOChecker) Check(ctx context.Context, stdout string) CheckerResult {
	if stdout == c.ExpectedOutput {
		return CheckerResult{
			Type:    CheckerTypeIO,
			Success: true,
			Message: "Output matches expected output",
		}
	}

	return CheckerResult{
		Type:    CheckerTypeIO,
		Success: false,
		Message: fmt.Sprintf("Expected output: %s, Actual output: %s", c.ExpectedOutput, stdout),
	}
}
