package service

import (
	"context"

	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *TaskService) RetrieveScript(ctx context.Context, req *pb.ScriptKey) (*pb.Script, error) {
	script, err := s.s3.GetScript(ctx, req.Key)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve a script: %v", err)
	}

	return &pb.Script{
		Content: script,
	}, nil

}
