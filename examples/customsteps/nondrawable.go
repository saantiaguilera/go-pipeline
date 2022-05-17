package main

import "github.com/saantiaguilera/go-pipeline"

type (
	// nonDrawableStep is a decoration of a UnitStep that doesn't draw itself, thus
	// when representing a graph this step won't be shown.
	//
	// this may be useful if you have a step that mutates or transforms an input without
	// enriching the graph visualization/representation but rather keeping steps
	// decoupled between them
	nonDrawableStep[I, O any] struct {
		pipeline.UnitStep[I, O]
	}
)

func (s nonDrawableStep[I, O]) Draw(g pipeline.Graph) {
	// Nothing
}
