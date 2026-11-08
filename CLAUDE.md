# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```sh
make build                       # builds into ./bin
make test                        # go test -v ./...
go test -race ./...              # what CI runs; make test omits -race
go test ./internal/esport -run TestHandlerUploadsFilteredCalendar
UPDATE_SNAPS=true make test      # rewrite go-snaps snapshots
make lint                        # fails unless golangci-lint matches .golangci-lint-version
make run JOBS=nba,f1             # local run, needs godotenv and .env
```

`make run` uploads the generated `.ics` files to the real bucket named by `S3_BUCKET_NAME`. `cmd/lambda` only runs inside the Lambda runtime, so use `cmd/calendar` locally.

## Architecture

A daily EventBridge-scheduled Lambda fetches sports schedules, builds `.ics` files and uploads them to S3, where CloudFront serves them. `cmd/calendar` and `cmd/lambda` are both thin wrappers around `app.Run`.

`internal/app` holds the job registry. Each job is a sport package (`football`, `nba`, `esport`, `basketball`, `f1`) whose `Handler` fetches events, filters them into a `calendar.Calendar` and uploads a single `<sport>.ics`. All jobs except `f1` go through `internal/rapidapi`; `f1` uses the HTTP client directly.

- Jobs run concurrently and errors are collected, so one failing source doesn't stop the other calendars from updating.
- `calendar.ErrNoEvents` isn't a failure: the previous calendar is kept and a `no_events` log event is emitted, which the alarm in `infra/alarms.tf` filters on.
- Adding a job means registering it in `jobs` in `internal/app/app.go` and adding it to the `no_events` alarm.
- The RapidAPI key comes from `RAPIDAPI_KEY` locally and from the SSM parameter named by `RAPIDAPI_KEY_PARAMETER` in Lambda. It is only required when a selected job uses RapidAPI.

## Tests

Handler tests use `testutil.NewMockHttpClient` with `httpmock`, saved API responses from `internal/testdata`, and `testutil.MemoryStorage` to inspect uploaded calendars. Snapshots live in each package's `__snapshots__/`. Review snapshot diffs rather than accepting them, since they can hide dropped events.

## Deployment

- Application: publishing a version tag (`v1.2.3`) builds the image and pushes it to ECR. The manually triggered Deploy GitHub Actions workflow then rolls out a chosen tag.
- Infrastructure: `make deploy` runs `terraform apply` in `infra/` against production. State is remote in S3.
