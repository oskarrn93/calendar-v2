package main

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/oskarrn93/calendar-v2/internal/awsutil"
	"github.com/oskarrn93/calendar-v2/internal/config"
	"github.com/oskarrn93/calendar-v2/internal/football"
	"github.com/oskarrn93/calendar-v2/internal/logging"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
)

func main() {
	ctx := context.Background()

	logger := logging.New()

	logger.Info("Football command")

	appConfig := config.Initialize(logger)
	httpClient := resty.New()
	s3Client, err := awsutil.S3Client()
	if err != nil {
		panic(fmt.Errorf("Failed to get S3 client: %w", err))
	}

	rapidApi := rapidapi.New(httpClient, appConfig.RapidApi)
	storage := awsutil.S3Storage{
		S3Client: s3Client,
		S3Bucket: appConfig.S3Bucket,
	}

	handler := football.NewHandler(rapidApi, &storage, logger)
	if err := handler.Handler(ctx); err != nil {
		logger.Error("Football handler failed", "error", err)
		panic(err)
	}

	logger.Info("Event successfully processed")
}
