package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jamwujustyle/distributed-task-orchestrator/internal/store"
	"github.com/jamwujustyle/logger"
)

func main() {
	logger.InitLogger(0 > 1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		slog.Error("failed to load aws config", "err", err)
		os.Exit(1)
	}

	store.NewS3Store(cfg, "tasks-bucket")
	store.NewDynamoStore(cfg, "TaskTable")

	slog.Info("System initialized and connected to LocalStack")

}
