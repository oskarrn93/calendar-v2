package main

import (
	"context"
	"encoding/json"
	"fmt"

	"oskarrn93/calendar-v2/internal/config"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-resty/resty/v2"
)

func handler(ctx context.Context, event json.RawMessage) error {
	logger := NewLogger()

	logger.Info("Received Event", "event", event)

	appConfig := config.Initialize(logger)
	httpClient := resty.New()
	s3Client, err := getS3Client()
	if err != nil {
		return fmt.Errorf("Failed to get S3 client: %w", err)
	}

	rapidApi := RapidApi{httpClient: httpClient, config: appConfig.RapidApi}
	storage := S3Storage{
		s3Client: s3Client,
		s3Bucket: appConfig.S3Bucket,
		logger:   logger,
	}

	nbaHandler := NbaHandler{rapidApi: rapidApi, storage: &storage, logger: logger}
	if err := nbaHandler.handler(ctx); err != nil {
		logger.Error("NBA handler failed", "error", err)
		return err
	}

	logger.Info("Event successfully processed")
	return nil
}

func main() {
	lambda.Start(handler)
}
