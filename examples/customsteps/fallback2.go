package pipeline

import (
	"context"

	"github.com/saantiaguilera/go-pipeline"
)

type (
	// fallback2Step is a step that will execute a fallback step when an error occurs in the main one
	// in contrast to the fallbackStep, this allows us to mutate the input when an error occurs
	// and go to the fallback with a mutated input
	fallback2Step[I, E, O any] struct {
		step      pipeline.Step[I, O]
		transform func(context.Context, I, error) (E, error)
		fallback  pipeline.Step[E, O]
	}
)

func (c fallback2Step[I, E, O]) Draw(graph pipeline.Graph) {
	graph.AddDecision(
		"errored",
		func(graph pipeline.Graph) {
			if c.step != nil {
				c.step.Draw(graph)
			}
		},
		func(graph pipeline.Graph) {
			if c.fallback != nil {
				c.fallback.Draw(graph)
			}
		},
	)
}

func (c fallback2Step[I, E, O]) Run(ctx context.Context, in I) (O, error) {
	res, err := c.step.Run(ctx, in)
	if err != nil {
		mr, err := c.transform(ctx, in, err)
		if err != nil {
			return *new(O), err
		}
		return c.fallback.Run(ctx, mr)
	}
	return res, nil
}
