package pipeline

import (
	"context"

	"github.com/google/uuid"
)

type (
	// Unit of work to yield a result of type O (or an error in case of a failure) from a given input I
	Unit[I, O any] func(context.Context, I) (O, error)

	// UnitStep for making a unit of work.
	UnitStep[I, O any] struct {
		id   string
		name string
		fn   Unit[I, O]
	}

	// Step is runnable element that yields a result or error from a given input
	// A step can be drawn into a graph to represent it.
	Step[I, O any] interface {
		DrawableGraph

		// Run a step. Returns an error if this step fails to complete.
		// An input I is provided as a mean of communication between different units of work
		Run(context.Context, I) (O, error)
	}
)

// NewUnitStep creates an immutable stateless unit of work based on a Unit function
// You can use this implementation when your use-cases will be completely stateless (they don't rely on a service
// or anything that can be injected at the start and stay immutable for the lifetime of the process)
func NewUnitStep[I, O any](name string, run Unit[I, O]) UnitStep[I, O] {
	return UnitStep[I, O]{
		id:   uuid.New().String(),
		name: name,
		fn:   run,
	}
}

// ID is a unique identifier of this step. You can safely assume it wont be repeated and use it in any custom steps
// to enrich logic (eg. a circuit breaker / cache for IDs)
func (s UnitStep[I, O]) ID() string {
	return s.id
}

// Name to identify a step. You shouldn't assume this name is unique per step but rather use it to understand what this is / does / represent
func (s UnitStep[I, O]) Name() string {
	return s.name
}

// Run a step and yield a result of type O or an error if it failed.
// This operation is context-aware.
func (s UnitStep[I, O]) Run(ctx context.Context, in I) (O, error) {
	if err := ctx.Err(); err != nil {
		return *new(O), err
	}
	return s.fn(ctx, in)
}

// Draw this step in a graph
func (s UnitStep[I, O]) Draw(graph Graph) {
	graph.AddActivity(s.Name())
}
