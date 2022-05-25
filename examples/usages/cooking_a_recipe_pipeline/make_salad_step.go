package main

import (
	"context"
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

// We can create a step normally through an inner function that will do the work. Here we are encapsulating
// the whole process in a single constructor method which helps us mantain isolation, although we are still
// coupling this behavior to the pipeline API.
//
// If you want to avoid that, see the cut_carrots_step.go
func newMakeSaladStep() pipeline.UnitStep[Vegetables, Salad] {
	return pipeline.NewUnitStep("make_salad_step", func(ctx context.Context, in Vegetables) (Salad, error) {
		fmt.Printf("Making salad with %d eggs and %d carrots\n", len(in.Eggs), len(in.Carrots))
		time.Sleep(1 * time.Second) // Simulate time it takes to do this action
		return Salad{
			Vegetables: in,
			Mixed:      true,
		}, nil
	})
}
