package main

import (
	"context"

	"github.com/saantiaguilera/go-pipeline"
)

type (
	// EventStep is a custom step that, given this particular scenario, will allow us to forward
	// an initial input EventID (that later mutates into other stuff and gets discarded) to the latest step
	// in this chain
	//
	// This allows us to create a "kind of" defer scenario, where we use something and at the end
	// we use it again, forgetting about it during the whole other process.
	//
	// This is needed because we use initially the EventID to unmarshall it into the data we will process
	// but at the end of the pipeline, if everything went correctly, we want to flag it as processed
	// to avoid processing it again.
	EventStep[T any] struct {
		pipe    pipeline.Step[EventID, T]
		process pipeline.Step[EventID, EventID]
	}
)

func NewEventStep[T any](p pipeline.Step[EventID, T], e pipeline.Step[EventID, EventID]) EventStep[T] {
	return EventStep[T]{
		pipe:    p,
		process: e,
	}
}

func (s EventStep[T]) Draw(g pipeline.Graph) {
	s.pipe.Draw(g)
	s.process.Draw(g)
}

// Run runs first the pipe step and if it goes correctly it runs another step
// that only handles the eventID (regardless of the output struct it returned the
// first pipe)
func (s EventStep[T]) Run(ctx context.Context, in EventID) (T, error) {
	out, err := s.pipe.Run(ctx, in)
	if err != nil {
		return *new(T), err
	}
	if _, err = s.process.Run(ctx, in); err != nil {
		return out, err
	}
	return out, nil
}
