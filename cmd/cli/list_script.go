package main

import (
	"context"
	"io"
	"log/slog"

	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func doListScripts(ctx context.Context, c pb.TaskServiceClient) {
	stream, err := c.ListScripts(ctx, &emptypb.Empty{})
	if err != nil {
		slog.Info("failed to list scripts")
		return
	}

	for {
		key, err := stream.Recv()
		if err == io.EOF {
			slog.Info("no scripts left")
			return
		}
		if err != nil {
			slog.Info("failed to stream scripts", "err", err)
			return
		}
		slog.Info("script", "key", key.GetKey())

	}
}
