package main

import (
	"context"

	"github.com/saantiaguilera/go-pipeline"
)

type (
	// varidicSequenceStep is a container of steps that will run sequentially.
	// If one of them fails, the chain is halted and the error is returned immediately
	//
	// all this steps must receive and return the same type of element since we use a slice
	// this is useful if we have a bunch of steps that all use the same type, as it is way less verbose to construct it
	// ```
	// newVaridicSequence[string](
	//   step1,
	//   step2,
	//   step3,
	//   step4,
	//   step5,
	//   ...
	// )
	// ```
	varidicSequenceStep[T any] []pipeline.Step[T, T]
)

func newVaridicSequence[T any](ss ...pipeline.Step[T, T]) varidicSequenceStep[T] {
	return ss
}

func (c varidicSequenceStep[T]) Draw(graph pipeline.Graph) {
	for _, s := range c {
		s.Draw(graph)
	}
}

func (c varidicSequenceStep[T]) Run(ctx context.Context, in T) (T, error) {
	res := in
	var err error
	for _, s := range c {
		if res, err = s.Run(ctx, res); err != nil {
			break
		}
	}
	return res, err
}
