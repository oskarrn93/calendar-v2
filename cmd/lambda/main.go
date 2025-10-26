package lambda

import (
	"context"
	"encoding/json"
	"fmt"

	"oskarrn93/calendar-v2/internal/awsutil"
	"oskarrn93/calendar-v2/internal/config"
	"oskarrn93/calendar-v2/internal/logging"
	"oskarrn93/calendar-v2/internal/nba"
	"oskarrn93/calendar-v2/internal/rapidapi"

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

	nbaHandler := nba.NewHandler(rapidApi, &storage, logger)
	if err := nbaHandler.Handler(ctx); err != nil {
		logger.Error("NBA handler failed", "error", err)
		return err
	}

	logger.Info("Event successfully processed")
	return nil
}

func main() {
	lambda.Start(handler)
}
