package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-resty/resty/v2"
	"github.com/oskarrn93/calendar-v2/internal/awsutil"
	"github.com/oskarrn93/calendar-v2/internal/basketball"
	"github.com/oskarrn93/calendar-v2/internal/config"
	"github.com/oskarrn93/calendar-v2/internal/esport"
	"github.com/oskarrn93/calendar-v2/internal/football"
	"github.com/oskarrn93/calendar-v2/internal/logging"
	"github.com/oskarrn93/calendar-v2/internal/nba"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
	"github.com/oskarrn93/calendar-v2/internal/util"
)

func handler(ctx context.Context, event json.RawMessage) error {
	logger := logging.New()

	logger.Info("Received Event", "event", event)

	appConfig, err := config.Initialize()
	if err != nil {
		return err
	}

	storage, err := awsutil.NewS3Storage(ctx, appConfig.S3Bucket, logger)
	if err != nil {
		return fmt.Errorf("failed to create S3 storage: %w", err)
	}

	rapidApi := rapidapi.New(resty.New(), appConfig.RapidApi)

	handlers := map[string]func(context.Context) error{
		"Football":   football.NewHandler(rapidApi, storage, logger).Handler,
		"NBA":        nba.NewHandler(rapidApi, storage, logger).Handler,
		"Esport":     esport.NewHandler(rapidApi, storage, logger).Handler,
		"Basketball": basketball.NewHandler(rapidApi, storage, logger).Handler,
	}

	wg := util.NewWaitGroup()

	for name, sportHandler := range handlers {
		wg.Run(func() error {
			if err := sportHandler(ctx); err != nil {
				return fmt.Errorf("%s handler failed: %w", name, err)
			}

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		logger.Error("One or more handlers failed", "error", err)
		return fmt.Errorf("one or more handlers failed: %w", err)
	}

	logger.Info("Event successfully processed")
	return nil
}

func main() {
	lambda.Start(handler)
}
