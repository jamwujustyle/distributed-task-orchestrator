package store

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Store struct {
	Client *s3.Client
	Bucket string
}

func NewS3Store(cfg aws.Config, bucketName string) *S3Store {
	return &S3Store{
		Client: s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UsePathStyle = true
			if ep := os.Getenv("AWS_ENDPOINT_URL"); ep != "" {
				o.BaseEndpoint = aws.String(ep)
			}
		}),
		Bucket: bucketName,
	}
}
func (s *S3Store) UploadScript(ctx context.Context, key string, content []byte) error {
	slog.Info("uploading script", "bucket", s.Bucket, "key", key, "size", len(content))
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Body:   bytes.NewReader(content),
		Key:    aws.String(key),
	})
	if err != nil {
		slog.Error("upload failed", "error", err)
		return fmt.Errorf("failed to upload script to S3 : %w", err)
	}
	slog.Info("upload succeeded", "bucket", s.Bucket, "key", key)
	return nil
}

func (s *S3Store) GetScript(ctx context.Context, key string) ([]byte, error) {
	slog.Info("retrieving script", "bucket", s.Bucket, "key", key)
	r, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		slog.Error("retrieve failed", "error", err)
		return nil, fmt.Errorf("failed to retrieve an object: %w", err)
	}
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("read failed", "error", err)
		return nil, fmt.Errorf("failed to read object body: %w", err)
	}
	slog.Info("retrieve succeeded", "size", len(data))
	return data, err
}

func (s *S3Store) ListScripts(ctx context.Context) ([]string, error) {
	output, err := s.Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.Bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	keys := make([]string, 0, len(output.Contents))
	for _, obj := range output.Contents {
		if obj.Key != nil {
			keys = append(keys, *obj.Key)
		}
	}
	return keys, nil
}

func (s *S3Store) Ping(ctx context.Context) error {
	_, err := s.Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.Bucket),
	})
	return err
}
