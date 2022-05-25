package main

import (
	"context"
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

// We can create a step normally through an inner function that will do the work. Here we are encapsulating
// the whole process in a single constructor method which helps us mantain isolation.
// In this particular case, we also inject a parameter 'times' that is static throughout the graph evaluation
//
// Note that we are still coupling this behavior to the pipeline API.
//
// If you want to avoid that, see the cut_carrots_step.go
func newCutEggsStep(times int) pipeline.UnitStep[[]Egg, []CutEgg] {
	return pipeline.NewUnitStep("cut_eggs_step", func(ctx context.Context, in []Egg) ([]CutEgg, error) {
		fmt.Printf("Cutting %d eggs\n", len(in)*times)
		time.Sleep(1 * time.Second) // Simulate time it takes to do this action

		return make([]CutEgg, len(in)*times), nil
	})
}
