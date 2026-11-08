---
name: add-job
description: Add a new sport or schedule source to the calendar generator, covering the package, job registry, no_events alarm, tests and fixtures. Use whenever the user wants a new calendar, sport, league or data source, even if they don't say "job".
---

Add the job named in `$ARGUMENTS`. Read the closest existing package first and follow it: `internal/basketball` for a RapidAPI source, `internal/f1` for a direct HTTP source. If you have neither API docs nor a sample response, ask for one, because `.env` is off limits and the response shape can't be guessed.

## Wiring

A job is only complete when all of these are done. Nothing fails at build time if one is missed:

1. `internal/<sport>/handler.go`, plus `constant.go` for the base URL and IDs. `Handler` returns the result of `Upload(ctx, h.storage, "<sport>.ics")` directly so `calendar.ErrNoEvents` reaches `runJobs`.
2. `internal/app/app.go`: add a `<sport>Job` builder and an entry in `jobs`. The key is the `-jobs` flag value and the alarm dimension.
3. `infra/alarms.tf`: add the same key to `for_each` on the `no_events` alarm. Without it, a source that silently breaks never alerts.
4. `internal/testdata/testdata.go`: add a `//go:embed <sport>/*` line for the new fixture directory.
5. `README.md` sport list and the package list under Architecture in `CLAUDE.md`.

## Behaviour to get right up front

Most past fixes to these jobs were cases the first version missed, so check each against the API's real data:

- **Cancelled or postponed events are skipped.** Find the API's status values rather than assuming them. Fixtures may only show the normal status.
- **To-be-decided participants or times** are marked or skipped, not rendered as an empty title or a bogus start time.
- **A failed page or team fails the whole job.** Return the error from `GetGames`/`GetRaces`, because uploading partial results would remove real events from subscribers' calendars.
- **Duplicates** are dropped when a source lists the same event twice, for example across teams or live and upcoming feeds.
- **Event IDs are stable** and built from the API's own ID as `<sport>-<id>`. Calendar clients use the UID to track events, so a changed format duplicates every event.
- **Pagination is capped** with a `maxPages` constant to protect the RapidAPI quota if the end-of-results flag never arrives.

## Tests

Write handler tests as in `internal/basketball/handler_test.go`: `testutil.NewMockHttpClient`, `testutil.MemoryStorage`, a saved response under `internal/testdata/<sport>/`. Cover the uploaded calendar's IDs, titles and times, plus one test for each skip rule above. Snapshots only capture parsed API results, so filtering needs explicit assertions.

Run `go test -race ./...` and `make lint`, then read the new snapshot in full.

## Finish

Leave the work uncommitted. Tell the user that it needs two separate production steps, and don't run either: `make deploy` for the alarm and a version tag for the application.
