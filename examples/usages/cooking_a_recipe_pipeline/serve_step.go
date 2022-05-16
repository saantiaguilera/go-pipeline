package main

import (
	"context"
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

func newServeStep() pipeline.UnitStep[DishContents, Dish] {
	return pipeline.NewUnitStep("serve_step", func(ctx context.Context, in DishContents) (Dish, error) {
		fmt.Printf("Serving dish with a salad of %d eggs and %d carrots (mixed: %v) and a cooked meat\n", len(in.Salad.Eggs), len(in.Salad.Carrots), in.Salad.Mixed)
		time.Sleep(1 * time.Second) // Simulate time it takes to do this action
		return Dish(in), nil
	})
}
