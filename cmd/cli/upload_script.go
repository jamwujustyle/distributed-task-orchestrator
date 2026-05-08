package main

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
)

func doUploadScript(ctx context.Context, c pb.TaskServiceClient, data []byte) (string, error) {
	k, err := c.UploadScript(ctx, &pb.Script{
		Content: data,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload script: %w", err)
	}
	slog.Info("script uploaded successfully", "key", k.Key)
	return k.Key, nil
}
