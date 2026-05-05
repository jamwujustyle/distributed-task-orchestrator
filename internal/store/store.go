package store

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type TaskStore interface {
	SaveTask(ctx context.Context, task Task) error
	GetTask(ctx context.Context, id string) (*Task, error)
	UpdateTask(ctx context.Context, id string, status TaskStatus) error
	GetPendingTasks(ctx context.Context) ([]Task, error)
}
type ArtifactStore interface {
	UploadScript(ctx context.Context, key string, data []byte) error
	GetScript(ctx context.Context, key string) ([]byte, error)
	ListScripts(ctx context.Context) (*s3.ListObjectsV2Output, error)
}
