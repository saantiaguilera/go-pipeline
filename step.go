package pipeline

import (
	"context"

	"github.com/google/uuid"
)

type (
	Unit[I, O any] func(context.Context, I) (O, error)

	// Step interface for making a unit work.
	UnitStep[I, O any] struct {
		id   string
		name string
		fn   Unit[I, O]
	}

	// Step is runnable element that yields a result or error from a given input
	Step[I, O any] interface {
		DrawableDiagram

		// Run a step. Returns an error if this step fails to complete.
		// An input I is provided as a mean of communication between different units of work
		Run(context.Context, I) (O, error)
	}
)

// NewStep creates an immutable stateless unit of work based on a function that matches the Runnable contract.
// You can use this implementation when your use-cases will be completely stateless (they don't rely on a service
// or anything that can be injected at the start and stay immutable for the lifetime of the process)
func NewUnitStep[I, O any](name string, run Unit[I, O]) UnitStep[I, O] {
	return UnitStep[I, O]{
		id:   uuid.New().String(),
		name: name,
		fn:   run,
	}
}

func (s UnitStep[I, O]) ID() string {
	return s.id
}

func (s UnitStep[I, O]) Name() string {
	return s.name
}

func (s UnitStep[I, O]) Run(ctx context.Context, in I) (O, error) {
	if err := ctx.Err(); err != nil {
		return *new(O), err
	}
	return s.fn(ctx, in)
}

func (s UnitStep[I, O]) Draw(graph GraphDiagram) {
	graph.AddActivity(s.Name())
}
