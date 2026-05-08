package service

import (
	"github.com/jamwujustyle/distributed-task-orchestrator/internal/store"
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
)

type TaskService struct {
	pb.UnimplementedTaskServiceServer
	dynamo store.TaskStore
	s3     store.ArtifactStore
}

func NewTaskService(dynamo store.TaskStore, s3 store.ArtifactStore) *TaskService {
	return &TaskService{dynamo: dynamo, s3: s3}
}
