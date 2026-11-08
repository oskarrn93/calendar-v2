package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/oskarrn93/calendar-v2/internal/awsutil"
	"github.com/oskarrn93/calendar-v2/internal/basketball"
	"github.com/oskarrn93/calendar-v2/internal/calendar"
	"github.com/oskarrn93/calendar-v2/internal/config"
	"github.com/oskarrn93/calendar-v2/internal/esport"
	"github.com/oskarrn93/calendar-v2/internal/f1"
	"github.com/oskarrn93/calendar-v2/internal/football"
	"github.com/oskarrn93/calendar-v2/internal/httpclient"
	"github.com/oskarrn93/calendar-v2/internal/nba"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
)

type job func(context.Context) error

type deps struct {
	httpClient  *resty.Client
	rapidApiKey string
	storage     awsutil.Storage
	logger      *slog.Logger
}

// rapidApi errors on a missing key so only jobs that need RapidAPI require it.
func (d deps) rapidApi() (rapidapi.RapidApi, error) {
	if d.rapidApiKey == "" {
		return rapidapi.RapidApi{}, errors.New("RAPIDAPI_KEY is not set")
	}
	return rapidapi.New(d.httpClient, d.rapidApiKey), nil
}

type jobBuilder interface {
	build(deps) (job, error)
}

type footballJob struct{}

func (footballJob) build(d deps) (job, error) {
	r, err := d.rapidApi()
	if err != nil {
		return nil, err
	}
	return football.NewHandler(r, d.storage).Handler, nil
}

type nbaJob struct{}

func (nbaJob) build(d deps) (job, error) {
	r, err := d.rapidApi()
	if err != nil {
		return nil, err
	}
	return nba.NewHandler(r, d.storage).Handler, nil
}

type esportJob struct{}

func (esportJob) build(d deps) (job, error) {
	r, err := d.rapidApi()
	if err != nil {
		return nil, err
	}
	return esport.NewHandler(r, d.storage, d.logger).Handler, nil
}

type basketballJob struct{}

func (basketballJob) build(d deps) (job, error) {
	r, err := d.rapidApi()
	if err != nil {
		return nil, err
	}
	return basketball.NewHandler(r, d.storage).Handler, nil
}

type f1Job struct{}

func (f1Job) build(d deps) (job, error) {
	return f1.NewHandler(d.httpClient, d.storage, d.logger).Handler, nil
}

// The keys are the job names accepted by the CLI's -jobs flag. New jobs also
// need adding to the no_events alarm in infra/alarms.tf.
var jobs = map[string]jobBuilder{
	"football":   footballJob{},
	"nba":        nbaJob{},
	"esport":     esportJob{},
	"basketball": basketballJob{},
	"f1":         f1Job{},
}

// Run runs the named jobs concurrently, or every job when names is empty.
// Names must match exactly; they are validated before any setup so a typo
// isn't masked by missing config.
func Run(ctx context.Context, logger *slog.Logger, names []string) error {
	selected, err := selectJobs(names)
	if err != nil {
		return err
	}

	appConfig, err := config.Initialize()
	if err != nil {
		return err
	}

	storage, err := awsutil.NewS3Storage(ctx, appConfig.S3Bucket, logger)
	if err != nil {
		return fmt.Errorf("failed to create S3 storage: %w", err)
	}

	d := deps{
		httpClient:  httpclient.New(),
		rapidApiKey: appConfig.RapidApiKey,
		storage:     storage,
		logger:      logger,
	}

	built := make(map[string]job, len(selected))
	for name, builder := range selected {
		j, err := builder.build(d)
		if err != nil {
			return fmt.Errorf("%s job: %w", name, err)
		}
		built[name] = j
	}

	return runJobs(ctx, logger, built)
}

func runJobs(ctx context.Context, logger *slog.Logger, jobs map[string]job) error {
	logger.Info("Running jobs", "jobs", slices.Sorted(maps.Keys(jobs)))

	// Collect every job's error rather than cancel on the first, so one
	// failing source doesn't stop the other calendars from updating.
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	for name, run := range jobs {
		wg.Go(func() {
			err := run(ctx)
			if errors.Is(err, calendar.ErrNoEvents) {
				// The event field is what the no_events alarm in infra/alarms.tf filters on.
				logger.Warn("No events found, keeping previous calendar", "event", "no_events", "job", name, "error", err)
				return
			}
			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("%s job failed: %w", name, err))
				mu.Unlock()
			}
		})
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return err
	}

	logger.Info("Jobs successfully processed")
	return nil
}

func selectJobs(names []string) (map[string]jobBuilder, error) {
	if len(names) == 0 {
		return jobs, nil
	}

	selected := make(map[string]jobBuilder, len(names))
	for _, name := range names {
		builder, ok := jobs[name]
		if !ok {
			return nil, fmt.Errorf("unknown job %q, valid jobs: %s", name, strings.Join(slices.Sorted(maps.Keys(jobs)), ", "))
		}
		selected[name] = builder
	}

	return selected, nil
}
