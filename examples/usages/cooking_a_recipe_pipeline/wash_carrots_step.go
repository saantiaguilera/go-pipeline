package main

import (
	"context"
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

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
