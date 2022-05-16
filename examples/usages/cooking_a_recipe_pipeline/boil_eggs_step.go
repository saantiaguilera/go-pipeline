package main

import (
	"context"
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

// We can create a step normally through an inner function that will do the work. Here we are encapsulating
// the whole process in a single constructor method which helps us mantain isolation
func newBoilEggsStep() pipeline.UnitStep[[]Egg, []Egg] {
	return pipeline.NewUnitStep("boil_eggs_step", func(ctx context.Context, in []Egg) ([]Egg, error) {
		fmt.Printf("Boiling %d eggs\n", len(in))
		time.Sleep(1 * time.Second) // Simulate time it takes to do this action

		for _, e := range in {
			e.Boiled = true
		}
		return in, nil
	})
}
