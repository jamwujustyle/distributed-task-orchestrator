package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jamwujustyle/distributed-task-orchestrator/internal/service"
	"github.com/jamwujustyle/distributed-task-orchestrator/internal/store"
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"github.com/jamwujustyle/logger"
	"google.golang.org/grpc"
)

var addr string = "0.0.0.0:50051"

func main() {
	logger.InitLogger(0 > 1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		slog.Error("failed to load aws config", "err", err)
		os.Exit(1)
	}

	dbStore := store.NewDynamoStore(cfg, "TaskTable")
	if err := dbStore.Ping(ctx); err != nil {
		slog.Error("DynamoDB storage offline", "err", err)
		os.Exit(1)
	}
	s3Store := store.NewS3Store(cfg, "tasks-bucket")
	if err := s3Store.Ping(ctx); err != nil {
		slog.Error("S3 storage offline", "err", err)
		os.Exit(1)
	}

	svc := service.NewTaskService(dbStore, s3Store)

	slog.Info("Storage layer verified", "provider", "LocalStack", "server", "controller")

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("could not listen", "err", err)
		os.Exit(1)
	}
	slog.Info("gRPC listening", "addr", addr)

	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, svc)

	if err := s.Serve(lis); err != nil {
		slog.Error("gRPC Failed to serve", "err", err)
		os.Exit(1)
	}

}
