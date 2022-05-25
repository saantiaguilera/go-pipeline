package main

import (
	"context"
	"sync"

	"github.com/saantiaguilera/go-pipeline"
)

// processVegetablesConcurrently is a custom step that decorates 2 inner steps and runs them concurrently
//
// We use this custom step instead of a ConcurrentStep to showcase how can we achieve custom behaviors
// on our own
// In this case, we get the benefit of having different input/ouptut between each step, hence creating
// a more decoupled environment.
//   This processVegetablesConcurrently step is a pipeline.Step[MealMaterials, Vegetables]
//     Inner SaladStep is a pipeline.Step[[]Egg, []CutEgg]
//     Inner MeatStep is a pipeline.Step[[]Carrot, []CutCarrot]
//
// If we wanted to achieve this with a ConcurrentStep, we would need to wrap both inner steps with one
// that changes the MealMaterials input into a []Egg / []Carrot since they don't receive a MealMaterials.
// Also we would need to wrap their outputs since they return a []CutEgg / []CutCarrot and we expect a
// Vegetables type. Finally we would need to reduce both results into a single one.
// This approach with a ConcurrentStep can be seen inside the `main.go` for the 'process meat' pipeline
type processVegetablesConcurrently struct {
	EggStep    pipeline.Step[[]Egg, []CutEgg]
	CarrotStep pipeline.Step[[]Carrot, []CutCarrot]
}

func (p processVegetablesConcurrently) Draw(g pipeline.Graph) {
	g.AddConcurrency(p.EggStep.Draw, p.CarrotStep.Draw)
}

func (p processVegetablesConcurrently) Run(ctx context.Context, in MealMaterials) (Vegetables, error) {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	errs := make(chan error, 2)
	var cutEggs *[]CutEgg
	var cutCarrots *[]CutCarrot

	go func() {
		defer wg.Done()
		v, err := p.EggStep.Run(ctx, in.Eggs)
		if err != nil {
			errs <- err
			return
		}
		cutEggs = &v
	}()

	go func() {
		defer wg.Done()
		v, err := p.CarrotStep.Run(ctx, in.Carrots)
		if err != nil {
			errs <- err
			return
		}
		cutCarrots = &v
	}()

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return Vegetables{}, <-errs
	}
	return Vegetables{
		Eggs:    *cutEggs,
		Carrots: *cutCarrots,
	}, nil
}
