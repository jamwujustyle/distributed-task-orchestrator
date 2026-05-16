package consumer

import (
	"context"
	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type TaskConsumer struct {
	consumer *kafka.Reader
}

func NewConsumer(brokers []string, topic string, groupID string) *TaskConsumer {
	return &TaskConsumer{
		consumer: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
	}
}

func (c *TaskConsumer) FetchTask(ctx context.Context) (kafka.Message, *pb.Task, error) {
	msg, err := c.consumer.FetchMessage(ctx)
	if err != nil {
		return kafka.Message{}, nil, err
	}
	var task pb.Task

	if err := proto.Unmarshal(msg.Value, &task); err != nil {
		return msg, nil, err
	}
	return msg, &task, nil

}

func (c *TaskConsumer) Commit(ctx context.Context, msg kafka.Message) error {
	return c.consumer.CommitMessages(ctx, msg)
}

func (c *TaskConsumer) Close() error {
	return c.consumer.Close()
}
