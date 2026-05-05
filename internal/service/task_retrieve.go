package service

import (
	"context"

	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *TaskService) RetrieveTask(ctx context.Context, req *pb.TaskId) (*pb.Task, error) {
	t, err := s.dynamo.GetTask(ctx, req.ID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "task not found: %s", err)
	}

	return &pb.Task{
		Id:          t.ID,
		ScriptS3Key: t.ScriptS3Key,
		Status:      string(t.Status),
	}, nil
}
