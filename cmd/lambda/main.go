package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/oskarrn93/calendar-v2/internal/app"
	"github.com/oskarrn93/calendar-v2/internal/logging"
	"github.com/oskarrn93/calendar-v2/internal/util"
)

func handler(ctx context.Context, event json.RawMessage) error {
	logger := logging.New()

	logger.Info("Received Event", "event", event)

	handlers, err := app.Handlers(ctx, logger)
	if err != nil {
		return err
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
