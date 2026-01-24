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
)

func Execute(
	ctx context.Context,
	cmdPrefix string,
	nsjailConfigTemplate string,
	testUtilsFile string,
	executionRequest models.ExecutionRequest,
) (models.ExecuteResponse, error) {
	jobIDStr := executionRequest.JobID.String()
	jobPath := fmt.Sprintf("/jobs/%s", jobIDStr)

	// Build nsjail config with replaced placeholders
	config := prepareNsjailConfig(nsjailConfigTemplate, jobIDStr, jobIDStr, executionRequest.EntryPoint)

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

	defer DeleteFile(ctx, cmdPrefix, cfgPath)

	// Write files to container
	if err := WriteFiles(ctx, cmdPrefix, jobPath, executionRequest.Source); err != nil {
		return models.ExecuteResponse{}, fmt.Errorf("failed to write files: %w", err)
	}

	// Execute nsjail
	r, err := ExecuteNsjail(ctx, cmdPrefix, cfgPath)

	if err != nil {
		return models.ExecuteResponse{}, err
	}

	runCheckers(
		ctx,
		executionRequest,
		&r,
		cmdPrefix,
		nsjailConfigTemplate,
		jobPath,
		testUtilsFile,
	)

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

func DeleteFile(ctx context.Context, cmdPrefix string, filePath string) error {
	cmd := executeCommand(ctx, cmdPrefix, "rm", "-f", filePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete config file %s: %w", filePath, err)
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
		Stdout:   strings.TrimSpace(stdout.String()),
		Stderr:   strings.TrimSpace(stderr.String()),
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

func runCheckers(
	ctx context.Context,
	request models.ExecutionRequest,
	response *models.ExecuteResponse,
	cmdPrefix string,
	nsjailConfigTemplate string,
	jobPath string,
	testUtilsFile string,
) error {
	if request.IOChecker != nil {
		r := request.IOChecker.Check(ctx, response.Stdout)
		response.CheckerResults = append(response.CheckerResults, r)
	}

	if request.CodeChecker != nil {
		testFilePath := fmt.Sprintf("%s/%s", jobPath, request.CodeChecker.FileName)
		err := WriteFile(ctx, cmdPrefix, testFilePath, request.CodeChecker.Code)
		if err != nil {
			return err
		}

		testUtilsFilePath := fmt.Sprintf("%s/%s", jobPath, "test_utils.py")
		err = WriteFile(ctx, cmdPrefix, testUtilsFilePath, testUtilsFile)
		if err != nil {
			return err
		}

		testJobId := fmt.Sprintf("%s-tests", request.JobID.String())
		cfgPath := fmt.Sprintf("/tmp/config-%s.cfg", testJobId)
		config := prepareNsjailConfig(nsjailConfigTemplate, testJobId, request.JobID.String(), request.CodeChecker.FileName)

		err = CreateConfigFile(ctx, cmdPrefix, cfgPath, config)
		if err != nil {
			return err
		}

		defer DeleteFile(ctx, cmdPrefix, cfgPath)

		r, err := ExecuteNsjail(ctx, cmdPrefix, cfgPath)
		if err != nil {
			return err
		}

		rs := request.CodeChecker.Check(ctx, r.Stdout)

		response.CheckerResults = append(response.CheckerResults, rs...)
	}

	return nil
}

func prepareNsjailConfig(config string, jobId string, jobFolder string, entryPoint string) string {
	config = strings.ReplaceAll(config, "{{JOB_ID}}", jobId)
	config = strings.ReplaceAll(config, "{{JOB_ID_FOLDER}}", jobFolder)
	config = strings.ReplaceAll(config, "{{ENTRY_POINT}}", entryPoint)
	return config
}
