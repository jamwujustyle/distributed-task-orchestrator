package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"

	"github.com/aws/aws-lambda-go/lambda"
)

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
