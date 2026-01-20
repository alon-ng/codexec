package cmd

import (
	"bytes"
	"codim/pkg/executors/drivers/models"
	"codim/pkg/fs"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func Execute(ctx context.Context, cmdPrefix string, nsjailConfigTemplate string, jobID uuid.UUID, src fs.Directory, entryPoint string) (models.ExecuteResponse, error) {
	jobIDStr := jobID.String()
	jobPath := fmt.Sprintf("/jobs/%s", jobIDStr)

	// Create job directory in container
	if err := CreateJobDirectory(ctx, cmdPrefix, jobPath); err != nil {
		return models.ExecuteResponse{}, fmt.Errorf("failed to create job directory: %w", err)
	}

	defer DeleteJobDirectory(ctx, cmdPrefix, jobPath)

	// Write files to container
	if err := WriteFiles(ctx, cmdPrefix, jobPath, src); err != nil {
		return models.ExecuteResponse{}, fmt.Errorf("failed to write files: %w", err)
	}

	// Build nsjail config with replaced placeholders
	config := strings.ReplaceAll(nsjailConfigTemplate, "{{JOB_ID}}", jobIDStr)
	config = strings.ReplaceAll(config, "{{ENTRY_POINT}}", entryPoint)

	// Execute nsjail
	r, err := ExecuteNsjail(ctx, cmdPrefix, config)

	return r, err
}

func DeleteJobDirectory(ctx context.Context, cmdPrefix string, jobPath string) error {
	cmd := executeCommand(ctx, cmdPrefix, "rm", "-rf", jobPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete job directory %s: %w", jobPath, err)
	}

	return nil
}

func CreateJobDirectory(ctx context.Context, cmdPrefix string, basePath string) error {
	cmd := executeCommand(ctx, cmdPrefix, "mkdir", "-p", basePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", basePath, err)
	}
	return nil
}

func WriteFiles(ctx context.Context, cmdPrefix string, basePath string, dir fs.Directory) error {
	return WriteDirectory(ctx, cmdPrefix, basePath, dir)
}

func WriteDirectory(ctx context.Context, cmdPrefix string, basePath string, dir fs.Directory) error {
	// Write files in current directory
	for _, file := range dir.Files {
		filePath := fmt.Sprintf("%s/%s%s", basePath, file.Name, file.Ext)
		cmd := executeCommand(ctx, cmdPrefix, "sh", "-c", fmt.Sprintf("cat > %s", filePath))
		cmd.Stdin = strings.NewReader(file.Content)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	// Recursively write subdirectories
	for _, subdir := range dir.Directories {
		subdirPath := fmt.Sprintf("%s/%s", basePath, subdir.Name)
		// Create subdirectory
		cmd := executeCommand(ctx, cmdPrefix, "mkdir", "-p", subdirPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create subdirectory %s: %w", subdirPath, err)
		}
		// Recursively write subdirectory contents
		if err := WriteDirectory(ctx, cmdPrefix, subdirPath, subdir); err != nil {
			return err
		}
	}

	return nil
}

func ExecuteNsjail(ctx context.Context, cmdPrefix string, config string) (models.ExecuteResponse, error) {
	// Execute nsjail with config piped to stdin
	// Use -Q flag to suppress nsjail's verbose logging (only show errors)
	cmd := executeCommand(ctx, cmdPrefix, "nsjail", "-Q", "--config", "/dev/stdin")
	cmd.Stdin = strings.NewReader(config)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			return models.ExecuteResponse{}, fmt.Errorf("failed to execute nsjail: %w", err)
		}
	}

	stderrStr, maybeMetrics, err := seperateStderr(stderr.String())

	if err != nil {
		return models.ExecuteResponse{}, fmt.Errorf("failed to seperate stderr: %w", err)
	}

	time, cpu, memoryKb, err := parseMetrics(maybeMetrics)
	if err != nil {
		return models.ExecuteResponse{}, fmt.Errorf("failed to parse metrics: %w", err)
	}

	return models.ExecuteResponse{
		Stdout:   stdout.String(),
		Stderr:   stderrStr,
		ExitCode: exitCode,
		Time:     time,
		CPU:      cpu,
		Memory:   memoryKb,
	}, nil
}

func parseMetrics(metricsLine string) (float64, float64, int64, error) {
	metrics := strings.Split(metricsLine, ",")

	if len(metrics) != 4 {
		return 0, 0, 0, fmt.Errorf("failed to parse metrics, expected 4 values, got %d", len(metrics))
	}

	time, err := strconv.ParseFloat(metrics[0], 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse time: %w", err)
	}
	usrCpu, err := strconv.ParseFloat(metrics[1], 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse user cpu: %w", err)
	}
	sysCpu, err := strconv.ParseFloat(metrics[2], 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse system cpu: %w", err)
	}
	memoryKb, err := strconv.ParseInt(metrics[3], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse memory: %w", err)
	}

	cpu := (usrCpu + sysCpu) / time

	return time, cpu, memoryKb, nil
}

func seperateStderr(stderr string) (string, string, error) {
	lines := strings.Split(stderr, "\n")

	if len(lines) < 2 {
		return "", "", fmt.Errorf("failed to parse stderr, expected at least 2 lines, got %d", len(lines))
	}

	// Everything but the 2 last lines
	errors := lines[:len(lines)-2]

	return strings.Join(errors, "\n"), lines[len(lines)-2], nil
}

func executeCommand(ctx context.Context, cmdPrefix string, cmdArgs ...string) *exec.Cmd {
	if len(cmdArgs) == 0 {
		return nil
	}

	if cmdPrefix == "" {
		baseCmd := cmdArgs[0]
		return exec.CommandContext(ctx, baseCmd, cmdArgs[1:]...)
	}

	baseCmd, baseCmdArgs := deconstructCmdPrefix(cmdPrefix)

	args := append(baseCmdArgs, cmdArgs...)
	return exec.CommandContext(ctx, baseCmd, args...)
}

func deconstructCmdPrefix(cmdPrefix string) (string, []string) {
	parts := strings.Split(cmdPrefix, " ")
	return parts[0], parts[1:]
}
