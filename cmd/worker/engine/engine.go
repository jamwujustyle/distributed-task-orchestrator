package engine

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type ScriptRequest struct {
	Script string `json:"script"`
}
type ScriptResponse struct {
	Stdout   string `json:"stdout"`
	ExitCode int    `json:"exit_code"`
}

type Engine struct {
	lambda *lambda.Client
	fn     string
}

func NewEngine(cfg aws.Config, fnName string) *Engine {
	return &Engine{
		lambda: lambda.NewFromConfig(cfg, func(o *lambda.Options) {
			o.BaseEndpoint = aws.String("http://localstack:4566")
		}),
		fn: fnName,
	}
}

func (e *Engine) ExecuteScript(ctx context.Context, script []byte) (string, error) {
	payload, err := json.Marshal(ScriptRequest{
		Script: base64.StdEncoding.EncodeToString(script),
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}
	r, err := e.lambda.Invoke(ctx, &lambda.InvokeInput{
		FunctionName: aws.String(e.fn),
		Payload:      payload,
	})
	if err != nil {
		return "", fmt.Errorf("failed to invoke lambda: %w", err)
	}

	if r.FunctionError != nil {
		return "", fmt.Errorf("labmda runtime error: %s", *r.FunctionError)
	}

	var res ScriptResponse

	if err := json.Unmarshal(r.Payload, &res); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if res.ExitCode != 0 {
		return "", fmt.Errorf("script failed with exit code: %d, %s", res.ExitCode, res.Stdout)
	}

	slog.Info(res.Stdout)
	return strings.TrimSpace(res.Stdout), nil
}
