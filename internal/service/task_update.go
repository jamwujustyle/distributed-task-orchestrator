package service

import (
	"context"

	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *TaskService) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*emptypb.Empty, error) {
	err := s.dynamo.UpdateTask(ctx, req.Id, req.Status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update task status: %s", err)
	}
	return &emptypb.Empty{}, nil
}
