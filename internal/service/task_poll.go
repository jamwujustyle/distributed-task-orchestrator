package service

import (
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *TaskService) PollTasks(_ *emptypb.Empty, stream pb.TaskService_PollTasksServer) error {
	tasks, err := s.dynamo.GetPendingTasks(stream.Context())
	if err != nil {
		return status.Errorf(codes.Internal, "failed to get pending tasks: %v", err)
	}

	for _, t := range tasks {
		err := stream.Send(&pb.Task{
			Id:          t.ID,
			ScriptS3Key: t.ScriptS3Key,
			Status:      t.Status,
		})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to send task: %v", err)
		}
	}
	return nil
}
