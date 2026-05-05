package store

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type TaskStatus string

const (
	StatusPending   TaskStatus = "PENDING"
	StatusRunning   TaskStatus = "RUNNING"
	StatusCompleted TaskStatus = "COMPLETED"
	StatusFailed    TaskStatus = "FAILED"
)

type Task struct {
	ID          string     `dynamodbav:"ID"`
	ScriptS3Key string     `dynamodbav:"ScriptS3Key"`
	Status      TaskStatus `dynamodbav:"Status"`
	CreatedAt   int64      `dynamodbav:"CreatedAt"`
	UpdatedAt   int64      `dynamodbav:"UpdatedAt"`
}

type DynamoStore struct {
	Client *dynamodb.Client
	Table  string
}

func NewDynamoStore(cfg aws.Config, tableName string) *DynamoStore {
	return &DynamoStore{
		Client: dynamodb.NewFromConfig(cfg),
		Table:  tableName,
	}
}

func (d *DynamoStore) SaveTask(ctx context.Context, task Task) error {
	now := time.Now().Unix()
	task.CreatedAt = now
	task.UpdatedAt = now

	item, err := attributevalue.MarshalMap(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.Table),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put item in dynamo: %w", err)
	}

	return nil
}

func (d *DynamoStore) GetTask(ctx context.Context, id string) (*Task, error) {
	k, err := attributevalue.MarshalMap(map[string]string{"ID": id})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key: %w", err)
	}
	r, err := d.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.Table),
		Key:       k,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}
	if r.Item == nil {
		return nil, fmt.Errorf("task not found: %s", id)
	}

	var t Task
	if err := attributevalue.UnmarshalMap(r.Item, &t); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return &t, nil
}

func (d *DynamoStore) UpdateTask(ctx context.Context, id string, status TaskStatus) error {
	key, err := attributevalue.MarshalMap(map[string]string{"ID": id})
	if err != nil {
		return fmt.Errorf("failed to marshal key: %w", err)
	}

	values, err := attributevalue.MarshalMap(map[string]any{
		":status":    status,
		":updatedAt": time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal status || updatedAt: %w", err)
	}

	_, err = d.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        aws.String(d.Table),
		Key:              key,
		UpdateExpression: aws.String("SET #status = :status, UpdatedAt = :updatedAt"),
		ExpressionAttributeNames: map[string]string{
			"#status": "Status",
		},
		ExpressionAttributeValues: values,
	})
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	return nil
}

func (d *DynamoStore) GetPendingTasks(ctx context.Context) ([]Task, error) {
	vs, err := attributevalue.MarshalMap(map[string]any{":status": StatusPending})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal status: %w", err)
	}

	r, err := d.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName:                 aws.String(d.Table),
		FilterExpression:          aws.String("#status = :status"),
		ExpressionAttributeNames:  map[string]string{"#status": "Status"},
		ExpressionAttributeValues: vs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan tasks: %w", err)
	}
	var ts []Task

	if err = attributevalue.UnmarshalListOfMaps(r.Items, &ts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
	}

	return ts, nil
}

func (d *DynamoStore) Ping(ctx context.Context) error {
	_, err := d.Client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(d.Table),
	})
	return err
}
