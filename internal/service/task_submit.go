package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jamwujustyle/distributed-task-orchestrator/internal/store"
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *TaskService) SubmitTask(ctx context.Context, req *pb.Task) (*pb.TaskId, error) {
	t := store.Task{
		ID:          uuid.New().String(),
		ScriptS3Key: req.ScriptS3Key,
		Status:      store.StatusPending,
	}

	if err := s.dynamo.SaveTask(ctx, t); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to save task: %v", err)
	}

	return &pb.TaskId{ID: t.ID}, nil
}
