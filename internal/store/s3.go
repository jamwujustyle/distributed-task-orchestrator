package store

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Store struct {
	Client *s3.Client
	Bucket string
}

func NewS3Store(cfg aws.Config, bucketName string) *S3Store {
	return &S3Store{
		Client: s3.NewFromConfig(cfg),
		Bucket: bucketName,
	}
}

func (s *S3Store) ListObjects(ctx context.Context) (*s3.ListObjectsV2Output, error) {
	output, err := s.Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.Bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to list objects: %w", err)
	}
	return output, nil
}

func (s *S3Store) UploadTaskScript(ctx context.Context, key string, content []byte) error {
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(content),
	})
	if err != nil {
		return fmt.Errorf("failed to upload script to S3 : %w", err)
	}
	return nil
}

func (s *S3Store) Ping(ctx context.Context) error {
	_, err := s.Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.Bucket),
	})
	return err
}
