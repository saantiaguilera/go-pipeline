package pipeline

import (
	"context"
)

type (
	// OptionalStep is a step that may or may not run depending on a statement
	OptionalStep[I, O any] struct {
		statement Statement[I]
		step      Step[I, O]
		def       Unit[I, O]
	}
)

// NewOptionalStep creates a step that may run the provided step if the statement evaluates correctly
// if the statement yields false, then the same input will be forwarded as output.
func NewOptionalStep[T any](stmt Statement[T], s Step[T, T]) OptionalStep[T, T] {
	return NewOptionalStepWithDefault(stmt, s, func(_ context.Context, in T) (T, error) {
		return in, nil
	})
}

// NewOptionalStepWithDefault creates a step that may run the provided step if the statement evaluates correctly
// if the statement yields false, then a default unit will be run to forward an output O
func NewOptionalStepWithDefault[I, O any](stmt Statement[I], s Step[I, O], def Unit[I, O]) OptionalStep[I, O] {
	return OptionalStep[I, O]{
		statement: stmt,
		step:      s,
		def:       def,
	}
}

func (c OptionalStep[I, O]) Draw(graph Graph) {
	graph.AddDecision(
		c.statement.Name(),
		func(graph Graph) {
			if c.step != nil {
				c.step.Draw(graph)
			}
		},
		func(graph Graph) {},
	)
}

func (c OptionalStep[I, O]) Run(ctx context.Context, in I) (O, error) {
	if c.statement.Evaluate(ctx, in) {
		return c.step.Run(ctx, in)
	}
	return c.def(ctx, in)
}
