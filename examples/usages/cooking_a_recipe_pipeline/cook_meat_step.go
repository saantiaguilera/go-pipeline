package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

func newCookMeatStep() pipeline.UnitStep[CookingTools, CookedMeat] {
	return pipeline.NewUnitStep("cook_meat_step", func(ctx context.Context, in CookingTools) (CookedMeat, error) {
		fmt.Printf("Cooking meat\n")
		time.Sleep(1 * time.Second) // Simulate time it takes to do this action

		if !in.Oven.Ignited {
			return CookedMeat{}, errors.New("cannot cook meat with the oven turned off")
		}
		return CookedMeat{ /* stuff */ }, nil
	})
}
