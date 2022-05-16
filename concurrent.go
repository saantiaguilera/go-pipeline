package pipeline

import (
	"context"
	"errors"
)

type (
	ConcurrentStep[I, O any] struct {
		steps  []Step[I, O]
		reduce reducer[O]
	}

	reducer[O any] func(context.Context, O, O) (O, error)
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

	mres, err := spawnAsync(c.steps, func(s Step[I, O]) (O, error) {
		return s.Run(ctx, in)
	})

	if err != nil {
		return *new(O), err
	}

	acc := mres[0]
	for _, v := range mres[1:] {
		acc, err = c.reduce(ctx, acc, v)
		if err != nil {
			return *new(O), err
		}
	}
	return acc, nil
}
