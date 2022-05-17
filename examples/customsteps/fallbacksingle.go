package main

import (
	"context"

	"github.com/saantiaguilera/go-pipeline"
)

type (
	// singlefallback is a step that will execute a single unit and if it errors it will fallback with the returned error
	singlefallback[I, O any] struct {
		unit     pipeline.UnitStep[I, O]
		fallback func(context.Context, I, error) (O, error)
	}
)

func (c singlefallback[I, O]) Draw(graph pipeline.Graph) {
	graph.AddDecision(
		"errored",
		func(graph pipeline.Graph) {
			graph.AddActivity(c.unit.Name())
		},
		func(graph pipeline.Graph) {},
	)
}

func (c singlefallback[I, O]) Run(ctx context.Context, in I) (O, error) {
	res, err := c.unit.Run(ctx, in)
	if err != nil {
		return c.fallback(ctx, in, err)
	}
	return res, nil
}
