package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jamwujustyle/distributed-task-orchestrator/internal/store"
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"github.com/jamwujustyle/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var target string = "0.0.0.0:50051"

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
	}

	s3Store := store.NewS3Store(cfg, "tasks-bucket")
	if err := s3Store.Ping(ctx); err != nil {
		slog.Error("S3 storage offline", "err", err)
	}

	slog.Info("Storage layer verified", "provider", "LocalStack", "server", "worker")

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("failed to create a client to", "target", target)
		os.Exit(1)
	}
	defer conn.Close()
	c := pb.NewTaskServiceClient(conn)

}
