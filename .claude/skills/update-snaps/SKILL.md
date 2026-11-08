---
name: update-snaps
description: Regenerate go-snaps test snapshots and review every resulting change for regressions in fetched and parsed API results.
disable-model-invocation: true
---

Regenerate the snapshots and audit the result. Scope: `$ARGUMENTS` if given (a package such as `./internal/nba`), otherwise `./...`.

The snapshots hold the parsed API results from `GetGames`/`GetRaces`, so they say nothing about filtering, UID construction or `.ics` output. A clean diff only clears fetching and parsing.

1. Run `go test <scope>` without updating. Outside CI this already writes snapshots for new tests. If only snapshot assertions fail, or nothing fails, continue. If anything else fails, stop and report it, because regenerating won't fix it.
2. Run `UPDATE_SNAPS=true go test <scope>`.
3. Read `git diff HEAD -- '*.snap'` alongside the code diff (`git diff HEAD -- '*.go'` and `git diff main...HEAD -- '*.go'`). Run `git status --short` and read any untracked `.snap` or `.go` files too, because `git diff` omits them.
4. For each hunk, say what changed and which code change explains it. Flag these prominently:
   - games, races or fields removed, or a lower count
   - changed API IDs, which event UIDs are built from
   - changed start times, especially timezone shifts
   - any hunk with no matching code change
5. Finish with a verdict: safe to commit as far as the snapshots show, or the hunks that need a decision. If the code diff touches filtering, event building or `internal/calendar`, say the snapshots don't cover it.

Leave the regenerated snapshots uncommitted, and don't revert them unless asked.
