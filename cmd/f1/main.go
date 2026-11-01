package main

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/oskarrn93/calendar-v2/internal/awsutil"
	"github.com/oskarrn93/calendar-v2/internal/config"
	"github.com/oskarrn93/calendar-v2/internal/f1"
	"github.com/oskarrn93/calendar-v2/internal/logging"
)

func main() {
	ctx := context.Background()

	logger := logging.New()

	logger.Info("F1 command")

	appConfig, err := config.Initialize()
	if err != nil {
		panic(err)
	}

	storage, err := awsutil.NewS3Storage(ctx, appConfig.S3Bucket, logger)
	if err != nil {
		panic(fmt.Errorf("failed to create S3 storage: %w", err))
	}

	handler := f1.NewHandler(resty.New(), appConfig.F1, storage, logger)
	if err := handler.Handler(ctx); err != nil {
		panic(fmt.Errorf("F1 handler failed: %w", err))
	}

	logger.Info("Event successfully processed")
}
