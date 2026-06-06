package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/google/uuid"
	"github.com/jamwujustyle/distributed-task-orchestrator/cmd/worker/client"
	"github.com/jamwujustyle/distributed-task-orchestrator/cmd/worker/consumer"
	"github.com/jamwujustyle/distributed-task-orchestrator/cmd/worker/engine"
	"github.com/jamwujustyle/distributed-task-orchestrator/internal/queue"
	"github.com/jamwujustyle/distributed-task-orchestrator/internal/service"
	"github.com/jamwujustyle/distributed-task-orchestrator/internal/store"
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"testing"
	"time"
)

func TestTaskLifecycle_Integration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	eKey := "AWS_ENDPOINT_URL"
	awsUrl := os.Getenv(eKey)
	if awsUrl == "" {
		awsUrl = "http://localhost:4566"
	}
	os.Setenv(eKey, awsUrl)

	broker := os.Getenv("KAFKA_BROKERS")
	if broker == "" {
		broker = "localhost:9092"
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		t.Fatalf("Error loading AWS configuration: %v", err)
	}

	dynamodb := store.NewDynamoStore(cfg, "TaskTable")
	s3store := store.NewS3Store(cfg, "tasks-bucket")

	if err := dynamodb.Ping(ctx); err != nil {
		t.Fatalf("Error reaching dynamodb: %v", err)
	}
	if err := s3store.Ping(ctx); err != nil {
		t.Fatalf("Error reaching s3store: %v", err)
	}

	producer := queue.NewTaskProducer([]string{broker}, "tasks")
	defer producer.Close()

	svc := service.NewTaskService(dynamodb, s3store, producer)
	scriptContent := []byte("echo \"i say quack - i'm a duck\"\n")

	sKey, err := svc.UploadScript(ctx, &pb.Script{Content: scriptContent})
	if err != nil {
		t.Fatalf("Error uploading script: %v", err)
	}

	tId, err := svc.SubmitTask(ctx, &pb.Task{ScriptS3Key: sKey.GetKey()})
	if err != nil {
		t.Fatalf("Error submitting task: %v", err)
	}

	task, err := svc.RetrieveTask(ctx, &pb.TaskId{ID: tId.GetID()})
	if err != nil {
		t.Fatalf("Error retrieving task: %v", err)
	}
	if task.GetStatus() != pb.TaskStatus_PENDING {
		t.Fatalf("Task status mismatch. expected: %v, got: %v", pb.TaskStatus_PENDING, task.GetStatus())
	}

	cons := consumer.NewConsumer([]string{broker}, "tasks", "test-integration-group-"+uuid.NewString())
	defer cons.Close()

	msg, kTask, err := cons.FetchTask(ctx)
	if err != nil {
		t.Fatalf("Error fetching tasks: %v", err)
	}

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Error listening: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, svc)

	go func() {
		_ = s.Serve(lis)
	}()
	defer s.Stop()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Error establishing client connection: %v", err)
	}

	c := pb.NewTaskServiceClient(conn)
	e := engine.NewEngine(cfg, "task-executor")

	err = client.RunTaskLifecycle(ctx, c, e, kTask)
	if err != nil {
		t.Fatalf("Error in RunTaskLifecycle: %v", err)
	}

	if err := cons.Commit(ctx, msg); err != nil {
		t.Fatalf("Failed to commit offset to Kafka: %v", err)
	}

	fRecord, err := dynamodb.GetTask(ctx, kTask.GetId())
	if err != nil {
		t.Fatalf("Error retriving task: %v", err)
	}

	if fRecord.Status != pb.TaskStatus_COMPLETED {
		t.Fatalf("Status mismatch, expected: %v, got: %v", pb.TaskStatus_COMPLETED, fRecord.Status)
	}

	expected := "i say quack - i'm a duck"
	if fRecord.Result != string(expected) {
		t.Fatalf("Result mismatch, expected: %v, got %v", expected, fRecord.Result)
	}
}
