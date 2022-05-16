package main

import (
	"context"
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

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
