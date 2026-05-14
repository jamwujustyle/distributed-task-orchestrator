package main

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/jamwujustyle/distributed-task-orchestrator/internal/engine"
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func doPollTasks(ctx context.Context, c pb.TaskServiceClient, e *engine.Engine) {
	for {
		stream, err := c.PollTasks(ctx, &emptypb.Empty{})
		if err != nil {
			slog.Error("failed to poll tasks, retrying", "err", err)
			time.Sleep(5 * time.Second)
			continue
		}
		for {
			t, err := stream.Recv()
			if err == io.EOF {
				slog.Info("no more pending tasks")
				time.Sleep(5 * time.Second)
				break
			}
			if err != nil {
				slog.Error("failed to stream, retrying poll", "err", err)
				break
			}

			slog.Info("received task", "id", t.GetId())

			if err := doUpdateTask(ctx, c, t.GetId(), pb.TaskStatus_RUNNING, t.GetResult()); err != nil {
				slog.Error("failed to update task to running", "err", err)
				continue
			}

			script, err := doRetrieveScript(ctx, c, t.ScriptS3Key)
			if err != nil {
				slog.Error("error pulling script from s3", "task", t.GetId(), "key", t.GetScriptS3Key(), "err", err)
				doUpdateTask(ctx, c, t.GetId(), pb.TaskStatus_FAILED, t.GetResult())
				continue
			}

			r, err := e.ExecuteScript(ctx, script)
			if err != nil {
				slog.Error("script execution failed", "err", err)
				doUpdateTask(ctx, c, t.GetId(), pb.TaskStatus_FAILED, r)
				continue
			}
			doUpdateTask(ctx, c, t.GetId(), pb.TaskStatus_COMPLETED, r)
		}
	}
}
