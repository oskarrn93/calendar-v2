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

	appConfig, err := config.Initialize()
	if err != nil {
		panic(err)
	}
	httpClient := resty.New()
	storage, err := awsutil.NewS3Storage(ctx, appConfig.S3Bucket, logger)
	if err != nil {
		panic(fmt.Errorf("failed to create S3 storage: %w", err))
	}

	rapidApi := rapidapi.New(httpClient, appConfig.RapidApi)

	handler := football.NewHandler(rapidApi, storage, logger)
	if err := handler.Handler(ctx); err != nil {
		logger.Error("Football handler failed", "error", err)
		panic(err)
	}

	logger.Info("Event successfully processed")
}
