package client

import (
	"context"

	"github.com/jamwujustyle/distributed-task-orchestrator/cmd/worker/engine"
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
)

func RunTaskLifecycle(ctx context.Context, c pb.TaskServiceClient, e *engine.Engine, task *pb.Task) error {
	if err := doUpdateTask(ctx, c, task.GetId(), pb.TaskStatus_RUNNING, ""); err != nil {
		return err
	}

	script, err := doRetrieveScript(ctx, c, task.GetScriptS3Key())
	if err != nil {
		doUpdateTask(ctx, c, task.GetId(), pb.TaskStatus_FAILED, err.Error())
		return err
	}
	res, err := e.ExecuteScript(ctx, script)
	if err != nil {
		doUpdateTask(ctx, c, task.GetId(), pb.TaskStatus_FAILED, err.Error())
		return err
	}

	return doUpdateTask(ctx, c, task.GetId(), pb.TaskStatus_COMPLETED, res)
}
