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

type Storage interface {
	Upload(ctx context.Context, s3Key string, data []byte) error
}

type S3Storage struct {
	s3Client *s3.Client
	s3Bucket string
	logger   *slog.Logger
}

func NewS3Storage(ctx context.Context, bucket string, logger *slog.Logger) (*S3Storage, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &S3Storage{
		s3Client: s3.NewFromConfig(cfg),
		s3Bucket: bucket,
		logger:   logger,
	}, nil
}

func (s *S3Storage) Upload(ctx context.Context, s3Key string, data []byte) error {
	response, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.s3Bucket),
		Key:         aws.String(s3Key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("text/calendar; charset=utf-8"),
		// Overrides CloudFront's one-day default TTL so the daily refresh reaches subscribers within the hour.
		CacheControl: aws.String("max-age=3600"),
	})
	s.logger.Debug("Upload S3 response", "s3Key", s3Key, "response", response)

	if err != nil {
		return fmt.Errorf("failed to upload to s3: %w", err)
	}

	s.logger.Info("Successfully uploaded file to S3", "s3Key", s3Key)

	return nil
}
