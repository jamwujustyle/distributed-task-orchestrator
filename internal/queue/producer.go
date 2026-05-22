package queue

import (
	"context"
	"fmt"
	"time"

	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"github.com/jamwujustyle/distributed-task-orchestrator/pkg/telemetry"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type TaskProducer struct {
	Writer *kafka.Writer
}

func NewTaskProducer(brokers []string, topic string) *TaskProducer {
	return &TaskProducer{
		Writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			Async:        true,
			RequiredAcks: kafka.RequireAll,
			MaxAttempts:  5,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (p *TaskProducer) PublishTask(ctx context.Context, task *pb.Task) error {
	key := []byte(task.GetId())

	value, err := proto.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task to protobuf: %w", err)
	}
	msg := kafka.Message{
		Key:   key,
		Value: value,
	}

	telemetry.InjectTraceContext(ctx, &msg)

	return p.Writer.WriteMessages(ctx, msg)
}

func (p *TaskProducer) Close() error { return p.Writer.Close() }
