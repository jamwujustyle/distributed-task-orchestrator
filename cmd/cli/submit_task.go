package main

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
)

func doSubmitTask(ctx context.Context, c pb.TaskServiceClient, key string) error {
	t, err := c.SubmitTask(ctx, &pb.Task{
		ScriptS3Key: key,
	})
	if err != nil {
		return fmt.Errorf("failed to submit task: %w", err)
	}
	slog.Info("task submitted successfully", "key", t.GetID())
	return nil
}
