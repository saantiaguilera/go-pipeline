package main

import (
	"context"
	"sync"

	"github.com/saantiaguilera/go-pipeline"
)

// processDishConcurrently is a custom step that decorates 2 inner steps and runs them concurrently
//
// We use this custom step instead of a ConcurrentStep to showcase how can we achieve custom behaviors
// on our own
// In this case, we get the benefit of having different input/ouptut between each step, hence creating
// a more decoupled environment.
//   This processDishConcurrently step is a pipeline.Step[MealMaterials, DishContent]
//     Inner SaladStep is a pipeline.Step[MealMaterials, Salad]
//     Inner MeatStep is a pipeline.Step[MealMaterials, CookedMeat]
//
// If we wanted to achieve this with a ConcurrentStep, we would need to wrap both inner steps with one
// that changes their output into a DishContent since they don't return that type, and later reduce
// both results into a single one.
// This approach with a ConcurrentStep can be seen inside the `main.go` for the 'process meat' pipeline
type processDishConcurrently struct {
	SaladStep pipeline.Step[MealMaterials, Salad]
	MeatStep  pipeline.Step[MealMaterials, CookedMeat]
}

func (p processDishConcurrently) Draw(g pipeline.Graph) {
	g.AddConcurrency(p.SaladStep.Draw, p.MeatStep.Draw)
}

func (p processDishConcurrently) Run(ctx context.Context, in MealMaterials) (DishContents, error) {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	errs := make(chan error, 2)
	var salad *Salad
	var meat *CookedMeat

	go func() {
		defer wg.Done()
		v, err := p.SaladStep.Run(ctx, in)
		if err != nil {
			errs <- err
			return
		}
		salad = &v
	}()

	go func() {
		defer wg.Done()
		v, err := p.MeatStep.Run(ctx, in)
		if err != nil {
			errs <- err
			return
		}
		meat = &v
	}()

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return DishContents{}, <-errs
	}
	return DishContents{
		Salad: *salad,
		Meat:  *meat,
	}, nil
}
