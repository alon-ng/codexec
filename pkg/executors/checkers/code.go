package checkers

import (
	"context"
	"encoding/json"
	"strings"
)

type CodeChecker struct {
	Code     string `json:"code"`
	FileName string `json:"file_name"`
}

type codeCheckerResult struct {
	IsTest  bool   `json:"is_test"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (c *CodeChecker) Check(ctx context.Context, stdout string) []CheckerResult {
	lines := strings.Split(stdout, "\n")
	results := make([]CheckerResult, 0)
	for _, line := range lines {
		var result codeCheckerResult
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			continue
		}

		if !result.IsTest {
			continue
		}

		results = append(results, CheckerResult{
			Type:    CheckerTypeCode,
			Success: result.Success,
			Message: result.Message,
		})
	}

	return results
}
