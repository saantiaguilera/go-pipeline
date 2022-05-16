package main

import (
	"context"
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

func newTurnOnOvenStep() pipeline.UnitStep[Oven, Oven] {
	return pipeline.NewUnitStep("turn_on_oven_step", func(ctx context.Context, in Oven) (Oven, error) {
		fmt.Print("Turning oven on\n")
		time.Sleep(1 * time.Second) // Simulate time it takes to do this action
		in.Ignited = true
		return in, nil
	})
}
