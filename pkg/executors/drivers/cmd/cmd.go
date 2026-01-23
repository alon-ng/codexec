package cmd

import (
	"bytes"
	"codim/pkg/executors/drivers/models"
	"codim/pkg/fs"
	"context"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
)

func Execute(ctx context.Context, cmdPrefix string, nsjailConfigTemplate string, jobID uuid.UUID, src fs.Entry, entryPoint string) (models.ExecuteResponse, error) {
	jobIDStr := jobID.String()
	jobPath := fmt.Sprintf("/jobs/%s", jobIDStr)

	// Build nsjail config with replaced placeholders
	config := strings.ReplaceAll(nsjailConfigTemplate, "{{JOB_ID}}", jobIDStr)
	config = strings.ReplaceAll(config, "{{ENTRY_POINT}}", entryPoint)

	// Create job directory in container
	if err := CreateJobDirectory(ctx, cmdPrefix, jobPath); err != nil {
		return models.ExecuteResponse{}, fmt.Errorf("failed to create job directory: %w", err)
	}

	defer DeleteJobDirectory(ctx, cmdPrefix, jobPath)

	// Create config file in container
	cfgPath := fmt.Sprintf("/tmp/config-%s.cfg", jobIDStr)
	if err := CreateConfigFile(ctx, cmdPrefix, cfgPath, config); err != nil {
		return models.ExecuteResponse{}, fmt.Errorf("failed to create config file: %w", err)
	}

	defer DeleteConfigFile(ctx, cmdPrefix, cfgPath)

	// Write files to container
	if err := WriteFiles(ctx, cmdPrefix, jobPath, src); err != nil {
		return models.ExecuteResponse{}, fmt.Errorf("failed to write files: %w", err)
	}

	// Execute nsjail
	r, err := ExecuteNsjail(ctx, cmdPrefix, cfgPath)

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

func WriteFiles(ctx context.Context, cmdPrefix string, basePath string, entry fs.Entry) error {
	return WriteEntry(ctx, cmdPrefix, basePath, entry)
}

func WriteEntry(ctx context.Context, cmdPrefix string, basePath string, entry fs.Entry) error {
	for _, child := range entry.Children {
		childPath := fmt.Sprintf("%s/%s", basePath, child.Name)

		if len(child.Children) > 0 {
			cmd := executeCommand(ctx, cmdPrefix, "mkdir", "-p", childPath)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to create subdirectory %s: %w", childPath, err)
			}
			if err := WriteEntry(ctx, cmdPrefix, childPath, child); err != nil {
				return err
			}
		} else if child.Content != "" {
			if err := WriteFile(ctx, cmdPrefix, childPath, child.Content); err != nil {
				return err
			}
		} else {
			cmd := executeCommand(ctx, cmdPrefix, "mkdir", "-p", childPath)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to create empty directory %s: %w", childPath, err)
			}
		}
	}

	if entry.Content != "" {
		filePath := fmt.Sprintf("%s/%s", basePath, entry.Name)
		if err := WriteFile(ctx, cmdPrefix, filePath, entry.Content); err != nil {
			return err
		}
	}

	return nil
}

func DeleteConfigFile(ctx context.Context, cmdPrefix string, cfgPath string) error {
	cmd := executeCommand(ctx, cmdPrefix, "rm", "-f", cfgPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete config file %s: %w", cfgPath, err)
	}
	return nil
}

func CreateConfigFile(ctx context.Context, cmdPrefix string, cfgPath string, config string) error {
	return WriteFile(ctx, cmdPrefix, cfgPath, config)
}

func ExecuteNsjail(ctx context.Context, cmdPrefix string, cfgPath string) (models.ExecuteResponse, error) {
	// Execute nsjail with config piped to stdin
	// Use -Q flag to suppress nsjail's verbose logging (only show errors)
	cmd := executeCommand(ctx, cmdPrefix, "nsjail", "-Q", "--config", cfgPath)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	var cpuTime time.Duration
	var maxMemory int64
	exitCode := 0

	err := cmd.Run()
	wallTime := time.Since(start)

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			return models.ExecuteResponse{}, fmt.Errorf("failed to execute nsjail: %w", err)
		}
	}

	if cmd.ProcessState != nil {
		cpuTime = cmd.ProcessState.UserTime() + cmd.ProcessState.SystemTime()
		usage := cmd.ProcessState.SysUsage().(*syscall.Rusage)
		maxMemory = usage.Maxrss / 1024
	}

	return models.ExecuteResponse{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Time:     wallTime.Seconds(),
		CPU:      cpuTime.Seconds(),
		Memory:   maxMemory,
	}, nil
}

func WriteFile(ctx context.Context, cmdPrefix string, filePath string, content string) error {
	base64Content := base64.StdEncoding.EncodeToString([]byte(content))
	cmd := executeCommand(ctx, cmdPrefix, "sh", "-c", fmt.Sprintf("echo '%s' | base64 -d > %s", base64Content, filePath))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}
	return nil
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
