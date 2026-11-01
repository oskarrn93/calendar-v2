package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/oskarrn93/calendar-v2/internal/awsutil"
	"github.com/oskarrn93/calendar-v2/internal/basketball"
	"github.com/oskarrn93/calendar-v2/internal/calendar"
	"github.com/oskarrn93/calendar-v2/internal/config"
	"github.com/oskarrn93/calendar-v2/internal/esport"
	"github.com/oskarrn93/calendar-v2/internal/f1"
	"github.com/oskarrn93/calendar-v2/internal/football"
	"github.com/oskarrn93/calendar-v2/internal/httpclient"
	"github.com/oskarrn93/calendar-v2/internal/logging"
	"github.com/oskarrn93/calendar-v2/internal/nba"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
	"github.com/oskarrn93/calendar-v2/internal/shutdown"
)

type Handler func(context.Context) error

func Handlers(ctx context.Context, logger *slog.Logger) (map[string]Handler, error) {
	appConfig, err := config.Initialize()
	if err != nil {
		return nil, err
	}

	storage, err := awsutil.NewS3Storage(ctx, appConfig.S3Bucket, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 storage: %w", err)
	}

	httpClient := httpclient.New()
	rapidApi := rapidapi.New(httpClient, appConfig.RapidApiKey)

	return map[string]Handler{
		"Football":   football.NewHandler(rapidApi, storage).Handler,
		"NBA":        nba.NewHandler(rapidApi, storage).Handler,
		"Esport":     esport.NewHandler(rapidApi, storage, logger).Handler,
		"Basketball": basketball.NewHandler(rapidApi, storage).Handler,
		"F1":         f1.NewHandler(httpClient, storage, logger).Handler,
	}, nil
}

// RunCommand runs a single handler, for the per-sport binaries in cmd/.
func RunCommand(name string) {
	ctx, stop := shutdown.Context(context.Background())
	defer stop()

	logger := logging.New()

	logger.Info(name + " command")

	handlers, err := Handlers(ctx, logger)
	if err != nil {
		panic(err)
	}

	handler, ok := handlers[name]
	if !ok {
		panic(fmt.Errorf("unknown handler %q", name))
	}

	err = handler(ctx)
	if errors.Is(err, calendar.ErrNoEvents) {
		logger.Warn("No events found, keeping previous calendar", "handler", name, "error", err)
		return
	}
	if err != nil {
		panic(fmt.Errorf("%s handler failed: %w", name, err))
	}

	logger.Info("Event successfully processed")
}
