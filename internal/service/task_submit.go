package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jamwujustyle/distributed-task-orchestrator/internal/store"
	"github.com/jamwujustyle/distributed-task-orchestrator/pkg/metrics"
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *TaskService) SubmitTask(ctx context.Context, req *pb.Task) (*pb.TaskId, error) {
	generatedID := uuid.New().String()
	t := store.Task{
		ID:          generatedID,
		ScriptS3Key: req.ScriptS3Key,
		Status:      pb.TaskStatus_PENDING,
	}

	if err := s.dynamo.SaveTask(ctx, t); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to save task: %v", err)
	}

	req.Id = generatedID

	err := s.producer.PublishTask(ctx, req)
	if err != nil {
		metrics.TasksPublishErrors.Inc()
		return nil, status.Errorf(codes.Internal, "failed to publish task: %v", err)
	}
	metrics.TasksSubmitted.Inc()
	return &pb.TaskId{ID: t.ID}, nil
}
