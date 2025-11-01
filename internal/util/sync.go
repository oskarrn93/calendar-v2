package util

import (
	"context"
	"errors"
	"slices"
	"sync"

	"golang.org/x/sync/errgroup"
)

const MaxConcurrentLimit = 10

type SyncSlice[T any] struct {
	data []T
	mu   sync.RWMutex
}

func (s *SyncSlice[T]) Append(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = append(s.data, item)
}

func (s *SyncSlice[T]) Get() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// return a copy of the slice since slices are reference types (pointer)
	return slices.Clone(s.data)
}

type WaitGroup[T any] struct {
	Ctx     context.Context
	eg      *errgroup.Group
	results SyncSlice[T]
	errors  SyncSlice[error]
}

func (wg *WaitGroup[T]) Run(f func() (T, error)) {
	wg.eg.Go(func() error {
		result, err := f()
		if err != nil {
			wg.errors.Append(err)
		}
		wg.results.Append(result)

		// Always return nil to prevent cancelling other goroutines if one fails, let them all finish
		return nil
	})
}

func (wg *WaitGroup[T]) Wait() ([]T, error) {
	err := wg.eg.Wait()

	errs := wg.errors.Get()
	resultError := errors.Join(err, errors.Join(errs...))

	result := wg.results.Get()

	// Return all results and the combined error, let the client handle partial failures
	return result, resultError
}

func (wg *WaitGroup[T]) SetLimit(limit int) {
	if limit <= 0 {
		limit = -1 // negative value means no limit
	}

	wg.eg.SetLimit(limit)
}

func NewWaitGroup[T any](ctx context.Context) *WaitGroup[T] {
	eg, egCtx := errgroup.WithContext(ctx)
	eg.SetLimit(MaxConcurrentLimit)

	wg := &WaitGroup[T]{
		Ctx: egCtx,
		eg:  eg,
	}

	return wg
}
