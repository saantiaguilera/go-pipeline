package pipeline

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

type (
	ConcurrentStep[I, O any] struct {
		steps  []Step[I, O]
		reduce reducer[O]
	}

	reducer[O any] func(context.Context, O, O) (O, error)

	// concurrentSlice that allows appending values of type T concurrently
	concurrentSlice[T any] struct {
		sync.RWMutex
		v []T
	}
)

// NewConcurrentStep creates a step that will run each of the inner steps concurrently.
// The step will wait for all of the steps to finish before returning.
//
// If one of them fails, the step will wait until everyone finishes and after that return the first encountered error.
func NewConcurrentStep[I, O any](steps []Step[I, O], reduce reducer[O]) ConcurrentStep[I, O] {
	return ConcurrentStep[I, O]{
		steps:  steps,
		reduce: reduce,
	}
}

func (c ConcurrentStep[I, O]) Draw(graph Graph) {
	if len(c.steps) > 0 {
		var forkSteps []GraphDrawer
		for _, s := range c.steps {
			forkSteps = append(forkSteps, s.Draw)
		}

		graph.AddConcurrency(forkSteps...)
	}
}

func (c ConcurrentStep[I, O]) Run(ctx context.Context, in I) (O, error) {
	if len(c.steps) == 0 {
		return *new(O), errors.New("cannot run with empty concurrent steps")
	}

	mres, err := c.runConcurrently(ctx, c.steps, in)
	if err != nil {
		return *new(O), err
	}

	acc := mres[0]
	for i := 1; i < len(mres); i++ {
		acc, err = c.reduce(ctx, acc, mres[i])
		if err != nil {
			return *new(O), err
		}
	}
	return acc, nil
}

// Run a number of workers concurrently, waiting for all of them to finish.
// After they're all done, if one of them failed the error is returned.
// If more than one fails, the last error is returned
func (c ConcurrentStep[I, O]) runConcurrently(
	ctx context.Context,
	workers []Step[I, O],
	in I,
) ([]O, error) {

	var wg sync.WaitGroup
	var errResult atomic.Value
	var mergedRes concurrentSlice[O]

	wg.Add(len(workers))
	if len(workers) > 1 {
		for i := 0; i < len(workers); i++ {
			go c.runStep(&wg, &errResult, &mergedRes, ctx, in, workers[i])
		}
		wg.Wait()
	} else { // avoid concurrency, no need to spawn and wait just use current
		c.runStep(&wg, &errResult, &mergedRes, ctx, in, workers[0])
	}

	if err, ok := errResult.Load().(error); ok && err != nil {
		return nil, err
	}
	return mergedRes.Get(), nil
}

func (c ConcurrentStep[I, O]) runStep(
	wg *sync.WaitGroup,
	errResult *atomic.Value,
	mergedRes *concurrentSlice[O],
	ctx context.Context,
	in I,
	step Step[I, O],
) {

	res, err := step.Run(ctx, in)
	if err != nil {
		errResult.CompareAndSwap(nil, err)
	}

	mergedRes.Append(res)
	wg.Done()
}

// Read slice T stored inside
func (s *concurrentSlice[T]) Get() []T {
	s.RLock()
	defer s.RUnlock()
	return s.v
}

// Append safely inside T slice
func (s *concurrentSlice[T]) Append(t T) {
	s.Lock()
	defer s.Unlock()
	s.v = append(s.v, t)
}
