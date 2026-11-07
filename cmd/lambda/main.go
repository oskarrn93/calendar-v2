package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/oskarrn93/calendar-v2/internal/app"
	"github.com/oskarrn93/calendar-v2/internal/logging"
)

func handler(ctx context.Context, event json.RawMessage) error {
	logger := logging.New()

	logger.Info("Received Event", "event", event)

	return app.Run(ctx, logger, nil)
}

func main() {
	lambda.Start(handler)
}
