package pipeline

import (
	"context"
	"errors"
)

type (
	// ConcurrentStep wraps multiple steps of a given Input/Output and runs them concurrently, later
	// reducing them into a single output of the same type.
	ConcurrentStep[I, O any] struct {
		steps  []Step[I, O]
		reduce reducer[O]
	}

	// reducer reduces two values of the same type in a single one
	reducer[O any] func(context.Context, O, O) (O, error)

	// concurrentResult is a discriminated union of a result or error.
	concurrentResult[T any] struct {
		Ret T
		Err error
	}
)

// NewConcurrentStep creates a step that will run each of the inner steps concurrently.
// The step will wait for all of the steps to finish before returning.
//
// If one of them fails, the step will wait until everyone finishes and after that return the first encountered error.
//
// This step (as all the others) doesn't handle panics. Be careful since this step creates goroutines and the panics
// not necessarilly will be signaled in the same goroutine as the origin call.
// Make sure to handle panics on your own if your code is unsafe (through decorations / deferrals in steps / etc)
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

// Run the step concurrently, if one of them fails an error will be returned.
//
// This step waits for all of the concurrent ones to finish.
//
// Note that this step may use goroutines and (as all other steps) doesn't handle panics,
// hence it is advise to handle them on your own if you can't guarantee a panic-safe environment.
func (c ConcurrentStep[I, O]) Run(ctx context.Context, in I) (O, error) {
	if len(c.steps) == 0 {
		return *new(O), errors.New("cannot run with empty concurrent steps")
	}

	mch := c.runConcurrently(ctx, c.steps, in)

	var acc O
	var err error
	for i := 0; i < len(c.steps); i++ {
		v := <-mch
		if err != nil {
			continue // we want all steps to finish, so simply cut here and wait for next step to end.
		}

		if v.Err != nil {
			err = v.Err // step errored.
			continue
		}

		if i == 0 {
			acc = v.Ret
		} else {
			acc, err = c.reduce(ctx, acc, v.Ret)
		}
	}
	return acc, err
}

// Run a number of workers concurrently, waiting for all of them to finish.
// After they're all done, if one of them failed the error is returned.
// If more than one fails, the last error is returned
//
// Note: this method doesn't recover from panics in goroutines. To be cohesive across the whole
// API none of the steps handle panics to let the client handle them on their own
// (through decorations / same steps with deferrals / whatever he wants to)
func (c ConcurrentStep[I, O]) runConcurrently(
	ctx context.Context,
	workers []Step[I, O],
	in I,
) <-chan concurrentResult[O] {

	ch := make(chan concurrentResult[O], len(workers))
	if len(workers) > 1 {
		for i := 0; i < len(workers); i++ {
			go c.runStep(ctx, in, workers[i], ch)
		}
	} else { // avoid concurrency, no need to spawn and wait just use current
		c.runStep(ctx, in, workers[0], ch)
	}
	return ch
}

func (c ConcurrentStep[I, O]) runStep(
	ctx context.Context,
	in I,
	step Step[I, O],
	ch chan<- concurrentResult[O],
) {

	res, err := step.Run(ctx, in)
	ch <- concurrentResult[O]{
		Ret: res,
		Err: err,
	}
}
