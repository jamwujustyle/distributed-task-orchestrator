package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/aws/aws-lambda-go/lambda"
)

var allowedCommands = []string{
	"echo",
	"cat",
	"grep",
	"sleep",
	"date",
}

func setupSandboxBin() (string, error) {
	sandboxDir := "/tmp/sandbox_bin"

	if err := os.MkdirAll(sandboxDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create sandbox bin: %w", err)
	}

	for _, cmdName := range allowedCommands {
		fullPath, err := exec.LookPath(cmdName)
		if err != nil {
			continue
		}
		targetLink := filepath.Join(sandboxDir, cmdName)

		_ = os.Remove(targetLink)

		if err := os.Symlink(fullPath, targetLink); err != nil {
			return "", fmt.Errorf("failed to symlink command: %s: %w", cmdName, err)
		}
	}
	return sandboxDir, nil
}

func validateScript(scriptContent string) error {
	lines := strings.Split(scriptContent, "\n")
	fmt.Printf("%#v\n", lines)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.Fields(trimmed)
		if len(parts) == 0 {
			continue
		}
		firstWord := parts[0]
		if strings.Contains(firstWord, "/") {
			return fmt.Errorf("use of absolute or relative paths with '/' in command names is forbidden: %q", firstWord)
		}
	}
	return nil
}

type ScriptRequest struct {
	Script string `json:"script"`
}

type ScriptResponse struct {
	Stdout   string `json:"stdout"`
	ExitCode int    `json:"exit_code"`
}

func handler(ctx context.Context, req ScriptRequest) (ScriptResponse, error) {
	data, err := base64.StdEncoding.DecodeString(req.Script)
	if err != nil {
		return ScriptResponse{}, fmt.Errorf("failed to decode script: %w", err)
	}

	scriptStr := string(data)

	if err := validateScript(scriptStr); err != nil {
		return ScriptResponse{
			Stdout:   fmt.Sprintf("Security error: %v\n", err),
			ExitCode: 126,
		}, nil
	}

	sandboxDir, err := setupSandboxBin()
	if err != nil {
		return ScriptResponse{}, fmt.Errorf("failed to configure sandbox environment: %w", err)
	}

	tmp, err := os.CreateTemp("", "task-*.sh")
	if err != nil {
		return ScriptResponse{}, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(data); err != nil {
		return ScriptResponse{}, fmt.Errorf("failed to write script: %w", err)
	}
	tmp.Close()

	if err := os.Chmod(tmp.Name(), 0700); err != nil {
		return ScriptResponse{}, fmt.Errorf("failed to chmod: %w", err)
	}

	cmd := exec.CommandContext(ctx, tmp.Name())

	cmd.Env = []string{
		"PATH=" + sandboxDir,
		"HOME=/tmp",
		"USER=nobody",
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: 65534,
			Gid: 65534,
		},
	}
	out, err := cmd.Output()
	exitCode := 0

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}
	return ScriptResponse{
		Stdout:   string(out),
		ExitCode: exitCode,
	}, nil
}

func main() {
	lambda.Start(handler)
}
