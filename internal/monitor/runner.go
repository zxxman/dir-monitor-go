package monitor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"dir-monitor-go/internal/logger"
	"dir-monitor-go/internal/model"
)

const (
	DefaultCommandTimeout   = 30
	CommandOutputBufferSize = 4096
)

type CommandExecutor struct {
	logger     *logger.Logger
	workingDir string
	envVars    map[string]string
}

func NewCommandExecutor(logger *logger.Logger, workingDir string) *CommandExecutor {
	return &CommandExecutor{
		logger:     logger,
		workingDir: workingDir,
		envVars:    make(map[string]string),
	}
}

func (ce *CommandExecutor) SetEnvVar(key, value string) {
	ce.envVars[key] = value
}

func (ce *CommandExecutor) ExecuteCommand(command string, event *model.FileEvent, timeout int) error {
	return ce.ExecuteCommandWithContext(context.Background(), command, event, timeout)
}

func (ce *CommandExecutor) ExecuteCommandWithContext(ctx context.Context, command string, event *model.FileEvent, timeout int) error {
	ce.logger.Info("Execute command: %s - File: %s", command, event.Path)

	if _, err := os.Stat(event.Path); err != nil {
		ce.logger.Error("File not found, skip execution: %s", event.Path)
		return fmt.Errorf("file not found: %s", event.Path)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	if err := validateCommand(command); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	cmd, err := ce.buildCommand(ctx, command, event)
	if err != nil {
		return fmt.Errorf("build command failed: %w", err)
	}

	ce.logger.Debug("Execute command: %v", cmd.Args)

	if ce.workingDir != "" {
		cmd.Dir = ce.workingDir
	}

	env := os.Environ()
	for k, v := range ce.envVars {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = env

	return ce.runCommand(ctx, cmd)
}

func (ce *CommandExecutor) buildCommand(ctx context.Context, command string, event *model.FileEvent) (*exec.Cmd, error) {
	commandLine := ce.replaceCommandVariables(command, event)
	if strings.TrimSpace(commandLine) == "" {
		return nil, fmt.Errorf("empty command")
	}

	var cmd *exec.Cmd
	if isWindows() {
		cmd = exec.CommandContext(ctx, "cmd.exe", "/c", commandLine)
	} else {
		cmd = exec.CommandContext(ctx, "/bin/sh", "-c", commandLine)
	}
	setProcessGroup(cmd)
	return cmd, nil
}

func (ce *CommandExecutor) runCommand(ctx context.Context, cmd *exec.Cmd) error {
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("command start failed: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		output := strings.TrimSpace(outBuf.String() + "\n" + errBuf.String())
		if output != "" {
			ce.logger.Debug("Command output: %s", output)
		}
		if err != nil {
			return fmt.Errorf("command execution failed: %w", err)
		}
		return nil

	case <-ctx.Done():
		if cmd.Process != nil {
			pid := cmd.Process.Pid
			_ = killProcessTree(pid)
		}
		output := strings.TrimSpace(outBuf.String() + "\n" + errBuf.String())
		if output != "" {
			ce.logger.Debug("Command terminated output: %s", output)
		}
		return fmt.Errorf("command cancelled or timeout: %w", ctx.Err())
	}
}

func (ce *CommandExecutor) replaceCommandVariables(command string, event *model.FileEvent) string {
	command = strings.ReplaceAll(command, "${EVENT_TYPE}", string(event.Type))
	command = strings.ReplaceAll(command, "${FILE_PATH}", event.Path)
	command = strings.ReplaceAll(command, "${FILE_NAME}", filepath.Base(event.Path))
	command = strings.ReplaceAll(command, "${FILE_DIR}", filepath.Dir(event.Path))
	command = strings.ReplaceAll(command, "${EVENT_TIME}", event.Timestamp.Format(time.RFC3339))

	for k, v := range ce.envVars {
		command = strings.ReplaceAll(command, "${"+k+"}", v)
	}

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			command = strings.ReplaceAll(command, "${"+parts[0]+"}", parts[1])
		}
	}

	return command
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func validateCommand(command string) error {
	if strings.TrimSpace(command) == "" {
		return fmt.Errorf("empty command")
	}
	return nil
}
