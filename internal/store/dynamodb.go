package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Store struct {
	Client *s3.Client
	Bucket string
}

func NewS3Store(bucketName string) (*S3Store, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		return nil, fmt.Errorf("Unable to load S3 bucket: %w", err)
	}
	client := s3.NewFromConfig(cfg)

	return &S3Store{
		Client: client,
		Bucket: bucketName,
	}, nil
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
