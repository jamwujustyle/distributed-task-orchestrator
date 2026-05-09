package main

import (
	"context"
	"log/slog"
	"os"

	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"github.com/jamwujustyle/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var target string = "controller:50051"

func main() {
	logger.InitLogger(0 > 1)

	ctx := context.Background()

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("failed to create a client to", "target", target)
		os.Exit(1)
	}
	defer conn.Close()

	c := pb.NewTaskServiceClient(conn)
	doPollTasks(ctx, c)
}
