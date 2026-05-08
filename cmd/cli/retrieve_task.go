package main

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
)

func doRetrieveTask(ctx context.Context, c pb.TaskServiceClient, id string) error {
	t, err := c.RetrieveTask(ctx, &pb.TaskId{
		ID: id,
	})
	if err != nil {
		return fmt.Errorf("failed to retrieve task: %w", err)
	}
	slog.Info("task retrieved", "task", t.GetId(), "status", t.GetStatus(), "script_key", t.GetScriptS3Key())
	return nil
}
