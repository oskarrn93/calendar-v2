package awsutil

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func S3Client() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	return client, nil
}

type Storage interface {
	Upload(ctx context.Context, filename string, data []byte, logger *slog.Logger) error
}

type S3Storage struct {
	S3Client *s3.Client
	S3Bucket string
}

func (s *S3Storage) Upload(ctx context.Context, s3Key string, data []byte, logger *slog.Logger) error {
	response, err := s.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.S3Bucket),
		Key:    aws.String(s3Key),
		Body:   bytes.NewReader(data),
	})
	logger.Debug("Upload S3 response", "s3Key", s3Key, "response", response)

	if err != nil {
		return fmt.Errorf("failed to upload to s3: %w", err)
	}

	logger.Info("Successfully uploaded file to S3", "s3Key", s3Key)

	return nil
}
