package main

import (
	"context"
	"flag"
	"os"
	"strings"

	"github.com/oskarrn93/calendar-v2/internal/app"
	"github.com/oskarrn93/calendar-v2/internal/logging"
	"github.com/oskarrn93/calendar-v2/internal/shutdown"
)

func main() {
	jobs := flag.String("jobs", "", "comma-separated jobs to run, e.g. nba,f1 (default: all)")
	flag.Parse()

	if err := run(parseJobs(*jobs)); err != nil {
		os.Exit(1)
	}
}

// run is split from main so deferred cleanup happens before os.Exit.
func run(names []string) error {
	ctx, stop := shutdown.Context(context.Background())
	defer stop()

	logger := logging.New()

	err := app.Run(ctx, logger, names)
	if err != nil {
		logger.Error("Run failed", "error", err)
	}
	return err
}

func parseJobs(s string) []string {
	var names []string
	for name := range strings.SplitSeq(s, ",") {
		if name = strings.ToLower(strings.TrimSpace(name)); name != "" {
			names = append(names, name)
		}
	}
	return names
}
