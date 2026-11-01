package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"
	"github.com/oskarrn93/calendar-v2/internal/awsutil"
	"github.com/oskarrn93/calendar-v2/internal/basketball"
	"github.com/oskarrn93/calendar-v2/internal/config"
	"github.com/oskarrn93/calendar-v2/internal/esport"
	"github.com/oskarrn93/calendar-v2/internal/f1"
	"github.com/oskarrn93/calendar-v2/internal/football"
	"github.com/oskarrn93/calendar-v2/internal/logging"
	"github.com/oskarrn93/calendar-v2/internal/nba"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
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

	rapidApi := rapidapi.New(resty.New(), appConfig.RapidApi)

	return map[string]Handler{
		"Football":   football.NewHandler(rapidApi, storage, logger).Handler,
		"NBA":        nba.NewHandler(rapidApi, storage, logger).Handler,
		"Esport":     esport.NewHandler(rapidApi, storage, logger).Handler,
		"Basketball": basketball.NewHandler(rapidApi, storage, logger).Handler,
		"F1":         f1.NewHandler(resty.New(), appConfig.F1, storage, logger).Handler,
	}, nil
}

// RunCommand runs a single handler, for the per-sport binaries in cmd/.
func RunCommand(name string) {
	ctx := context.Background()
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

	if err := handler(ctx); err != nil {
		panic(fmt.Errorf("%s handler failed: %w", name, err))
	}

	logger.Info("Event successfully processed")
}
