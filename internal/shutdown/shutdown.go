// Package shutdown provides graceful shutdown on OS termination signals.
package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Context returns a copy of parent that is cancelled when the process receives
// SIGINT or SIGTERM. The returned stop function releases the signal handler
// and should be deferred by the caller.
func Context(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, stop := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)

	// Restore default signal behaviour once cancelled, so a second signal
	// force-kills the process instead of being swallowed while work unwinds.
	context.AfterFunc(ctx, stop)

	return ctx, stop
}
