package main

import (
	"context"
	"io"
	"log/slog"

	"github.com/jamwujustyle/distributed-task-orchestrator/internal/engine"
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func doPollTasks(ctx context.Context, c pb.TaskServiceClient) {
	stream, err := c.PollTasks(ctx, &emptypb.Empty{})
	if err != nil {
		slog.Error("failed to poll tasks")
		return
	}

	for {
		t, err := stream.Recv()
		if err == io.EOF {
			slog.Info("no more pending tasks")
			return
		}
		if err != nil {
			slog.Error("failed to stream")
			return
		}

		slog.Info("received tasks", "id", t.GetId())

		if err := doUpdateTask(ctx, c, t.GetId(), pb.TaskStatus_RUNNING); err != nil {
			continue
		}

		script, err := doRetrieveScript(ctx, c, t.ScriptS3Key)
		if err != nil {
			slog.Error("error pulling script from s3", "task", t.GetId(), "key", t.GetScriptS3Key())
			doUpdateTask(ctx, c, t.GetId(), pb.TaskStatus_FAILED)
			continue
		}

		if err := engine.ExecuteScript(script); err != nil {
			slog.Error("script execution failed", "err", err)
			doUpdateTask(ctx, c, t.GetId(), pb.TaskStatus_FAILED)
			continue
		}
		doUpdateTask(ctx, c, t.GetId(), pb.TaskStatus_COMPLETED)
	}
}
