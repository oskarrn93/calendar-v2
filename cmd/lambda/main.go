package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/oskarrn93/calendar-v2/internal/awsutil"
	"github.com/oskarrn93/calendar-v2/internal/config"
	"github.com/oskarrn93/calendar-v2/internal/logging"
	"github.com/oskarrn93/calendar-v2/internal/nba"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
	"github.com/oskarrn93/calendar-v2/internal/util"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-resty/resty/v2"
)

func handler(ctx context.Context, event json.RawMessage) error {
	logger := logging.New()

	logger.Info("Received Event", "event", event)

	appConfig := config.Initialize(logger)
	httpClient := resty.New()
	s3Client, err := awsutil.S3Client()
	if err != nil {
		return fmt.Errorf("Failed to get S3 client: %w", err)
	}

	rapidApi := rapidapi.New(httpClient, appConfig.RapidApi)
	storage := awsutil.S3Storage{
		S3Client: s3Client,
		S3Bucket: appConfig.S3Bucket,
	}

	// use any for results since we don't care about them
	wg := util.NewWaitGroup[any](ctx)

	wg.Run(func() (any, error) {
		nbaHandler := nba.NewHandler(rapidApi, &storage, logger)
		if err := nbaHandler.Handler(ctx); err != nil {
			return nil, fmt.Errorf("NBA lambda handler failed: %w", err)
		}

		return nil, nil
	})

	_, err = wg.Wait()

	if err != nil {
		logger.Error("One or more handlers failed", "error", err)
		return fmt.Errorf("One or more handlers failed: %w", err)
	}

	logger.Info("Event successfully processed")
	return nil
}

func main() {
	lambda.Start(handler)
}
