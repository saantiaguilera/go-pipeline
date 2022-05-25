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
func newWashCarrotsStep() pipeline.UnitStep[[]Carrot, []Carrot] {
	return pipeline.NewUnitStep("wash_carrots_step", func(ctx context.Context, in []Carrot) ([]Carrot, error) {
		fmt.Printf("Washing %d carrots\n", len(in))
		time.Sleep(1 * time.Second) // Simulate time it takes to do this action

		for _, c := range in {
			c.Washed = true
		}
		return in, nil
	})
}
