---
name: update-snaps
description: Regenerate go-snaps test snapshots and review every resulting change for regressions such as dropped events.
disable-model-invocation: true
---

Regenerate the snapshots and audit the result. Scope: `$ARGUMENTS` if given (a package such as `./internal/nba`), otherwise `./...`.

1. Run `go test <scope>` without updating. If only snapshot assertions fail, continue. If anything else fails, stop and report it, because regenerating won't fix it.
2. Run `UPDATE_SNAPS=true go test <scope>`.
3. Read `git diff -- '*.snap'` alongside the code diff (`git diff -- '*.go'` and `git diff main...HEAD -- '*.go'`).
4. For each hunk, say what changed and which code change explains it. Flag these prominently:
   - events or fields removed, or a lower event count
   - changed event IDs, which subscribers' calendars use to identify events
   - changed start or end times, especially timezone shifts
   - any hunk with no matching code change
5. Finish with a verdict: safe to commit, or the hunks that need a decision.

Leave the regenerated snapshots uncommitted, and don't revert them unless asked.
