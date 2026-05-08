package main

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
)

func doRetrieveScript(ctx context.Context, c pb.TaskServiceClient, key string) error {
	s, err := c.RetrieveScript(ctx, &pb.ScriptKey{
		Key: key,
	})
	if err != nil {
		return fmt.Errorf("failed to retrieve script, %w", err)
	}
	slog.Info("retrieved script successfully", "content", s.GetContent())
	return nil
}
