package service

import (
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *TaskService) ListScripts(_ *emptypb.Empty, stream pb.TaskService_ListScriptsServer) error {
	scripts, err := s.s3.ListScripts(stream.Context())
	if err != nil {
		return status.Errorf(codes.Internal, "failed to list scripts: %v", err)
	}

	for _, key := range scripts {
		if err := stream.Send(&pb.ScriptKey{
			Key: key,
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to send over stream: %v", err)
		}
	}
	return nil
}
