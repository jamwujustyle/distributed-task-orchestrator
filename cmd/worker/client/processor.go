package client

import (
	"context"
	"log/slog"
	"time"

	"github.com/jamwujustyle/distributed-task-orchestrator/cmd/worker/engine"
	"github.com/jamwujustyle/distributed-task-orchestrator/pkg/metrics"
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
)

func RunTaskLifecycle(ctx context.Context, c pb.TaskServiceClient, e *engine.Engine, task *pb.Task) error {
	start := time.Now()
	defer func() {
		metrics.TaskDuration.Observe(time.Since(start).Seconds())
	}()

	slog.Info("running task lifecycle", "taskId", task.GetId(), "scriptKey", task.GetScriptS3Key())

	if err := doUpdateTask(ctx, c, task.GetId(), pb.TaskStatus_RUNNING, ""); err != nil {
		metrics.TasksProcessed.WithLabelValues("failed").Inc()
		return err
	}

	script, err := doRetrieveScript(ctx, c, task.GetScriptS3Key())
	if err != nil {
		doUpdateTask(ctx, c, task.GetId(), pb.TaskStatus_FAILED, err.Error())
		metrics.TasksProcessed.WithLabelValues("failed").Inc()
		return err
	}
	res, err := e.ExecuteScript(ctx, script)
	if err != nil {
		doUpdateTask(ctx, c, task.GetId(), pb.TaskStatus_FAILED, err.Error())
		metrics.TasksProcessed.WithLabelValues("failed").Inc()
		return err
	}
	metrics.TasksProcessed.WithLabelValues("completed").Inc()

	return doUpdateTask(ctx, c, task.GetId(), pb.TaskStatus_COMPLETED, res)
}
