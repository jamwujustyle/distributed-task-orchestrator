package service

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *TaskService) UploadScript(ctx context.Context, req *pb.Script) (*pb.ScriptKey, error) {
	key := uuid.New().String()

	if err := s.s3.UploadScript(ctx, key, req.Content); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload the script: %v", err)
	}

	return &pb.ScriptKey{
		Key: key,
	}, nil
}
