package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-resty/resty/v2"
)

func handler(ctx context.Context, event json.RawMessage) error {
	logger := NewLogger()

	logger.Info("Received Event", "event", event)

	appConfig := InitializeConfig(logger)
	httpClient := resty.New()
	s3Client, err := getS3Client()
	if err != nil {
		return fmt.Errorf("Failed to get S3 client: %w", err)
	}

	rapidApi := RapidApi{httpClient: httpClient, config: appConfig.rapidApi}
	storage := S3Storage{
		s3Client: s3Client,
		s3Bucket: appConfig.s3Bucket,
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
