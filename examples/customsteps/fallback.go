package main

import (
	"context"

	"github.com/saantiaguilera/go-pipeline"
)

type (
	// fallbackStep is a step that will execute a fallback step when an error occurs in the main one
	fallbackStep[I, O any] struct {
		step     pipeline.Step[I, O]
		fallback pipeline.Step[I, O]
	}
)

func (c fallbackStep[I, O]) Draw(graph pipeline.Graph) {
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

func (c fallbackStep[I, O]) Run(ctx context.Context, in I) (O, error) {
	res, err := c.step.Run(ctx, in)
	if err != nil {
		return c.fallback.Run(ctx, in)
	}
	return res, nil
}
