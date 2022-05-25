package main

import "github.com/saantiaguilera/go-pipeline"

type (
	// nonDrawableStep is a custom step decorator over a UnitStep that overrides the drawing behavior
	// muting it.
	//
	// This is used in mutating steps that simply change the input/output between two steps so they
	// don't couple to another bussiness unit, hence we don't want them visualized in the graph
	// as they don't provide any value
	nonDrawableStep[I, O any] struct {
		pipeline.UnitStep[I, O]
	}
)

func newNonDrawableStep[I, O any](name string, unit pipeline.Unit[I, O]) nonDrawableStep[I, O] {
	return nonDrawableStep[I, O]{
		UnitStep: pipeline.NewUnitStep(name, unit),
	}
}

func (s nonDrawableStep[I, O]) Draw(g pipeline.Graph) {
	// Nothing
}
