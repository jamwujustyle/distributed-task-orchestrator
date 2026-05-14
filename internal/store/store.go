package store

import (
	"context"

	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
)

type TaskStore interface {
	SaveTask(ctx context.Context, task Task) error
	GetTask(ctx context.Context, id string) (*Task, error)
	UpdateTask(ctx context.Context, id string, status pb.TaskStatus, result string) error
	GetPendingTasks(ctx context.Context) ([]Task, error)
}
type ArtifactStore interface {
	UploadScript(ctx context.Context, key string, data []byte) error
	GetScript(ctx context.Context, key string) ([]byte, error)
	ListScripts(ctx context.Context) ([]string, error)
}
