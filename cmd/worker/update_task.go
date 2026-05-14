package main

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
)

func doUpdateTask(ctx context.Context, c pb.TaskServiceClient, id string, status pb.TaskStatus, result string) error {
	_, err := c.UpdateTask(ctx, &pb.UpdateTaskRequest{
		Id:     id,
		Status: status,
		Result: result,
	})
	if err != nil {
		return fmt.Errorf("server returned an error: %w", err)

	}
	slog.Info("task updated successfully", "id", id, "status", status)

	return nil
}
