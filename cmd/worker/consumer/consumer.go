package consumer

import (
	"context"
	"net"
	"strings"

	pb "github.com/jamwujustyle/distributed-task-orchestrator/pkg/protocol/v1"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type TaskConsumer struct {
	consumer *kafka.Reader
}

func NewConsumer(brokers []string, topic string, groupID string) *TaskConsumer {
	var dialer *kafka.Dialer

	if len(brokers) > 0 && (strings.Contains(brokers[0], "localhost") || strings.Contains(brokers[0], "127.0.0.1")) {
		dialer = &kafka.Dialer{
			DualStack: true,
			DialFunc: func(ctx context.Context, network, address string) (net.Conn, error) {
				if strings.Contains(address, "kafka-controller") {
					address = "127.0.0.1:9092"
				}
				return (&net.Dialer{}).DialContext(ctx, network, address)
			},
		}
	}

	return &TaskConsumer{
		consumer: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3,
			MaxBytes: 10e6,
			Dialer:   dialer,
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
