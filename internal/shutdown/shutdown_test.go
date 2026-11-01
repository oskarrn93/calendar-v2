package shutdown

import (
	"context"
	"syscall"
	"testing"
	"time"
)

func TestContextCancelledOnSignal(t *testing.T) {
	for _, sig := range []syscall.Signal{syscall.SIGINT, syscall.SIGTERM} {
		t.Run(sig.String(), func(t *testing.T) {
			ctx, stop := Context(context.Background())
			defer stop()

			if err := syscall.Kill(syscall.Getpid(), sig); err != nil {
				t.Fatalf("failed to send signal: %v", err)
			}

			select {
			case <-ctx.Done():
			case <-time.After(2 * time.Second):
				t.Fatal("context was not cancelled after signal")
			}
		})
	}
}

func TestContextCancelledByStop(t *testing.T) {
	ctx, stop := Context(context.Background())
	stop()

	if ctx.Err() == nil {
		t.Fatal("expected context to be cancelled after stop")
	}
}
