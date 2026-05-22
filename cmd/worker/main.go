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
	"github.com/jamwujustyle/distributed-task-orchestrator/pkg/telemetry"
	"github.com/jamwujustyle/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TaskEnvelope struct {
	Msg  kafka.Message
	Task *pb.Task
	Ctx  context.Context
}

var target string = "controller:50051"

func main() {
	logger.InitLogger(0 > 1)

	ctx := context.Background()

	shutdown, err := telemetry.InitTracer(ctx, "worker", "jaeger:4317")
	if err != nil {
		slog.Error("failed to init tracer", "err", err)
	} else {
		defer shutdown(ctx)
	}
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		slog.Error("failed to load default config")
		os.Exit(1)
	}

	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
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
		msgCtx := telemetry.ExtractTraceContext(ctx, &msg)
		tasksChan <- TaskEnvelope{
			Msg:  msg,
			Task: task,
			Ctx:  msgCtx,
		}
	}

}

func worker(ctx context.Context, id int, c pb.TaskServiceClient, e *engine.Engine, cons *consumer.TaskConsumer, tasks <-chan TaskEnvelope) {
	for env := range tasks {

		workerCtx, span := telemetry.GetTracer().Start(env.Ctx, "WorkerProcessTask")

		err := client.RunTaskLifecycle(workerCtx, c, e, env.Task)
		if err == nil {
			cons.Commit(ctx, env.Msg)
			slog.Info("task commited", "worker", id, "task", env.Task.GetId())
			span.SetStatus(codes.Ok, "task completed successfully")
		} else {
			slog.Error("task failed", "worker", id, "task", env.Task.GetId())
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}
}
