package util

import (
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

const MaxConcurrentLimit = 10

// WaitGroup collects the error of every task instead of cancelling the
// remaining ones on the first failure, so a partial run still completes.
type WaitGroup struct {
	eg   errgroup.Group
	mu   sync.Mutex
	errs []error
}

func NewWaitGroup() *WaitGroup {
	wg := &WaitGroup{}
	wg.eg.SetLimit(MaxConcurrentLimit)

	return wg
}

func (wg *WaitGroup) Run(f func() error) {
	wg.eg.Go(func() error {
		if err := f(); err != nil {
			wg.mu.Lock()
			wg.errs = append(wg.errs, err)
			wg.mu.Unlock()
		}

		// Always return nil so errgroup never cancels the other tasks
		return nil
	})
}

func (wg *WaitGroup) Wait() error {
	_ = wg.eg.Wait()

	return errors.Join(wg.errs...)
}
