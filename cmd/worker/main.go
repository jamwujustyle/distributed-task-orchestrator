package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jamwujustyle/distributed-task-orchestrator/cmd/worker/client"
	"github.com/jamwujustyle/distributed-task-orchestrator/cmd/worker/consumer"
	"github.com/jamwujustyle/distributed-task-orchestrator/cmd/worker/engine"
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"github.com/jamwujustyle/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TaskEnvelope struct {
	Msg  kafka.Message
	Task *pb.Task
}

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

	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "kafka:9092"
	}
	cons := consumer.NewConsumer([]string{brokers}, "tasks", "worker-group")
	defer cons.Close()

	tasksChan := make(chan TaskEnvelope, 100)

	for i := range 5 {
		go worker(ctx, i, c, e, cons, tasksChan)
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		slog.Info("metrics server listening", "addr", ":9091")
		if err := http.ListenAndServe(":9091", nil); err != nil {
			slog.Error("metrics server failed", "err", err)
		}
	}()

	for {
		msg, task, err := cons.FetchTask(ctx)
		if err != nil {
			slog.Error("fetch failed", "err", err)
			time.Sleep(5 * time.Second)
			continue
		}

		tasksChan <- TaskEnvelope{
			Msg:  msg,
			Task: task,
		}
	}

}

func worker(ctx context.Context, id int, c pb.TaskServiceClient, e *engine.Engine, cons *consumer.TaskConsumer, tasks <-chan TaskEnvelope) {
	for env := range tasks {
		err := client.RunTaskLifecycle(ctx, c, e, env.Task)
		if err == nil {
			cons.Commit(ctx, env.Msg)
			slog.Info("task commited", "worker", id, "task", env.Task.GetId())
		} else {
			slog.Error("task failed", "worker", id, "task", env.Task.GetId())
		}
	}
}
