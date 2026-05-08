package main

import (
	"context"
	"log/slog"
	"os"
	"time"

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

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("failed to establish connection", "err", err)
	}

	c := pb.NewTaskServiceClient(conn)

	if len(os.Args) < 2 {
		slog.Error("usage: cli <command> [args]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "upload":
		data, err := os.ReadFile(os.Args[2])
		if err != nil {
			slog.Error("failed to read file", "err", err)
			os.Exit(1)
		}
		key, err := doUploadScript(ctx, c, data)
		if err != nil {
			slog.Error("failed 'doUploadScript' job", "err", err)
			//i wonder if the os.exit radicality justified, perhaps continue is a better fit?
			os.Exit(1)
		}
		slog.Info("use this key to to submit a task", "key", key)
	case "list":
		doListScripts(ctx, c)
	case "submit":
		doSubmitTask(ctx, c, os.Args[2])
	case "retrieve":
		doRetrieveTask(ctx, c, os.Args[2])
	}

}
