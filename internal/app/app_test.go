package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"testing"

	"github.com/oskarrn93/calendar-v2/internal/calendar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunJobs(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	errBoom := errors.New("boom")

	tests := []struct {
		name    string
		jobs    map[string]error
		wantErr error
	}{
		{name: "success", jobs: map[string]error{"a": nil}},
		{name: "no events is not an error", jobs: map[string]error{"a": nil, "b": fmt.Errorf("wrapped: %w", calendar.ErrNoEvents)}},
		{name: "job error is returned and others still run", jobs: map[string]error{"a": nil, "c": errBoom}, wantErr: errBoom},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ran := make(chan string, len(tt.jobs))
			jobs := make(map[string]job, len(tt.jobs))
			for name, err := range tt.jobs {
				jobs[name] = func(context.Context) error { ran <- name; return err }
			}

			err := runJobs(context.Background(), logger, jobs)
			close(ran)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
			var gotRan []string
			for name := range ran {
				gotRan = append(gotRan, name)
			}
			assert.ElementsMatch(t, slices.Collect(maps.Keys(tt.jobs)), gotRan)
		})
	}
}

func TestSelectJobs(t *testing.T) {
	tests := []struct {
		name    string
		names   []string
		want    []string
		wantErr string
	}{
		{name: "all when none selected", names: nil, want: slices.Collect(maps.Keys(jobs))},
		{name: "only selected", names: []string{"nba", "f1"}, want: []string{"nba", "f1"}},
		{name: "unknown job", names: []string{"nba", "nope"}, wantErr: `unknown job "nope"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := selectJobs(tt.names)

			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.ElementsMatch(t, tt.want, slices.Collect(maps.Keys(got)))
		})
	}
}

func TestBuildRequiresRapidApiKeyOnlyWhenUsed(t *testing.T) {
	d := deps{logger: slog.New(slog.DiscardHandler)}

	_, err := jobs["f1"].build(d)
	require.NoError(t, err)

	_, err = jobs["nba"].build(d)
	require.ErrorContains(t, err, "RAPIDAPI_KEY is not set")
}
