package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jamwujustyle/distributed-task-orchestrator/cmd/worker/client"
	"github.com/jamwujustyle/distributed-task-orchestrator/cmd/worker/consumer"
	"github.com/jamwujustyle/distributed-task-orchestrator/cmd/worker/engine"
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"github.com/jamwujustyle/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var target string = "controller:50051"

func main() {
	logger.InitLogger(0 > 1)

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		slog.Error("failed to load default config")
		os.Exit(1)
	}

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("failed to create a client to", "target", target)
		os.Exit(1)
	}
	defer conn.Close()

	e := engine.NewEngine(cfg, "task-executor")
	c := pb.NewTaskServiceClient(conn)
	cons := consumer.NewConsumer([]string{"kafka:29092"}, "tasks", "worker-group")
	defer cons.Close()

	for {
		msg, task, err := cons.FetchTask(ctx)
		if err != nil {
			slog.Error("failed to fetch tasks", "err", err)
			continue
		}
		slog.Info("received task from kafka", "id", task.GetId())

		err = client.RunTaskLifecycle(ctx, c, e, task)
		if err == nil {
			cons.Commit(ctx, msg)
			slog.Info("task finished and commited", "id", task.GetId())
		} else {
			slog.Error("task execution failed", "id", task.GetId())
		}
	}

}
