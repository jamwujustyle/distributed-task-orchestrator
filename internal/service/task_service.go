package service

import (
	"github.com/jamwujustyle/distributed-task-orchestrator/internal/queue"
	"github.com/jamwujustyle/distributed-task-orchestrator/internal/store"
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
)

type TaskService struct {
	pb.UnimplementedTaskServiceServer
	dynamo   store.TaskStore
	s3       store.ArtifactStore
	producer *queue.TaskProducer
}

func NewTaskService(dynamo store.TaskStore, s3 store.ArtifactStore, producer *queue.TaskProducer) *TaskService {
	return &TaskService{dynamo: dynamo, s3: s3, producer: producer}
}
