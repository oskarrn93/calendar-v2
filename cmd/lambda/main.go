package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/oskarrn93/calendar-v2/internal/app"
	"github.com/oskarrn93/calendar-v2/internal/calendar"
	"github.com/oskarrn93/calendar-v2/internal/logging"
)

func handler(ctx context.Context, event json.RawMessage) error {
	logger := logging.New()

	logger.Info("Received Event", "event", event)

	handlers, err := app.Handlers(ctx, logger)
	if err != nil {
		return err
	}

	// Collect every handler's error rather than cancel on the first, so one
	// failing source doesn't stop the other calendars from updating.
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	for name, sportHandler := range handlers {
		wg.Go(func() {
			err := sportHandler(ctx)
			if errors.Is(err, calendar.ErrNoEvents) {
				logger.Warn("No events found, keeping previous calendar", "handler", name, "error", err)
				return
			}
			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("%s handler failed: %w", name, err))
				mu.Unlock()
			}
		})
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		logger.Error("One or more handlers failed", "error", err)
		return fmt.Errorf("one or more handlers failed: %w", err)
	}

	logger.Info("Event successfully processed")
	return nil
}

func main() {
	lambda.Start(handler)
}
